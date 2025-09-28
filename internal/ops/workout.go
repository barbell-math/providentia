package ops

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

type (
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

	createSingleWorkoutParams struct {
		w                    *types.RawWorkout
		barPathCalcParams    *types.BarPathCalcHyperparams
		barTrackerCalcParams *types.BarPathTrackerHyperparams
		batch                *sbjobqueue.Batch
		clientCache          *dal.IdCache[string, int64]
		exerciseCache        *dal.IdCache[string, int32]
		bufWriter            *dal.BufferedWriter[dal.BulkCreateTrainingLogsParams]
	}
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	data ...types.RawWorkout,
) (opErr error) {
	batch, _ := sbjobqueue.BatchWithContext(ctxt)
	clientCache := dal.NewClientIdCache(state.Global.PerRequestIdCacheSize)
	exerciseCache := dal.NewExerciseIdCache(state.Global.PerRequestIdCacheSize)
	bufWriter := dal.NewBufferedWriter(
		state.Global.BatchSize,
		dal.Q.BulkCreateTrainingLogs,
		func() (err error) {
			if err := batch.Wait(); err != nil {
				return err
			}
			return
		},
	)
	params := createSingleWorkoutParams{
		barPathCalcParams:    barPathCalcParams,
		barTrackerCalcParams: barTrackerCalcParams,
		batch:                batch,
		clientCache:          &clientCache,
		exerciseCache:        &exerciseCache,
		bufWriter:            &bufWriter,
	}

	for _, iterW := range data {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		params.w = &iterW
		if opErr = createSingleWorkout(
			ctxt, state, queries, params,
		); opErr != nil {
			return
		}
	}

	if opErr = bufWriter.Flush(ctxt, queries); opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotAddWorkoutErr, dal.FormatErr(opErr),
		)
		return
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Added new workouts",
		"NumWorkouts", len(data),
	)
	return
}

func EnsureWorkoutsExist(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	data ...types.RawWorkout,
) (opErr error) {
	return
}

func createSingleWorkout(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	params createSingleWorkoutParams,
) (opErr error) {
	if opErr = validateWorkout(params); opErr != nil {
		opErr = sberr.AppendError(types.InvalidWorkoutErr, opErr)
		return
	}

	var iterClientId int64
	if iterClientId, opErr = params.clientCache.Get(
		ctxt, queries, params.w.ClientEmail,
	); opErr != nil {
		opErr = sberr.AppendError(
			types.InvalidWorkoutErr,
			sberr.Wrap(
				types.CouldNotFindRequestedClientErr,
				"Unknown Email: %s", params.w.ClientEmail,
			),
			opErr,
		)
		return
	}

	for i, iterE := range params.w.Exercises {
		var iterExerciseId int32
		if iterExerciseId, opErr = params.exerciseCache.Get(
			ctxt, queries, iterE.Name,
		); opErr != nil {
			opErr = sberr.AppendError(
				types.InvalidWorkoutErr,
				types.MalformedWorkoutExerciseErr,
				sberr.Wrap(
					types.CouldNotFindRequestedExerciseErr,
					"Unknown Exercise: %s", iterE.Name,
				),
				opErr,
			)
			return
		}

		if opErr = params.bufWriter.Write(
			ctxt, queries,
			dal.BulkCreateTrainingLogsParams{
				ClientID:         iterClientId,
				ExerciseID:       iterExerciseId,
				DatePerformed:    dal.TimeToPGDate(params.w.DatePerformed),
				InterSessionCntr: int16(params.w.Session),
				InterWorkoutCntr: int16(i + 1),
				Weight:           iterE.Weight,
				Sets:             iterE.Sets,
				Reps:             iterE.Reps,
				Effort:           iterE.Effort,
			},
		); opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotAddWorkoutErr, dal.FormatErr(opErr),
			)
			return
		}

		if len(iterE.BarPath) > 0 {
			state.PhysicsJobQueue.Schedule(&jobs.Physics{
				BarPath:              iterE.BarPath,
				Tl:                   params.bufWriter.Last(),
				B:                    params.batch,
				Q:                    queries,
				BarPathCalcParams:    params.barPathCalcParams,
				BarTrackerCalcParams: params.barTrackerCalcParams,
			})
		}
	}
	return
}

func validateWorkout(params createSingleWorkoutParams) (opErr error) {
	curExercise := -1
	curSet := -1
	wrapErr := func(err error, msg string, args ...any) error {
		msgStart := fmt.Sprintf("Exercise %d ", curExercise)
		if curSet >= 0 {
			msgStart = fmt.Sprintf("%s, Set %d ", msgStart, curSet)
		}
		return sberr.AppendError(
			sberr.Wrap(types.MalformedWorkoutExerciseErr, msgStart),
			sberr.Wrap(err, msg, args...),
		)
	}

	if params.w.Session <= 0 {
		opErr = sberr.Wrap(
			types.InvalidSessionErr,
			"Must be >0, Got: %d", params.w.Session,
		)
		return
	}

	for curExercise = range len(params.w.Exercises) {
		curSet = -1
		iterE := params.w.Exercises[curExercise]

		if len(iterE.BarPath) == 0 {
			// Not supplying any physics or bar path data is valid
			continue
		}

		ceilSets := int(math.Ceil(iterE.Sets))
		if len(iterE.BarPath) != ceilSets {
			opErr = wrapErr(
				types.InvalidBarPathsLenErr,
				"Number of sets %f -> %d, Length: %d",
				iterE.Sets, ceilSets, len(iterE.BarPath),
			)
			return
		}

		for _, curBarPath := range iterE.BarPath {
			if curBarPath.Source() == types.VideoBarPathData {
				videoPath, _ := curBarPath.VideoData()

				// These are surface level checks that make it so trivial errors
				// "fail fast" here rather than in a separate go routine in the
				// physics job queue.
				var fi fs.FileInfo
				if fi, opErr = os.Stat(videoPath); opErr != nil {
					return wrapErr(opErr, "")
				} else if fi.IsDir() {
					opErr = wrapErr(
						types.VideoPathDirNotFileErr,
						"Path: %s", videoPath,
					)
					return
				}
				// TODO - check video size/len limits
			} else if curBarPath.Source() == types.TimeSeriesBarPathData {
				timeSeriesData, _ := curBarPath.TimeSeriesData()
				lenTimeData := len(timeSeriesData.TimeData)
				lenPosData := len(timeSeriesData.PositionData)

				// These are surface level checks that make it so trivial errors
				// "fail fast" here rather than in a separate go routine in the
				// physics job queue.
				if lenTimeData != lenPosData {
					opErr = wrapErr(
						types.TimePositionDataMismatchErr,
						"len time data: %d, len pos data: %d",
						lenTimeData, lenPosData,
					)
					return
				}
				if lenTimeData < int(params.barPathCalcParams.MinNumSamples) {
					opErr = wrapErr(
						types.TimeDataLenErr,
						"minimum num samples: %d, got: %d",
						params.barPathCalcParams.MinNumSamples, lenTimeData,
					)
					return
				}
			}
		}
	}
	return
}

func CreateWorkoutsFromCSV(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	params := []types.RawWorkout{}
	opts.ReuseRecord = true

	for _, file := range files {
		clientName := strings.TrimSuffix(path.Base(file), path.Ext(file))

		if opErr = sbcsv.LoadCSVFile(file, &sbcsv.LoadOpts{
			Opts:          opts,
			RequestedCols: sbcsv.ReqColsForStruct[types.RawWorkout](),
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

				variants, err := parseDataDir(rawData.DataDir)
				if err != nil {
					return err
				}

				iterID := types.WorkoutID{
					ClientEmail:   clientName,
					Session:       rawData.Session,
					DatePerformed: rawData.DatePerformed,
				}
				if len(params) == 0 || params[len(params)-1].WorkoutID == iterID {
					params = append(params, types.RawWorkout{WorkoutID: iterID})
				}
				params[len(params)-1].Exercises = append(
					params[len(params)-1].Exercises,
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
			return opErr
		}
	}

	return CreateWorkouts(
		ctxt, state, queries,
		barPathCalcParams, barTrackerCalcParams,
		params...,
	)
}

func parseDataDir(dir string) (res []types.BarPathVariant, err error) {
	if dir == "" {
		return
	}
	var fi fs.FileInfo
	if fi, err = os.Stat(dir); err != nil {
		return
	} else if !fi.IsDir() {
		err = sberr.Wrap(
			types.InvalidDataDirErr, "'%s' was not a dir but must be", dir,
		)
		return
	}

	if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), "Set") {
			return nil
		}
		fmt.Println(path)
		return nil
	}); err != nil {
		err = sberr.AppendError(types.InvalidDataDirErr, err)
		return
	}

	return
}

func ReadClientTotalNumTrainingLogEntries(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = dal.Query1x2(
		dal.Q.GetTotalNumTrainingLogEntriesForClient, queries, ctxt, clientEmail,
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetTotalNumExercisesErr, dal.FormatErr(opErr),
		)
		return
	}
	state.Log.Log(ctxt, sblog.VLevel(3), "Read total num exercises for client")
	return
}

func ReadClientTotalNumPhysEntries(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = dal.Query1x2(
		dal.Q.GetTotalNumPhysicsEntriesForClient, queries, ctxt, clientEmail,
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetTotalNumPhysEntriesErr, dal.FormatErr(opErr),
		)
		return
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read total num phys entries for client",
	)
	return
}

func ReadClientNumWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = dal.Query1x2(
		dal.Q.GetNumWorkoutsForClient, queries, ctxt, clientEmail,
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetNumWorkoutsErr, dal.FormatErr(opErr),
		)
		return
	}
	state.Log.Log(ctxt, sblog.VLevel(3), "Read num workouts for client")
	return
}

func ReadWorkoutsByID(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	ids ...types.WorkoutID,
) (res []types.Workout, opErr error) {
	res = make([]types.Workout, len(ids))

	for i, id := range ids {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var rawData []dal.GetAllWorkoutDataRow
		rawData, opErr = dal.Query1x2(
			dal.Q.GetAllWorkoutData, queries, ctxt,
			dal.GetAllWorkoutDataParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int16(id.Session),
				DatePerformed:    dal.TimeToPGDate(id.DatePerformed),
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedWorkoutErr, dal.FormatErr(opErr),
			)
			return
		}
		if len(rawData) == 0 {
			opErr = sberr.Wrap(
				types.CouldNotFindRequestedWorkoutErr,
				"No data found for %+v", id,
			)
			return
		}
		res[i].WorkoutID = id
		res[i].Exercises = make([]types.ExerciseData, len(rawData))
		_ = dal.GetAllWorkoutDataRow(types.ExerciseData{})
		copy(
			res[i].Exercises,
			*(*[]types.ExerciseData)(unsafe.Pointer(&rawData)),
		)
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Read workouts by ID", "Num", len(ids))
	return
}

func FindWorkoutsByID(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	ids ...types.WorkoutID,
) (res []types.Found[types.Workout], opErr error) {
	res = make([]types.Found[types.Workout], len(ids))

	for i, id := range ids {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var rawData []dal.GetAllWorkoutDataRow
		rawData, opErr = dal.Query1x2(
			dal.Q.GetAllWorkoutData, queries, ctxt,
			dal.GetAllWorkoutDataParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int16(id.Session),
				DatePerformed:    dal.TimeToPGDate(id.DatePerformed),
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedWorkoutErr, dal.FormatErr(opErr),
			)
			return
		}
		if len(rawData) == 0 {
			res[i].Found = false
			continue
		}
		res[i].Found = true
		res[i].Value.WorkoutID = id
		res[i].Value.Exercises = make([]types.ExerciseData, len(rawData))
		_ = dal.GetAllWorkoutDataRow(types.ExerciseData{})
		copy(
			res[i].Value.Exercises,
			*(*[]types.ExerciseData)(unsafe.Pointer(&rawData)),
		)
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Found workouts by ID", "Num", len(ids))
	return
}

func ReadWorkoutsInDateRange(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clientEmail string,
	start time.Time,
	end time.Time,
) (res []types.Workout, opErr error) {
	var rawData []dal.GetAllWorkoutDataBetweenDatesRow

	if start.After(end) {
		opErr = sberr.Wrap(
			types.CouldNotFindRequestedWorkoutErr,
			"Start date (%s) must be after end date (%s)",
			start, end,
		)
		return
	}

	var ok bool
	ok, opErr = dal.Query1x2(dal.Q.ClientExists, queries, ctxt, clientEmail)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotFindRequestedWorkoutErr, dal.FormatErr(opErr),
		)
		return
	} else if !ok {
		opErr = sberr.AppendError(
			types.CouldNotFindRequestedWorkoutErr,
			sberr.Wrap(
				types.CouldNotFindRequestedClientErr, "Client: %s", clientEmail,
			),
		)
		return
	}

	rawData, opErr = dal.Query1x2(
		dal.Q.GetAllWorkoutDataBetweenDates, queries, ctxt,
		dal.GetAllWorkoutDataBetweenDatesParams{
			Email:  clientEmail,
			Start:  dal.TimeToPGDate(start),
			Ending: dal.TimeToPGDate(end),
		},
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotFindRequestedWorkoutErr, dal.FormatErr(opErr),
		)
		return
	}
	if len(rawData) == 0 {
		opErr = sberr.Wrap(
			types.CouldNotFindRequestedWorkoutErr,
			"No data found for date range: [%s, %s]", start, end,
		)
		return
	}

	res = make([]types.Workout, 0, 10)
	for i := range len(rawData) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		iterID := types.WorkoutID{
			ClientEmail:   clientEmail,
			Session:       uint16(rawData[i].InterSessionCntr),
			DatePerformed: rawData[i].DatePerformed.Time,
		}
		if len(res) == 0 || res[len(res)-1].WorkoutID != iterID {
			res = append(res, types.Workout{WorkoutID: iterID})
		}
		res[len(res)-1].Exercises = append(
			res[len(res)-1].Exercises,
			types.ExerciseData{
				Name:         rawData[i].Name,
				Weight:       rawData[i].Weight,
				Sets:         rawData[i].Sets,
				Reps:         rawData[i].Reps,
				Effort:       rawData[i].Effort,
				Volume:       rawData[i].Volume,
				Exertion:     rawData[i].Exertion,
				TotalReps:    rawData[i].TotalReps,
				Time:         rawData[i].Time,
				Position:     rawData[i].Position,
				Velocity:     rawData[i].Velocity,
				Acceleration: rawData[i].Acceleration,
				Jerk:         rawData[i].Jerk,
				Force:        rawData[i].Force,
				Impulse:      rawData[i].Impulse,
				Work:         rawData[i].Work,
				Power:        rawData[i].Power,
				RepSplits:    rawData[i].RepSplits,
			},
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read workouts from client with date range",
		"Num", len(res),
	)

	return
}

func UpdateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	data ...types.RawWorkout,
) (opErr error) {
	// TODO
	// Does the workout already exist??
	// What about upserting? That could get weird when the entire workout is
	// being replaced...
	// Delete then re-add? Would work well for checking if the workout exists
	return
}

func DeleteWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	ids ...types.WorkoutID,
) (opErr error) {
	for _, id := range ids {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var count int64
		count, opErr = dal.Query1x2(
			dal.Q.DeleteWorkout, queries, ctxt,
			dal.DeleteWorkoutParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int16(id.Session),
				DatePerformed:    dal.TimeToPGDate(id.DatePerformed),
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotDeleteRequestedWorkoutErr, dal.FormatErr(opErr),
			)
			return
		}
		if count == 0 {
			opErr = sberr.Wrap(
				types.CouldNotFindRequestedWorkoutErr,
				"No data found for %+v", id,
			)
			return
		}
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Deleted workouts", "Num", len(ids))
	return
}

func DeleteWorkoutsInDateRange(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clientEmail string,
	start time.Time,
	end time.Time,
) (res int64, opErr error) {
	if start.After(end) {
		opErr = sberr.Wrap(
			types.CouldNotFindRequestedWorkoutErr,
			"Start date (%s) must be after end date (%s)",
			start, end,
		)
		return
	}

	var ok bool
	ok, opErr = dal.Query1x2(dal.Q.ClientExists, queries, ctxt, clientEmail)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotFindRequestedWorkoutErr, dal.FormatErr(opErr),
		)
		return
	} else if !ok {
		opErr = sberr.AppendError(
			types.CouldNotFindRequestedWorkoutErr,
			sberr.Wrap(
				types.CouldNotFindRequestedClientErr, "Client: %s", clientEmail,
			),
		)
		return
	}

	res, opErr = dal.Query1x2(
		dal.Q.DeleteWorkoutsBetweenDates, queries, ctxt,
		dal.DeleteWorkoutsBetweenDatesParams{
			Email:  clientEmail,
			Start:  dal.TimeToPGDate(start),
			Ending: dal.TimeToPGDate(end),
		},
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotDeleteRequestedWorkoutErr, dal.FormatErr(opErr),
		)
		return
	}
	if res == 0 {
		opErr = sberr.Wrap(
			types.CouldNotFindRequestedWorkoutErr,
			"No data found for date range: [%s, %s]", start, end,
		)
		return
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Deleted workouts", "Num", res)
	return
}
