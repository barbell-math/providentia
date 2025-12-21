package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
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
// If any error occurs no changes will be made to the database.
func CreateHyperparams[T types.Hyperparams](
	ctxt context.Context,
	params ...T,
) (opErr error) {
	if len(params) == 0 {
		return
	}
	return runOp(ctxt, dal.CreateHyperparams[T], params)
}

// Checks that the supplied hyperparams are present in the database and adds
// them if they are not present. In order for the supplied hyperparams to be be
// considered already present the model type, version number, and parameter
// fields must all match. Any newly created hyperparams must satisfy the
// uniqueness constraints outlined by [CreateHyperparams].
//
// This function will be slower than [CreateHyperparams], so if you are working
// with large amounts of data and are ok with erroring on duplicated hyperparams
// consider using [CreateHyperparams].
//
// The context must have a [types.State] variable.
//
// Hyperparams will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func EnsureHyperparamsExist[T types.Hyperparams](
	ctxt context.Context,
	params ...T,
) (opErr error) {
	if len(params) == 0 {
		return
	}
	return runOp(ctxt, dal.EnsureHyperparamsExist[T], params)
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
	opts *sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, jobs.RunCSVLoaderJobs, jobs.CSVLoaderOpts[T]{
		Opts:    opts,
		Files:   files,
		Creator: dal.CreateHyperparams[T],
	})
}

// Checks that the supplied hyperparams are present in the database and adds
// them if they are not present. In order for the supplied hyperparams to be be
// considered already present the model type, version number, and parameter
// fields must all match. Any newly created hyperparams must satisfy the
// uniqueness constraints outlined by [CreateHyperparams]. All csv files must be
// valid as outlined by [CreateHyperparamsFromCSV].
//
// This function will be slower than [CreateHyperparams], so if you are working
// with large amounts of data and are ok with erroring on duplicated hyperparams
// consider using [CreateHyperparams].
//
// The context must have a [types.State] variable.
//
// Hyperparams will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func EnsureHyperparamsExistFromCSV[T types.Hyperparams](
	ctxt context.Context,
	opts *sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, jobs.RunCSVLoaderJobs, jobs.CSVLoaderOpts[T]{
		Opts:    opts,
		Files:   files,
		Creator: dal.EnsureHyperparamsExist[T],
	})
}

// Gets the total number of hyperparameters across all hyperparameter types in
// the database.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumHyperparams(ctxt context.Context) (res int64, opErr error) {
	opErr = runOp(ctxt, dal.ReadNumHyperparams, &res)
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
	opErr = runOp(ctxt, dal.ReadNumHyperparamsFor[T], &res)
	return
}

// Gets the hyperparameters associated with the supplied hyperparam type and
// versions if they exist. If they do not exist an error will be returned. The
// order of the returned hyperparams will match the order of the supplied
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
	opErr = runOp(
		ctxt, dal.ReadHyperparamsByVersionFor[T],
		dal.ReadHyperparamsByVersionForOpts[T]{
			Versions: versions,
			Params:   &res,
		},
	)
	return
}

// Gets the default hyperparameters associated with the supplied hyperparam type
// if they exist. If they do not exist an error will be returned. If the
// returned result is not consistent with what the providentia library expects
// an error will be returned.
//
// The context must have a [types.State].variable.
//
// No changes will be made to the database.
func ReadDefaultHyperparamsFor[T types.Hyperparams](
	ctxt context.Context,
) (res T, opErr error) {
	opErr = runOp(ctxt, dal.ReadDefaultHyperparamsFor[T], &res)
	return
}

// Gets the hyperparam data associated with the supplied versions for the
// provided type if they exist. If a hyperparam exists it will be put in the
// returned slice and the found flag will be set to true. If a hyperparam does
// not exist the value in the slice will be a zero initialized hyperparam and
// the found flag will be set to false. No error will be returned if a
// hyperparam does not exist. The order of the returned hyperparams will match
// the order of the supplied hyperparam versions.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func FindHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	versions ...int32,
) (res []types.Found[T], opErr error) {
	if len(versions) == 0 {
		return
	}
	opErr = runOp(
		ctxt, dal.FindHyperparamsByVersionFor[T],
		dal.FindHyperparamsByVersionForOpts[T]{
			Versions: versions,
			Params:   &res,
		},
	)
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
	return runOp(ctxt, dal.DeleteHyperparams[T], versions)
}
