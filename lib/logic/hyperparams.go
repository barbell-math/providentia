package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
)

// Adds the supplied hyperparameters to the database. The supplied hyperparams
// will be validated according to the rules for each type of parameters.
//
// For [types.BarPathCalcHyperparams] the following must be true:
//   - MinNumSamples >= 2
//   - TimeDeltaEps > 0
//   - ApproxErr must be a valid approx error enum value
//   - NearZeroFilter > 0
//
// For [types.BarPathTrackerHyperparams] the following must be true:
//   - MinLength > 0
//   - MaxFileSize > MinFileSize
//
// The pairing of the type of hyperparameter and version number must be unique,
// including the set of pairs of hyperparmeter type and version number already
// in the database.
//
// The context must have a [types.State] variable.
//
// Hyperparameters will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateHyperparams[T types.Hyperparams](
	ctxt context.Context,
	params ...T,
) (opErr error) {
	if len(params) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.CreateHyperparams(ctxt, state, queries, params...)
		},
	})
}

// Adds the hyperparams supplied in the csv files to the database. Has the same
// behavior as [CreateHyperparams] other than getting the clients from csv files.
// The csv files are expected to have column names on the first row. The fields
// of each hyperparams struct are used to determine the required columns, the
// column names, and types for each column.
//
// The `ReuseRecord` field on opts will be set to true before loading the csv
// file. All other options are left alone.
//
// The context must have a [types.State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateHyperparamsFromCSV[T types.Hyperparams](
	ctxt context.Context,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.CreateHyperparamsFromCSV[T](
				ctxt, state, queries, opts, files...,
			)
		},
	})
}

// Gets the total number of hyperparameters across all hyperparameter types in
// the database.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumHyperparams(ctxt context.Context) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadNumHyperparams(ctxt, state, queries)
			return err
		},
	})
	return
}

// Gets the total number of hyperparameters for the given hyperparameter type in
// the database.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumHyperparamsFor[T types.Hyperparams](
	ctxt context.Context,
) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadNumHyperparamsFor[T](ctxt, state, queries)
			return err
		},
	})
	return
}

// Gets the hyperparameters associated with the supplied hyperparam type and
// versions if they exist. If they do not exist an error will be returned. The
// order of the returned hyperparams may not match the order of the supplied
// hyperparams.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	versions ...int32,
) (res []T, opErr error) {
	if len(versions) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadHyperparamsByVersionFor[T](
				ctxt, state, queries, versions...,
			)
			return err
		},
	})
	return
}

// Deletes the supplied hyperparameters, as identified by their hyperparameter
// type and version number. All data associated with the hyperparameters will be
// deleted.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func DeleteHyperparams[T types.Hyperparams](
	ctxt context.Context,
	versions ...int32,
) (opErr error) {
	if len(versions) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.DeleteHyperparams[T](ctxt, state, queries, versions...)
		},
	})
}
