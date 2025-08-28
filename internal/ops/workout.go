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
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.RawWorkout,
) (opErr error) {
	syncQueries := dal.NewSyncQueries(queries)
	batch, _ := sbjobqueue.BatchWithContext(ctxt)
	clientCacheLoader := dal.NewClientCacheLoader(syncQueries)
	exerciseCacheLoader := dal.NewExerciseCacheLoader(syncQueries)
	bufWriter := dal.NewBufferedWriter(
		state.Global.BatchSize,
		func(
			ctxt context.Context,
			arg []dal.BulkCreateTrainingLogsParams,
		) (count int64, err error) {
			if err := batch.Wait(); err != nil {
				return 0, err
			}
			syncQueries.Run(func(q *dal.Queries) {
				count, err = q.BulkCreateTrainingLogs(ctxt, arg)
			})
			return
		},
	)

	for _, iterW := range data {
		// TODO - check if supplied ctxt was canceled; if it was break; look into
		// similar things in other ops

		if opErr = validateWorkout(state, &iterW); opErr != nil {
			opErr = sberr.AppendError(types.InvalidWorkoutErr, opErr)
			return
		}

		var iterClient types.IdWrapper[int64, types.Client]
		if iterClient, opErr = state.ClientCache.Get(
			ctxt, iterW.ClientEmail, clientCacheLoader,
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
			var iterExercise types.IdWrapper[int32, types.Exercise]
			if iterExercise, opErr = state.ExerciseCache.Get(
				ctxt, iterE.Name, exerciseCacheLoader,
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

			if opErr = bufWriter.Write(ctxt, dal.BulkCreateTrainingLogsParams{
				ClientID:   iterClient.Id,
				ExerciseID: iterExercise.Id,

				DatePerformed: pgtype.Date{
					Time:             iterW.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
				InterSessionCntr: int16(iterW.Session),
				InterWorkoutCntr: int16(i + 1),

				Weight: iterE.Weight,
				Sets:   iterE.Sets,
				Reps:   iterE.Reps,
				Effort: iterE.Effort,
			}); opErr != nil {
				opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
				return
			}

			if len(iterE.BarPath) > 0 {
				state.PhysicsJobQueue.Schedule(&jobs.Physics{
					BarPath: iterE.BarPath,
					Tl:      bufWriter.Last(),
					B:       batch,
					S:       state,
					Q:       syncQueries,
				})
			}
		}
	}

	if opErr = bufWriter.Flush(ctxt); opErr != nil {
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

func validateWorkout(state *types.State, w *types.RawWorkout) (opErr error) {
	curExercise := -1
	curSet := -1
	wrapErr := func(msg string, args ...any) error {
		msgStart := fmt.Sprintf("Exercise %d: ", curExercise)
		if curSet >= 0 {
			msgStart = fmt.Sprintf("%sSet %d: ", msgStart, curSet)
		}
		return sberr.Wrap(types.MalformedWorkoutExerciseErr, msgStart+msg, args...)
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
				"the bar paths list must either be empty or the same length as the ceiling of the number of sets (%f -> %d), Got: %d",
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
					return sberr.AppendError(
						types.MalformedWorkoutExerciseErr, opErr,
					)
				} else if fs.IsDir() {
					opErr = wrapErr(
						"expected a video file, got dir: %s",
						videoPath,
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
					opErr = wrapErr(fmt.Sprintf(
						"the length of the time data (%d) and position data (%d) must match",
						lenTimeData, lenPosData,
					))
					return
				}
				if lenTimeData < int(state.PhysicsData.MinNumSamples) {
					opErr = wrapErr(
						"the minimum number of samples (%d) was not provided, got %d samples",
						state.PhysicsData.MinNumSamples, lenTimeData,
					)
					return
				}
			}
		}
	}
	return
}

func ReadClientTotalNumExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = queries.GetTotalNumExercisesForClient(ctxt, clientEmail)
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
	queries *dal.Queries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = queries.GetTotalNumPhysicsEntriesForClient(ctxt, clientEmail)
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
	queries *dal.Queries,
	clientEmail string,
) (res int64, opErr error) {
	res, opErr = queries.GetNumWorkoutsForClient(ctxt, clientEmail)
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
	queries *dal.Queries,
	ids ...types.WorkoutID,
) (res []types.Workout, opErr error) {
	res = make([]types.Workout, len(ids))

	for i, id := range ids {
		var rawData []dal.GetAllWorkoutDataRow
		rawData, opErr = queries.GetAllWorkoutData(
			ctxt,
			dal.GetAllWorkoutDataParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int16(id.Session),
				DatePerformed: pgtype.Date{
					Time:             id.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
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
	queries *dal.Queries,
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
	if ok, opErr = queries.ClientExists(ctxt, clientEmail); opErr != nil {
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

	rawData, opErr = queries.GetAllWorkoutDataBetweenDates(
		ctxt,
		dal.GetAllWorkoutDataBetweenDatesParams{
			Email: clientEmail,
			Start: pgtype.Date{
				Time:             start,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			Ending: pgtype.Date{
				Time:             end,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
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
	queries *dal.Queries,
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
	queries *dal.Queries,
	ids ...types.WorkoutID,
) (opErr error) {
	for _, id := range ids {
		var count int64
		count, opErr = queries.DeleteWorkout(
			ctxt,
			dal.DeleteWorkoutParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int16(id.Session),
				DatePerformed: pgtype.Date{
					Time:             id.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
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
	queries *dal.Queries,
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
	if ok, opErr = queries.ClientExists(ctxt, clientEmail); opErr != nil {
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

	res, opErr = queries.DeleteWorkoutsBetweenDates(
		ctxt,
		dal.DeleteWorkoutsBetweenDatesParams{
			Email: clientEmail,
			Start: pgtype.Date{
				Time:             start,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
			Ending: pgtype.Date{
				Time:             end,
				InfinityModifier: pgtype.Finite,
				Valid:            true,
			},
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
