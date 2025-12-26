package jobs

import (
	"context"
	"io"
	"io/fs"
	"iter"
	"math"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	rawWorkoutData struct {
		DatePerformed time.Time
		Session       uint16
		Exercise      string
		Weight        types.Kilogram
		Sets          float64
		Reps          int32
		Effort        types.RPE
		DataDir       string
	}

	workoutFileChunk struct {
		Headers    []byte
		FullFile   []byte
		StartIdx   int
		EndIdx     int
		Started    bool
		Stop       bool
		ReachedEnd bool
		readPos    int
	}

	workoutCSVLoader struct {
		B           *sbjobqueue.Batch
		S           *types.State
		Tx          pgx.Tx
		UID         uint64
		ClientEmail string
		FileDir     string
		FileChunk   *workoutFileChunk
		Opts        *sbcsv.Opts
		*types.BarPathCalcHyperparams
		*types.BarPathTrackerHyperparams
	}

	CSVWorkoutLoaderOpts struct {
		*sbcsv.Opts
		*types.BarPathCalcHyperparams
		*types.BarPathTrackerHyperparams
		Files iter.Seq2[string, error]
		Batch *sbjobqueue.Batch
	}
)

var (
	setDataFileRe         = `^Set([0-9]+).(csv|mp4)$`
	compiledSetDataFileRe = regexp.MustCompile(setDataFileRe)
)

func NewWorkoutFileChunk(
	headers []byte,
	fullFile []byte,
	startIdx int,
	endIdx int,
) *workoutFileChunk {
	return &workoutFileChunk{
		Headers:    headers,
		FullFile:   fullFile,
		StartIdx:   startIdx,
		EndIdx:     endIdx,
		Started:    (startIdx == 0),
		ReachedEnd: false,
		Stop:       false,
	}
}

func (f *workoutFileChunk) Read(buf []byte) (n int, err error) {
	if f.Stop || f.readPos >= len(f.Headers)+len(f.FullFile[f.StartIdx:]) {
		err = io.EOF
		return
	}
	f.ReachedEnd = (f.readPos+f.StartIdx-len(f.Headers) > f.EndIdx)
	if f.readPos < len(f.Headers) {
		headersN := copy(buf, f.Headers[f.readPos:])
		f.readPos += headersN
		n += headersN
	}
	if n < len(buf) && f.readPos < len(f.Headers)+len(f.FullFile[f.StartIdx:]) {
		endIdx := slices.Index(
			f.FullFile[f.StartIdx+f.readPos-len(f.Headers):],
			'\n',
		)
		if endIdx == -1 {
			endIdx = len(f.FullFile)
		} else {
			endIdx += f.StartIdx + f.readPos - len(f.Headers) + 1
		}
		dataN := copy(buf[n:], f.FullFile[f.StartIdx+f.readPos-len(f.Headers):endIdx])
		f.readPos += dataN
		n += dataN
	}
	return
}

func UploadWorkoutsFromCSV(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *CSVWorkoutLoaderOpts,
) error {
	wait := false
	if opts.Batch == nil {
		wait = true
		opts.Batch, _ = sbjobqueue.BatchWithContext(ctxt)
	}

	for file, err := range opts.Files {
		if err != nil {
			return err
		}
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			formatJobLogLine("UploadWorkoutsFromCSV", 0, "Processing data file"),
			"File", file,
		)

		clientEmail := strings.TrimSuffix(path.Base(file), path.Ext(file))
		if _, err := mail.ParseAddress(clientEmail); err != nil {
			return sberr.AppendError(
				types.CSVLoaderJobQueueErr,
				sberr.Wrap(
					err,
					"The name of each workout file must follow the format: <client email>.csv\nGot: %s",
					clientEmail,
				),
			)
		}

		fileChunks, err := sbcsv.ChunkFile(
			file, NewWorkoutFileChunk, state.WorkoutCSVFileChunks,
		)
		if err != nil {
			return err
		}
		for _, chunk := range fileChunks {
			if chunk.EndIdx-chunk.StartIdx <= 0 {
				continue
			}
			state.CSVLoaderJobQueue.Schedule(&workoutCSVLoader{
				S:                         state,
				Tx:                        tx,
				B:                         opts.Batch,
				UID:                       UID_CNTR.Add(1),
				ClientEmail:               clientEmail,
				FileDir:                   path.Dir(file),
				FileChunk:                 chunk,
				Opts:                      opts.Opts,
				BarPathCalcHyperparams:    opts.BarPathCalcHyperparams,
				BarPathTrackerHyperparams: opts.BarPathTrackerHyperparams,
			})
		}
	}

	if wait {
		return opts.Batch.Wait()
	}
	return nil
}

func (w *workoutCSVLoader) JobType(_ types.CSVLoaderJob) {}

func (w *workoutCSVLoader) Batch() *sbjobqueue.Batch {
	return w.B
}

func (w *workoutCSVLoader) formatLogLine(msg string) string {
	return formatJobLogLine("workoutCSVLoader", w.UID, msg)
}

func (w *workoutCSVLoader) Run(ctxt context.Context) (opErr error) {
	w.S.Log.Log(ctxt, sblog.VLevel(3), w.formatLogLine("Starting..."))

	firstWorkoutIdSet, lastWorkoutIdSet := w.FileChunk.StartIdx == 0, false
	prevWorkoutId, lastWorkoutId := types.WorkoutId{}, types.WorkoutId{}

	params := []types.Workout{}
	if opErr = sbcsv.LoadReader(w.FileChunk, &sbcsv.LoadOpts{
		Opts:          *w.Opts,
		RequestedCols: sbcsv.ReqColsForStruct[rawWorkoutData](),
		Op: func(
			o *sbcsv.Opts,
			rowIdx int,
			row []string,
			reqCols []sbcsv.RequestedCols,
		) error {
			rawData, err := sbcsv.RowToStruct[rawWorkoutData](o, row, reqCols)
			if err != nil {
				return err
			}

			iterId := types.WorkoutId{
				ClientEmail:   w.ClientEmail,
				Session:       rawData.Session,
				DatePerformed: rawData.DatePerformed,
			}

			w.FileChunk.Stop = (lastWorkoutIdSet && iterId != lastWorkoutId)
			if w.FileChunk.Stop {
				return nil
			}

			if !firstWorkoutIdSet {
				firstWorkoutIdSet = true
				prevWorkoutId = iterId
			}
			if !lastWorkoutIdSet && w.FileChunk.ReachedEnd {
				lastWorkoutIdSet = true
				lastWorkoutId = iterId
			}
			if iterId != prevWorkoutId {
				w.FileChunk.Started = true
				params = append(params, types.Workout{WorkoutId: iterId})
				prevWorkoutId = iterId
			}
			if !w.FileChunk.Started {
				return nil
			}

			var variants []types.BarPathVariant
			if rawData.DataDir != "" {
				variants, err = w.parseWorkoutDataDir(
					path.Join(w.FileDir, rawData.DataDir),
					int(math.Ceil(rawData.Sets)),
				)
				if err != nil {
					return err
				}
			}

			iterExerciseData := types.ExerciseData{
				Name:   rawData.Exercise,
				Weight: rawData.Weight,
				Sets:   rawData.Sets,
				Reps:   rawData.Reps,
				Effort: rawData.Effort,
				PhysData: make(
					[]types.Optional[types.PhysicsData],
					int(math.Ceil(rawData.Sets)),
				),
			}

			if len(variants) > 0 {
				if err := RunPhysicsJobs(ctxt, w.S, w.Tx, PhysicsOpts{
					BarPathCalcParams:    w.BarPathCalcHyperparams,
					BarTrackerCalcParams: w.BarPathTrackerHyperparams,
					RawData:              variants,
					ExerciseData:         &iterExerciseData,
				}); err != nil {
					return err
				}
			}

			params[len(params)-1].Exercises = append(
				params[len(params)-1].Exercises,
				iterExerciseData,
			)

			return nil
		},
	}); opErr != nil {
		goto errReturn
	}

	// This is unfortunate... but it has to be done because a single transaction
	// is backed by a single conn which is not thread safe.
	w.B.Lock()
	if opErr = dal.CreateWorkouts(ctxt, w.S, w.Tx, params); opErr != nil {
		goto errReturn
	}
	w.B.Unlock()

	w.S.Log.Log(
		ctxt, sblog.VLevel(3),
		w.formatLogLine("Finished loading workout data"),
		"NumRows", len(params),
	)
	return

errReturn:
	w.S.Log.Error(w.formatLogLine("Encountered error"), "Error", opErr)
	return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
}

func (w *workoutCSVLoader) parseWorkoutDataDir(
	dir string,
	numSets int,
) (res []types.BarPathVariant, parseErr error) {
	if dir == "" {
		return
	}

	var fi fs.FileInfo
	if fi, parseErr = os.Stat(dir); parseErr != nil {
		return
	} else if !fi.IsDir() {
		parseErr = sberr.Wrap(
			types.InvalidDataDirErr, "'%s' was not a dir but must be", dir,
		)
		return
	}

	res = make([]types.BarPathVariant, numSets)
	if parseErr = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir {
			return nil
		}
		if d.IsDir() {
			return sberr.Wrap(
				types.InvalidDataDirErr,
				"Data dirs should not have sub-dirs: '%s' was a dir", path,
			)
		}

		idxs := compiledSetDataFileRe.FindStringSubmatchIndex(d.Name())
		if len(idxs) == 0 {
			return sberr.Wrap(
				types.InvalidDataDirErr,
				"Data file had invalid name '%s'. Name must match the following regex: %s",
				path, setDataFileRe,
			)
		}

		rawSetNum, err := strconv.ParseInt(d.Name()[idxs[2]:idxs[3]], 10, 0)
		setNum := int(rawSetNum)
		if err != nil {
			return sberr.AppendError(types.InvalidDataDirErr, err)
		}
		if setNum <= 0 || setNum > numSets+1 {
			return sberr.Wrap(
				types.InvalidDataDirErr,
				"Set num (%d) out of allowed range [1, %d]",
				setNum, numSets,
			)
		}

		ext := d.Name()[idxs[4]:idxs[5]]
		switch ext {
		case "mp4":
			res[setNum-1] = types.BarPathVariant{
				Flag:      types.VideoBarPathData,
				VideoPath: path,
			}
		// TODO
		// case "wla":
		// 	tsData, err := w.loadWLACSVData(path)
		// 	if err != nil {
		// 		return sberr.AppendError(
		// 			sberr.Wrap(
		// 				types.InvalidDataDirErr,
		// 				"Weight lifting analysis export file malformed",
		// 			),
		// 			err,
		// 		)
		// 	}
		// 	res[setNum-1] = tsData
		case "csv":
			tsData, err := w.loadTimeSeriesCSVData(path)
			if err != nil {
				return sberr.AppendError(
					sberr.Wrap(
						types.InvalidDataDirErr,
						"Time series csv file malformed",
					),
					err,
				)
			}
			res[setNum-1] = tsData
		}

		return err
	}); parseErr != nil {
		return
	}

	return
}

func (w *workoutCSVLoader) loadTimeSeriesCSVData(
	path string,
) (types.BarPathVariant, error) {
	reqCols := []sbcsv.RequestedCols{
		{Name: "Time"}, {Name: "XPos"}, {Name: "YPos"},
	}

	rawTimeSeriesData := types.RawTimeSeriesData{}
	if err := sbcsv.LoadFile(path, &sbcsv.LoadOpts{
		Opts:          *w.Opts,
		RequestedCols: reqCols,
		Op: func(
			o *sbcsv.Opts,
			rowIdx int,
			row []string,
			reqCols []sbcsv.RequestedCols,
		) error {
			time, err := strconv.ParseFloat(row[reqCols[0].Idx], 64)
			if err != nil {
				return err
			}
			xpos, err := strconv.ParseFloat(row[reqCols[1].Idx], 64)
			if err != nil {
				return err
			}
			ypos, err := strconv.ParseFloat(row[reqCols[2].Idx], 64)
			if err != nil {
				return err
			}

			rawTimeSeriesData.TimeData = append(
				rawTimeSeriesData.TimeData, types.Second(time),
			)
			rawTimeSeriesData.PositionData = append(
				rawTimeSeriesData.PositionData,
				types.Vec2[types.Meter, types.Meter]{
					X: types.Meter(xpos), Y: types.Meter(ypos),
				},
			)
			return nil
		},
	}); err != nil {
		return types.BarPathVariant{}, err
	}

	return types.BarPathVariant{
		Flag:       types.TimeSeriesBarPathData,
		TimeSeries: rawTimeSeriesData,
	}, nil
}
