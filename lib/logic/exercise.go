package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
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

// Checks that the supplied exercises are present in the database and adds them
// if they are not present. In order for the supplied exercises to be be
// considered already present the name, kind, and focus fields must all match.
// Any newly created exercises must satisfy the uniqueness constraints outlined
// by [CreateExercises].
//
// This function will be slower than [CreateExercises], so if you are working
// with large amounts of data and are ok with erroring on duplicated exercises
// consider using [CreateExercises].
//
// The context must have a [types.State] variable.
//
// Exercises will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func EnsureExercisesExist(
	ctxt context.Context,
	exercises ...types.Exercise,
) (opErr error) {
	if len(exercises) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.EnsureExercisesExist(ctxt, state, queries, exercises...)
		},
	})
}

// Adds the exercises supplied in the csv files to the database. Has the same
// behavior as [CreateExercises] other than getting the exercises from csv
// files. The csv files are expected to have column names on the first row and
// the following columns must be present as identified by the column name on the
// first row. More columns may be present, they will be ignored.
//   - Name (string): the name of the exercise
//   - KindID (string): one of MainCompound, MainCompoundAccessory,
//     CompoundAccessory, or Accessory
//   - FocusID (string): one of UnknownExerciseFocus, Squat, Bench, or Deadlift
//
// The `ReuseRecord` field on opts will be set to true before loading the csv
// file. All other options are left alone.
//
// The context must have a [types.State] variable.
//
// Exercises will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateExercisesFromCSV(
	ctxt context.Context,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.UploadExercisesFromCSV(
				ctxt, state, queries, ops.CreateExercises, opts, files...,
			)
		},
	})
}

// TODO - doc, test
func EnsureExercisesExistFromCSV(
	ctxt context.Context,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.UploadExercisesFromCSV(
				ctxt, state, queries, ops.EnsureExercisesExist, opts, files...,
			)
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
// exercises will match the order of the supplied exercise names.
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
			return
		},
	})
	return
}

// Gets the exercise data associated with the supplied exercises if they exist.
// If a exercise exists it will be put in the returned slice and the found flag
// will be set to true. If a exercise does not exist the value in the slice will
// be a zero initialized exercise and the found flag will be set to false. No
// error will be returned if a exercise does not exist. The order of the
// returned exercises will match the order of the supplied exercise emails.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func FindExercisesByName(
	ctxt context.Context,
	names ...string,
) (res []types.Found[types.Exercise], opErr error) {
	if len(names) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.FindExercisesByName(ctxt, state, queries, names...)
			return
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
