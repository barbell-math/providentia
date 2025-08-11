package ops

import (
	"context"
	"fmt"
	"math"
	"os"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	clientIdCache := map[string]int64{}
	exerciseCache := map[string]int32{}
	bufWriter := NewBufferedWriter[dal.BulkCreateTrainingLogParams](
		state.Global.BatchSize, queries.BulkCreateTrainingLog,
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

		var physicsData [][]types.PhysicsData
		var exerciseIds []int32
		_ = physicsData
		_ = exerciseIds
		// TODO - upload physics data...get ids...use ids when writing training
		// logs...

		// for i, iterE := range iterW.Exercises {
		// 	if opErr = bufWriter.Write(ctxt, dal.BulkCreateTrainingLogParams{
		// 		ExerciseID:      ids[i].ExerciseID,
		// 		ExerciseKindID:  ids[i].KindID,
		// 		ExerciseFocusID: ids[i].FocusID,
		// 		ClientID:        clientIdCache[iterW.ClientEmail],
		// 		// TODO -Videoid ???

		// 		DatePerformed: pgtype.Date{
		// 			Time:             iterW.DatePerformed,
		// 			InfinityModifier: pgtype.Finite,
		// 			Valid:            true,
		// 		},
		// 		Weight: iterE.Weight,
		// 		Sets:   iterE.Sets,
		// 		Reps:   iterE.Reps,
		// 		Effort: iterE.Effort,

		// 		InterSessionCntr: iterW.Session,
		// 		InterWorkoutCntr: int32(i + 1),
		// 	}); opErr != nil {
		// 		opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
		// 		return
		// 	}
		// }
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

		if len(iterE.VideoPaths) != len(iterE.TimeData) {
			opErr = wrapErr(
				"the length of the supplied time data (%d) and video paths (%d) must match",
				len(iterE.TimeData), len(iterE.VideoPaths),
			)
			return
		}
		if len(iterE.PositionData) != len(iterE.TimeData) {
			opErr = wrapErr(
				"the length of the supplied position data (%d) and time data (%d) must match",
				len(iterE.TimeData), len(iterE.PositionData),
			)
			return
		}
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

		for curSet := range len(iterE.VideoPaths) {
			if setErr := validateSet(
				ctxt, state, queries, &iterE, curSet,
			); setErr != "" {
				opErr = wrapErr(setErr)
				return
			}
		}
	}
	return
}

func validateSet(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	e *types.RawData,
	set int,
) (opErr string) {
	timeLen := len(e.TimeData[set])
	posLen := len(e.PositionData[set])

	if timeLen != 0 && posLen != 0 {
		if timeLen != posLen {
			opErr = fmt.Sprintf(
				"the length of the time data (%d) and position data (%d) must match",
				timeLen, posLen,
			)
			return
		}
		if timeLen < int(state.Physics.MinNumSamples) {
			opErr = fmt.Sprintf(
				"the minimum number of samples (%d) was not provided, got %d samples",
				state.Physics.MinNumSamples, timeLen,
			)
			return
		}
		delta := e.TimeData[set][1] - e.TimeData[set][0]
		for i := 1; i < timeLen; i++ {
			iterDelta := e.TimeData[set][i] - e.TimeData[set][i-1]
			if iterDelta < 0 {
				opErr = fmt.Sprintf(
					"time samples must be increasing, got a delta of %f",
					iterDelta,
				)
				return
			}
			if math.Abs(iterDelta-delta) < state.Physics.TimeDeltaEps {
				opErr = fmt.Sprintf(
					"time samples must all have the same delta (within %f variance), got delta of %f and %f",
					state.Physics.TimeDeltaEps, delta, iterDelta,
				)
			}
		}

		// TODO - calculate other physics values from raw data
		// set some cntr to wait for results later
	} else if e.VideoPaths[set] != "" {
		if fs, err := os.Stat(e.VideoPaths[set]); err != nil {
			opErr = err.Error()
			return
		} else if fs.IsDir() {
			opErr = fmt.Sprintf(
				"expected a video file, got dir: %s",
				e.VideoPaths[set],
			)
			return
		}
		// TODO - check video size limits
		// TODO - schedule video for processing in the queue and
		// set some cntr to wait for results later
		// Extract time and pos data only, use some kind of closure here to check
		// for min num samples, then pass on to same proc as above case
	}
	return
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
