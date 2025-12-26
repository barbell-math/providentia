package tests

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestBulkUpload(t *testing.T) {
	t.Run("failingNoWrites", bulkUploadFailingNoWrites)
	t.Run("passing", bulkUploadPassing)
}

func bulkUploadFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("badClientDir", bulkUploadBadClientDir(ctxt))
	t.Run("badClientFile", bulkUploadBadClientFile(ctxt))
	t.Run("malformedClientFile", bulkUploadMalformedClientFile(ctxt))
	t.Run("badExerciseDir", bulkUploadBadExerciseDir(ctxt))
	t.Run("badExerciseFile", bulkUploadBadExerciseFile(ctxt))
	t.Run("malformedExerciseFile", bulkUploadMalformedExerciseFile(ctxt))
	t.Run("badHyperparamsDir", bulkUploadBadHyperparamsDir(ctxt))
	t.Run("badHyperparamsFile", bulkUploadBadHyperparamsFile(ctxt))
	t.Run("badHyperparamsExt", bulkUploadBadHyperparamsExt(ctxt))
	t.Run("badHyperparamsType", bulkUploadBadHyperparamsType(ctxt))
	t.Run("malformedHyperparamFile", bulkUploadMalformedHyperparamFile(ctxt))
}

func bulkUploadBadClientDir(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ClientDir:        "./testData/dataBadDir",
			ClientCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`open ./testData/dataBadDir: no such file or directory`,
		)
	}
}

func bulkUploadBadClientFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ClientDir:        "./testData/dataBadFile",
			ClientCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`Supplied dir \(\./testData/dataBadFile\) contained non-csv files`,
		)
	}
}

func bulkUploadMalformedClientFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ClientDir:        "./testData/badClientData",
			ClientCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
		sbtest.ContainsError(t, sbcsv.InvalidCSVFileErr, err)
		sbtest.ContainsError(
			t, sbcsv.MissingColumnErr, err,
			`Requested column 'LastName' was not found`,
		)
	}
}

func bulkUploadBadExerciseDir(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ExerciseDir:        "./testData/dataBadDir",
			ExerciseCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`open ./testData/dataBadDir: no such file or directory`,
		)
	}
}

func bulkUploadBadExerciseFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ExerciseDir:        "./testData/dataBadFile",
			ExerciseCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`Supplied dir \(\./testData/dataBadFile\) contained non-csv files`,
		)
	}
}

func bulkUploadMalformedExerciseFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ExerciseDir:        "./testData/badExerciseData",
			ExerciseCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
		sbtest.ContainsError(t, sbcsv.InvalidCSVFileErr, err)
		sbtest.ContainsError(
			t, sbcsv.MissingColumnErr, err,
			`Requested column 'KindId' was not found`,
		)
	}
}

func bulkUploadBadHyperparamsDir(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			HyperparamsDir:        "./testData/dataBadDir",
			HyperparamsCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`open ./testData/dataBadDir: no such file or directory`,
		)
	}
}

func bulkUploadBadHyperparamsFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			HyperparamsDir:        "./testData/dataBadFile",
			HyperparamsCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.FilesWithExtInDirErr, err,
			`Supplied dir \(\./testData/dataBadFile\) contained non-csv files`,
		)
	}
}

func bulkUploadBadHyperparamsExt(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			HyperparamsDir:        "./testData/badHyperparamExtData",
			HyperparamsCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, types.UnknownFileInDataDirErr, err,
			`Invalid hyperparam file. File name must have the following format: <file name>.<hyperparam type>.csv`,
		)
	}
}

func bulkUploadBadHyperparamsType(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			HyperparamsDir:        "./testData/badHyperparamTypeData",
			HyperparamsCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, types.UnknownFileInDataDirErr, err,
			`Unknown hyperparam type: .*, must be one of .*. File: .*`,
		)
	}
}

func bulkUploadMalformedHyperparamFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			HyperparamsDir:        "./testData/badHyperparamData",
			HyperparamsCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
		sbtest.ContainsError(t, sbcsv.InvalidCSVFileErr, err)
		sbtest.ContainsError(
			t, sbcsv.MissingColumnErr, err,
			`Requested column 'MinNumSamples' was not found`,
		)
	}
}

func bulkUploadPassing(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
		ClientDir:             "./testData/clientData",
		ClientCreateType:      types.Create,
		ExerciseDir:           "./testData/exerciseData",
		ExerciseCreateType:    types.Create,
		HyperparamsDir:        "./testData/hyperparamData",
		HyperparamsCreateType: types.Create,
		WorkoutDir:            "./testData/workoutData",
		BarPathCalcHyperparams: &types.BarPathCalcHyperparams{
			MinNumSamples: 5,
			ApproxErr:     types.SecondOrder,
			NoiseFilter:   1,
		},
		BarPathTrackerHyperparams: &types.BarPathTrackerHyperparams{},
		Opts:                      sbcsv.Opts{TimeFormat: "1/2/2006"},
	})
	sbtest.Nil(t, err)

	numClients, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numClients, 2)

	numExercises, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numExercises, int64(len(migrations.ExerciseSetupData)+3))

	numHyperparams, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, numDefaultHyperparams+4)

	numHyperparams, err = logic.ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, 3)
	numHyperparams, err = logic.ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, 3)

	numWorkouts, err := logic.ReadNumWorkoutsForClient(ctxt, "two@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, numWorkouts, 3)
}
