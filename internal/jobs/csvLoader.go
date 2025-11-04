package jobs

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"time"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

type (
	CSVLoader[
		T types.Client |
			types.Exercise |
			types.Hyperparams |
			types.RawWorkout,
	] struct {
		B          *sbjobqueue.Batch
		S          *types.State
		Q          *dal.SyncQueries
		FileDir    string
		ClientName string
		FileChunk  io.Reader
		Opts       *sbcsv.Opts
		WriteFunc  func(
			ctxt context.Context,
			state *types.State,
			queries *dal.SyncQueries,
			values ...T,
		) error
	}

	rawTrainingLog struct {
		DatePerformed time.Time
		Session       uint16
		Exercise      string
		Weight        types.Kilogram
		Sets          float64
		Reps          int32
		Effort        types.RPE
		DataDir       string
	}

	csvTimeSeriesData struct {
		Time      float64
		XPosition float64
		YPosition float64
	}

	WorkoutFileChunk struct {
		Headers    []byte
		FullFile   []byte
		StartIdx   int
		EndIdx     int
		Started    bool
		Stop       bool
		ReachedEnd bool
		readPos    int
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
) *WorkoutFileChunk {
	return &WorkoutFileChunk{
		Headers:    headers,
		FullFile:   fullFile,
		StartIdx:   startIdx,
		EndIdx:     endIdx,
		Started:    (startIdx == 0),
		ReachedEnd: false,
		Stop:       false,
	}
}

func (f *WorkoutFileChunk) Read(buf []byte) (n int, err error) {
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

func (w *CSVLoader[T]) JobType(_ types.CSVLoaderJob) {}

func (w *CSVLoader[T]) Batch() *sbjobqueue.Batch {
	return w.B
}

func (w *CSVLoader[T]) Run(ctxt context.Context) (opErr error) {
	params := []T{}
	switch typedParams := any((*[]T)(&params)).(type) {
	case *[]types.RawWorkout:
		if w.ClientName == "" {
			opErr = sberr.Wrap(
				types.CSVLoaderJobQueueErr,
				"ClientName must not be empty when parsing workout data",
			)
			return
		}
		fileChunk, ok := w.FileChunk.(*WorkoutFileChunk)
		if !ok {
			opErr = sberr.Wrap(
				types.CSVLoaderJobQueueErr,
				"Workouts must be chunked with WorkoutFileChunk due to workout boundaries",
			)
			return
		}

		firstWorkoutIDSet, lastWorkoutIDSet := fileChunk.StartIdx == 0, false
		prevWorkoutID, lastWorkoutID := types.WorkoutID{}, types.WorkoutID{}

		if opErr = sbcsv.LoadReader(w.FileChunk, &sbcsv.LoadOpts{
			Opts:          *w.Opts,
			RequestedCols: sbcsv.ReqColsForStruct[rawTrainingLog](),
			Op: func(
				o *sbcsv.Opts,
				rowIdx int,
				row []string,
				reqCols []sbcsv.RequestedCols,
			) error {
				rawData, err := sbcsv.RowToStruct[rawTrainingLog](o, row, reqCols)
				if err != nil {
					return err
				}

				iterID := types.WorkoutID{
					ClientEmail:   w.ClientName,
					Session:       rawData.Session,
					DatePerformed: rawData.DatePerformed,
				}

				fileChunk.Stop = (lastWorkoutIDSet && iterID != lastWorkoutID)
				if fileChunk.Stop {
					return nil
				}

				if !firstWorkoutIDSet {
					firstWorkoutIDSet = true
					prevWorkoutID = iterID
				}
				if !lastWorkoutIDSet && fileChunk.ReachedEnd {
					lastWorkoutIDSet = true
					lastWorkoutID = iterID
				}
				if iterID != prevWorkoutID {
					fileChunk.Started = true
					(*typedParams) = append(
						(*typedParams),
						types.RawWorkout{WorkoutID: iterID},
					)
					prevWorkoutID = iterID
				}
				if !fileChunk.Started {
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

				(*typedParams)[len((*typedParams))-1].Exercises = append(
					(*typedParams)[len((*typedParams))-1].Exercises,
					types.RawExerciseData{
						Name:    rawData.Exercise,
						Weight:  rawData.Weight,
						Sets:    rawData.Sets,
						Reps:    rawData.Reps,
						Effort:  rawData.Effort,
						BarPath: variants,
					},
				)

				return nil
			},
		}); opErr != nil {
			return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
		}
	default:
		if opErr = sbcsv.LoadReader(w.FileChunk, &sbcsv.LoadOpts{
			Opts:          *w.Opts,
			RequestedCols: sbcsv.ReqColsForStruct[T](),
			Op:            sbcsv.RowToStructOp(&params),
		}); opErr != nil {
			return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
		}
	}

	if opErr = w.WriteFunc(ctxt, w.S, w.Q, params...); opErr != nil {
		return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
	}

	return
}

func (w *CSVLoader[T]) parseWorkoutDataDir(
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
			res[setNum-1] = types.BarPathVideo(path)
		case "wla":
			tsData, err := w.loadWLACSVData(path)
			if err != nil {
				return sberr.AppendError(
					sberr.Wrap(
						types.InvalidDataDirErr,
						"Weight lifting analysis export file malformed",
					),
					err,
				)
			}
			res[setNum-1] = tsData
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

func (w *CSVLoader[T]) loadTimeSeriesCSVData(
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

	return types.BarPathTimeSeriesData(rawTimeSeriesData), nil
}

func (w *CSVLoader[T]) loadWLACSVData(
	path string,
) (types.BarPathVariant, error) {
	reqCols := []sbcsv.RequestedCols{
		{Name: "Time (s)"},
		{Name: "displacement (horizontal, cm)"},
		{Name: "displacement (vertical, cm)"},
	}

	f, err := os.Open(path)
	if err != nil {
		return types.BarPathVariant{}, err
	}
	sc := bufio.NewScanner(f)
	i := 0
	for ; i < 17 && sc.Scan(); i++ {
		_ = sc.Text()
	}
	if err := sc.Err(); err != nil {
		return types.BarPathVariant{}, err
	}
	if i != 17 {
		return types.BarPathVariant{}, types.WLAFileMissingData
	}

	rawTimeSeriesData := types.RawTimeSeriesData{}
	if err := sbcsv.LoadReader(f, &sbcsv.LoadOpts{
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
					X: types.Meter(xpos / 100), Y: types.Meter(ypos / 100),
				},
			)
			return nil
		},
	}); err != nil {
		return types.BarPathVariant{}, err
	}

	return types.BarPathTimeSeriesData(rawTimeSeriesData), nil
}
