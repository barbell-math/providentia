package ops

import (
	"context"
	"os"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/barbell-math/providentia/lib/types"
	sberr "github.com/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	clientIdMap := map[string]int64{}
	exerciseMap := map[string]dal.GetExerciseIDsRow{}
	bufWriter := NewBufferedWriter[dal.BulkCreateTrainingLogParams](
		state.Conf.BatchSize, queries.BulkCreateTrainingLog,
	)

	for _, iterW := range data {
		if iterW.Session <= 0 {
			opErr = sberr.Wrap(
				types.InvalidSessionErr,
				"Must be >=0, Got: %d", iterW.Session,
			)
			return
		}
		for _, iterE := range iterW.Exercises {
			if iterE.VideoPath != "" {
				fi, err := os.Stat(iterE.VideoPath)
				if err != nil {
					opErr = sberr.AppendError(types.InvalidVideoFileErr, err)
					return
				} else if fi.IsDir() {
					opErr = sberr.Wrap(
						types.InvalidVideoFileErr,
						"Supplied path was a dir not a file",
					)
					return
				}
			}
		}

		if _, ok := clientIdMap[iterW.ClientEmail]; !ok {
			var clientID int64
			clientID, opErr = queries.GetClientIDFromEmail(ctxt, iterW.ClientEmail)
			if opErr != nil {
				opErr = sberr.AppendError(
					types.CouldNotFindRequestedClientErr, opErr,
				)
				return
			}
			clientIdMap[iterW.ClientEmail] = clientID
		}

		// TODO - check there is no existing workout

		ids := make([]dal.GetExerciseIDsRow, len(iterW.Exercises))
		for i, iterE := range iterW.Exercises {
			if iterIds, ok := exerciseMap[iterE.Name]; ok {
				ids[i] = iterIds
				continue
			}
			var iterIds dal.GetExerciseIDsRow
			iterIds, opErr = queries.GetExerciseIDs(ctxt, iterE.Name)
			if opErr != nil {
				opErr = sberr.AppendError(
					sberr.Wrap(
						types.CouldNotFindRequestedExerciseErr,
						"Missing exercise: %s", iterE.Name,
					),
					opErr,
				)
				return
			}
			exerciseMap[iterE.Name] = iterIds
			ids[i] = iterIds
		}

		// TODO - precalc things like video id (which will require some arbitrary
		// computaiton) and store the resulting ids in a tlData index->videoid map

		for i, iterE := range iterW.Exercises {
			if opErr = bufWriter.Write(ctxt, dal.BulkCreateTrainingLogParams{
				ExerciseID:      ids[i].ExerciseID,
				ExerciseKindID:  ids[i].KindID,
				ExerciseFocusID: ids[i].FocusID,
				ClientID:        clientIdMap[iterW.ClientEmail],
				// TODO -Videoid ???

				DatePerformed: pgtype.Date{
					Time:             iterW.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
				Weight: iterE.Weight,
				Sets:   iterE.Sets,
				Reps:   iterE.Reps,
				Effort: iterE.Effort,

				InterSessionCntr: iterW.Session,
				InterWorkoutCntr: int32(i + 1),
			}); opErr != nil {
				opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
				return
			}
		}
	}

	if opErr = bufWriter.Flush(ctxt); opErr != nil {
		opErr = sberr.AppendError(types.CouldNotAddWorkoutErr, opErr)
		return
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

func DeleteWorkouts(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	data ...types.Workout,
) (opErr error) {
	return
}
