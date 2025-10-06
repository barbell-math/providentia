package jobs

import (
	"context"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
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
		ClientName string
		FileChunk  []byte
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
)

var (
	setDataFileRe         = `^Set([0-9]+).(csv|mp4)$`
	compiledSetDataFileRe = regexp.MustCompile(setDataFileRe)
)

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
		if opErr = sbcsv.LoadBytes(w.FileChunk, &sbcsv.LoadOpts{
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

				variants, err := w.parseWorkoutDataDir(
					rawData.DataDir, int(math.Ceil(rawData.Sets)),
				)
				if err != nil {
					return err
				}

				iterID := types.WorkoutID{
					ClientEmail:   w.ClientName,
					Session:       rawData.Session,
					DatePerformed: rawData.DatePerformed,
				}
				if len((*typedParams)) == 0 || (*typedParams)[len((*typedParams))-1].WorkoutID != iterID {
					(*typedParams) = append((*typedParams), types.RawWorkout{WorkoutID: iterID})
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
		if opErr = sbcsv.LoadBytes(w.FileChunk, &sbcsv.LoadOpts{
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
			res[setNum] = types.BarPathVideo(path)
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
			res[setNum] = tsData
		}

		return err
	}); parseErr != nil {
		parseErr = sberr.AppendError(types.InvalidDataDirErr, parseErr)
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
		return types.BarPathVariant{}, sberr.AppendError(
			types.InvalidDataDirErr, err,
		)
	}

	return types.BarPathTimeSeriesData(rawTimeSeriesData), nil
}
