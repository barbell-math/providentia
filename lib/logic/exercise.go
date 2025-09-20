package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Adds the supplied exercises to the database. The supplied name for each
// exercise must not be an empty string. The supplied id fields must map to
// valid enum values. Exercise names must not be duplicated, including the set
// of exercises that are already in the database.
//
// The context must have a [types.State] variable.
//
// Exercises will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateExercises(
	ctxt context.Context,
	exercises ...types.Exercise,
) (opErr error) {
	if len(exercises) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			return ops.CreateExercises(ctxt, state, queries, exercises...)
		},
	})
}

// Gets the total number of exercises in the database.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumExercises(ctxt context.Context) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadNumExercises(ctxt, state, queries)
			return err
		},
	})
	return
}

// Gets the exercise data associated with the supplied names if they exist. If
// they do not exist an error will be returned. The order of the returned
// exercises may not match the order of the supplied exercise names.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadExercisesByName(
	ctxt context.Context,
	names ...string,
) (res []types.Exercise, opErr error) {
	if len(names) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadExercisesByName(ctxt, state, queries, names...)
			return err
		},
	})
	return
}

// Updates the supplied exercises, as identified by their name, with the data
// from the supplied structs. Names cannot be updated due to their uniqueness
// constraint. If an exercise is supplied with a name that does not exist in the
// database an error will be returned.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func UpdateExercises(
	ctxt context.Context,
	exercises ...types.Exercise,
) (opErr error) {
	if len(exercises) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.UpdateExercises(ctxt, state, queries, exercises...)
		},
	})
}

// Deletes the supplied exercises, as identified by their name. All data
// associated with the exercise will be deleted.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func DeleteExercises(ctxt context.Context, names ...string) (opErr error) {
	if len(names) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.DeleteExercises(ctxt, state, queries, names...)
		},
	})
}
