package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// TODO - add comment
func CreateWorkouts(
	ctxt context.Context,
	workouts ...types.RawWorkout,
) (opErr error) {
	if len(workouts) == 0 {
		return
	}
	return runOp(ctxt, func(state *types.State, queries *dal.Queries) error {
		return ops.CreateWorkouts(ctxt, state, queries, workouts...)
	})
}

// Gets the total number of exercies across all workouts in the database for a
// given client.
//
// The context must have a [State] variable.
//
// No changes will be made to the database.
func ReadClientTotalNumExercises(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		res, err = ops.ReadClientTotalNumExercises(
			ctxt, state, queries, clientEmail,
		)
		return err
	})
	return
}

func ReadClientTotalNumPhysEntries(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		res, err = ops.ReadClientTotalNumPhysEntries(
			ctxt, state, queries, clientEmail,
		)
		return err
	})
	return
}

// Gets the total number of workouts in the database for a given client.
//
// The context must have a [State] variable.
//
// No changes will be made to the database.
func ReadClientNumWorkouts(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		res, err = ops.ReadClientNumWorkouts(ctxt, state, queries, clientEmail)
		return err
	})
	return
}

// Gets the workout data associated with the supplied ids if they exist. If they
// do not exist an error will be returned.
//
// The context must have a [State] variable.
//
// No changes will be made to the database.
func ReadWorkoutsByID(
	ctxt context.Context,
	ids ...types.WorkoutID,
) (res []types.Workout, opErr error) {
	opErr = runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		res, err = ops.ReadWorkoutsByID(ctxt, state, queries, ids...)
		return err
	})
	return
}
