package ops

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	batch, _ := sbjobqueue.BatchWithContext(ctxt)
	clientIdCache := map[string]int64{}
	exerciseCache := map[string]int32{}
	bufWriter := NewBufferedWriter[dal.BulkCreateTrainingLogParams](
		state.Global.BatchSize,
		func(
			ctxt context.Context,
			arg []dal.BulkCreateTrainingLogParams,
		) (int64, error) {
			if err := batch.Wait(); err != nil {
				return 0, err
			}
			return queries.BulkCreateTrainingLog(ctxt, arg)
		},
	)

	for _, iterW := range data {
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
			clientID, opErr = queries.GetClientIDFromEmail(ctxt, iterW.ClientEmail)
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
			ctxt, state, queries,
			&iterW, exerciseCache,
		); opErr != nil {
			opErr = sberr.AppendError(types.InvalidWorkoutErr, opErr)
			return
		}

		for i, iterE := range iterW.Exercises {
			if opErr = bufWriter.Write(ctxt, dal.BulkCreateTrainingLogParams{
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

			if rawDataHasPhysData(&iterE) {
				state.PhysicsJobQueue.Schedule(&physicsJob{
					BarPath: iterE.BarPath,
					Tl:      bufWriter.Last(),
					B:       batch,
					S:       state,
				})
			}
		}
	}

	if opErr = bufWriter.Flush(ctxt); opErr != nil {
		opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
		return
	}

	return
}

func validateWorkout(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	w *types.Workout,
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
			iterId, opErr = queries.GetExerciseId(ctxt, iterE.Name)
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

		if !rawDataHasPhysData(&iterE) {
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

func rawDataHasPhysData(e *types.RawData) bool {
	return len(e.BarPath) != 0
}

func ReadWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	id ...types.WorkoutID,
) (opErr error) {
	return
}

func ReadWorkout(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	id types.WorkoutID,
) (opErr error) {
	return
}

func ReadNumWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	clientEmail string,
) (opErr error) {
	return
}

// Intent is to read number of exercises a given client has performed
// func ReadNumExercises(
// 	ctxt context.Context,
// 	state *types.State,
// 	queries *dal.Queries,
// 	clientEmail string,
// ) (opErr error) {
// 	return
// }

func UpdateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	return
}

// Eh?? -might make for quicker physics data updates but would require another
// user facing api type UpdatePhysData struct { WorkoutID, PhysData }
func UpdateWorkoutPhysicsData(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	return
}

func DeleteWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	return
}
