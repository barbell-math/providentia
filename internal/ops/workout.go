package ops

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"

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

	// TODO - make caches separate - part of init state??
	clientIdCache := map[string]int64{}
	exerciseCache := map[string]int32{}

	bufWriter := NewBufferedWriter[dal.BulkCreateTrainingLogsParams](
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
		if iterW.Session <= 0 {
			opErr = sberr.AppendError(
				types.InvalidWorkoutErr,
				sberr.Wrap(
					types.InvalidSessionErr,
					"Must be >0, Got: %d", iterW.Session,
				),
			)
			return
		}
		if _, ok := clientIdCache[iterW.ClientEmail]; !ok {
			var clientID int64
			syncQueries.Run(func(q *dal.Queries) {
				clientID, opErr = q.GetClientIDFromEmail(
					ctxt, iterW.ClientEmail,
				)
			})
			if opErr != nil {
				opErr = sberr.AppendError(
					types.InvalidWorkoutErr,
					sberr.AppendError(
						types.CouldNotFindRequestedClientErr, opErr,
					),
				)
				return
			}
			clientIdCache[iterW.ClientEmail] = clientID
		}

		if opErr = validateWorkout(
			ctxt, state, syncQueries,
			&iterW, exerciseCache,
		); opErr != nil {
			opErr = sberr.AppendError(types.InvalidWorkoutErr, opErr)
			return
		}

		for i, iterE := range iterW.Exercises {
			if opErr = bufWriter.Write(ctxt, dal.BulkCreateTrainingLogsParams{
				ClientID: clientIdCache[iterW.ClientEmail],
				// exerciseCache is populated by validateWorkout
				ExerciseID: exerciseCache[iterE.Name],

				DatePerformed: pgtype.Date{
					Time:             iterW.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
				InterSessionCntr: int32(iterW.Session),
				InterWorkoutCntr: int32(i + 1),

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

func validateWorkout(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	w *types.RawWorkout,
	exerciseCache map[string]int32,
) (opErr error) {
	curExercise := -1
	curSet := -1
	wrapErr := func(msg string, args ...any) error {
		msgStart := fmt.Sprintf("Exercise %d: ", curExercise)
		if curSet >= 0 {
			msgStart = fmt.Sprintf("%sSet %d: ", msgStart, curSet)
		}
		return sberr.Wrap(types.MalformedWorkoutExerciseErr, msgStart+msg, args...)
	}

	for curExercise = range len(w.Exercises) {
		curSet = -1
		iterE := w.Exercises[curExercise]

		if iterId, ok := exerciseCache[iterE.Name]; !ok {
			queries.Run(func(q *dal.Queries) {
				iterId, opErr = q.GetExerciseId(ctxt, iterE.Name)
			})
			if opErr != nil {
				opErr = sberr.AppendError(
					wrapErr(""),
					sberr.Wrap(
						types.CouldNotFindRequestedExerciseErr,
						"Unknown exercise: %s", iterE.Name,
					),
				)
				return
			}
			exerciseCache[iterE.Name] = iterId
		}

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
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read total num exercises for client",
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
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read num workouts for client",
	)
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
		// TODO - make sure exercises are returned in the correct order!!!
		rawData, opErr = queries.GetAllWorkoutData(
			ctxt,
			dal.GetAllWorkoutDataParams{
				Email:            id.ClientEmail,
				InterSessionCntr: int32(id.Session),
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
		res[i].BasicData = make([]types.BasicData, len(rawData))
		res[i].PhysData = make([]types.PhysicsData, len(rawData))
		for j := 0; j < len(rawData); j++ {
			res[i].BasicData[j] = types.BasicData{
				Name:      rawData[j].Name,
				Weight:    rawData[j].Weight,
				Sets:      rawData[j].Sets,
				Reps:      rawData[j].Reps,
				Effort:    rawData[j].Effort,
				Volume:    rawData[j].Volume,
				Exertion:  rawData[j].Exertion,
				TotalReps: rawData[j].TotalReps,
			}
			res[i].PhysData[j] = types.PhysicsData{
				Time:         rawData[j].Time,
				Position:     rawData[j].Position,
				Velocity:     rawData[j].Velocity,
				Acceleration: rawData[j].Acceleration,
				Jerk:         rawData[j].Jerk,
				Force:        rawData[j].Force,
				Impulse:      rawData[j].Impulse,
				Work:         rawData[j].Work,
			}
		}
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read workouts from client",
		"Num", len(ids),
	)

	return
}

func UpdateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.RawWorkout,
) (opErr error) {
	return
}

// // Eh?? -might make for quicker physics data updates but would require another
// // user facing api type UpdatePhysData struct { WorkoutID, PhysData }
// func UpdateWorkoutPhysicsData(
// 	ctxt context.Context,
// 	state *types.State,
// 	queries *dal.Queries,
// 	data ...types.RawWorkout,
// ) (opErr error) {
// 	return
// }

func DeleteWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.RawWorkout,
) (opErr error) {
	return
}
