package jobs

import (
	"context"
	"io/fs"
	"iter"
	"math"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbcsv"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sberrs"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbjobqueue"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sblog"
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

	workoutCSVLoader struct {
		B           *sbjobqueue.Batch
		S           *types.State
		Tx          pgx.Tx
		UID         int64
		ClientEmail string
		File        string
		Opts        *sbcsv.Opts
		*types.BarPathCalcHyperparams
		*types.BarPathTrackerHyperparams

		DataDirs      []string
		PhysDataBatch *sbjobqueue.Batch
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

		uid := UID_CNTR.Add(1)
		state.Log.Log(
			ctxt, sblog.VLevel(3),
			formatJobLogLine("UploadWorkoutsFromCSV", uid, "Processing data file"),
			"File", file,
		)
		state.CSVLoaderJobQueue.Schedule(&workoutCSVLoader{
			S:                         state,
			Tx:                        tx,
			B:                         opts.Batch,
			UID:                       uid,
			ClientEmail:               clientEmail,
			File:                      file,
			Opts:                      opts.Opts,
			BarPathCalcHyperparams:    opts.BarPathCalcHyperparams,
			BarPathTrackerHyperparams: opts.BarPathTrackerHyperparams,
		})
	}

	if wait {
		return opts.Batch.Wait()
	}
	return nil
}

func (w *workoutCSVLoader) JobType(_ types.CSVLoaderJob) {}
func (w *workoutCSVLoader) Batch() *sbjobqueue.Batch     { return w.B }

func (w *workoutCSVLoader) formatLogLine(msg string) string {
	return formatJobLogLine("workoutCSVLoader", w.UID, msg)
}

func (w *workoutCSVLoader) Run(ctxt context.Context) (opErr error) {
	w.S.Log.Log(ctxt, sblog.VLevel(3), w.formatLogLine("Starting..."))
	w.PhysDataBatch, _ = sbjobqueue.BatchWithContext(ctxt)

	var f *os.File
	params := []types.Workout{}
	prevWorkoutId := types.WorkoutId{}
	reqCols := sbcsv.ReqColsForStruct[rawWorkoutData]()
	fieldIdxs := []int{}

	fieldIdxs, opErr = sbcsv.ReqColsToFieldIdxs[rawWorkoutData](reqCols)
	if opErr != nil {
		goto errReturn
	}

	f, opErr = os.Open(w.File)
	if opErr != nil {
		goto errReturn
	}

	if opErr = sbcsv.LoadReader(f, &sbcsv.LoadOpts{
		Opts:          *w.Opts,
		RequestedCols: reqCols,
		Op: func(
			o *sbcsv.Opts,
			rowIdx int,
			row []string,
			reqCols []sbcsv.RequestedCols,
		) error {
			rawData, err := sbcsv.RowToStruct[rawWorkoutData](o, row, reqCols, fieldIdxs)
			if err != nil {
				return err
			}

			iterId := types.WorkoutId{
				ClientEmail:   w.ClientEmail,
				Session:       rawData.Session,
				DatePerformed: rawData.DatePerformed,
			}
			if iterId != prevWorkoutId {
				if err := w.scheduleWorkoutPhysDataJobs(ctxt, params); err != nil {
					return err
				}
				params = append(params, types.Workout{WorkoutId: iterId})
				prevWorkoutId = iterId
			}

			w.DataDirs = append(w.DataDirs, rawData.DataDir)
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

			params[len(params)-1].Exercises = append(
				params[len(params)-1].Exercises,
				iterExerciseData,
			)
			return nil
		},
	}); opErr != nil {
		goto errReturn
	}

	if opErr = w.scheduleWorkoutPhysDataJobs(ctxt, params); opErr != nil {
		goto errReturn
	}

	if opErr = w.PhysDataBatch.Wait(); opErr != nil {
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

func (w *workoutCSVLoader) scheduleWorkoutPhysDataJobs(
	ctxt context.Context,
	params []types.Workout,
) error {
	if len(params) == 0 {
		// This will happen when the very first workout is being created and the
		// prev workout ID is still zero initialized.
		return nil
	}
	lastWorkout := &params[len(params)-1]

	for i, dataDir := range w.DataDirs {
		var err error
		var variants []types.BarPathVariant

		if dataDir != "" {
			variants, err = w.parseWorkoutDataDir(
				path.Join(path.Dir(w.File), dataDir),
				int(math.Ceil(lastWorkout.Exercises[i].Sets)),
			)
			if err != nil {
				return err
			}
		}

		if len(variants) > 0 {
			if err := RunPhysicsJobs(ctxt, w.S, w.Tx, PhysicsOpts{
				Batch:                w.PhysDataBatch,
				BarPathCalcParams:    w.BarPathCalcHyperparams,
				BarTrackerCalcParams: w.BarPathTrackerHyperparams,
				RawData:              variants,
				ExerciseData:         &params[len(params)-1].Exercises[i],
			}); err != nil {
				return err
			}
		}
	}

	w.DataDirs = w.DataDirs[:0]
	return nil
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
