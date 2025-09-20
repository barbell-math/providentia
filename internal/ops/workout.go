package ops

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"
	"time"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
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

	for _, iterW := range data {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		if opErr = validateWorkout(
			&iterW, barPathCalcParams, barTrackerCalcParams,
		); opErr != nil {
			opErr = sberr.AppendError(types.InvalidWorkoutErr, opErr)
			return
		}

		var iterClientId int64
		if iterClientId, opErr = clientCache.Get(
			ctxt, queries, iterW.ClientEmail,
		); opErr != nil {
			opErr = sberr.AppendError(
				types.InvalidWorkoutErr,
				sberr.Wrap(
					types.CouldNotFindRequestedClientErr,
					"Unknown Email: %s", iterW.ClientEmail,
				),
				opErr,
			)
			return
		}

		for i, iterE := range iterW.Exercises {
			var iterExerciseId int32
			if iterExerciseId, opErr = exerciseCache.Get(
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

			if opErr = bufWriter.Write(
				ctxt, queries,
				dal.BulkCreateTrainingLogsParams{
					ClientID:         iterClientId,
					ExerciseID:       iterExerciseId,
					DatePerformed:    dal.TimeToPGDate(iterW.DatePerformed),
					InterSessionCntr: int16(iterW.Session),
					InterWorkoutCntr: int16(i + 1),
					Weight:           iterE.Weight,
					Sets:             iterE.Sets,
					Reps:             iterE.Reps,
					Effort:           iterE.Effort,
				},
			); opErr != nil {
				opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
				return
			}

			if len(iterE.BarPath) > 0 {
				state.PhysicsJobQueue.Schedule(&jobs.Physics{
					BarPath:              iterE.BarPath,
					Tl:                   bufWriter.Last(),
					B:                    batch,
					Q:                    queries,
					BarPathCalcParams:    barPathCalcParams,
					BarTrackerCalcParams: barTrackerCalcParams,
				})
			}
		}
	}

	if opErr = bufWriter.Flush(ctxt, queries); opErr != nil {
		opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
		return
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Added new workouts",
		"NumWorkouts", len(data),
	)
	return
}

func validateWorkout(
	w *types.RawWorkout,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
) (opErr error) {
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

	if w.Session <= 0 {
		opErr = sberr.Wrap(
			types.InvalidSessionErr,
			"Must be >0, Got: %d", w.Session,
		)
		return
	}

	for curExercise = range len(w.Exercises) {
		curSet = -1
		iterE := w.Exercises[curExercise]

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
				var fs fs.FileInfo
				if fs, opErr = os.Stat(videoPath); opErr != nil {
					return wrapErr(opErr, "")
				} else if fs.IsDir() {
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
				if lenTimeData < int(barPathCalcParams.MinNumSamples) {
					opErr = wrapErr(
						types.TimeDataLenErr,
						"minimum num samples: %d, got: %d",
						barPathCalcParams.MinNumSamples, lenTimeData,
					)
					return
				}
			}
		}
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
		opErr = sberr.AppendError(types.CouldNotGetTotalNumExercisesErr, opErr)
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
		opErr = sberr.AppendError(types.CouldNotGetTotalNumPhysEntriesErr, opErr)
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
		opErr = sberr.AppendError(types.CouldNotGetNumWorkoutsErr, opErr)
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
				types.CouldNotFindRequestedWorkoutErr, opErr,
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

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read workouts from client",
		"Num", len(ids),
	)

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
		opErr = sberr.AppendError(types.CouldNotFindRequestedWorkoutErr, opErr)
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
			types.CouldNotFindRequestedWorkoutErr, opErr,
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
				types.CouldNotDeleteRequestedWorkoutErr, opErr,
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
		opErr = sberr.AppendError(types.CouldNotFindRequestedWorkoutErr, opErr)
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
			types.CouldNotDeleteRequestedWorkoutErr, opErr,
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
