package tests

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestBulkUpload(t *testing.T) {
	t.Run("failingNoWrites", bulkUploadFailingNoWrites)
}

func bulkUploadFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("badClientDir", bulkUploadBadClientDir(ctxt))
	t.Run("badClientFile", bulkUploadBadClientFile(ctxt))
	// TODO - malformed client file
	t.Run("badExerciseDir", bulkUploadBadExerciseDir(ctxt))
	t.Run("badExerciseFile", bulkUploadBadExerciseFile(ctxt))
	// TODO - malformed exercise file
}

func bulkUploadBadClientDir(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.BulkUploadData(ctxt, &types.BulkUploadDataOpts{
			ClientDir:        "./testData/dataBadDir",
			ClientCreateType: types.Create,
		})
		sbtest.ContainsError(t, types.BulkDataUploadErr, err)
		sbtest.ContainsError(
			t, util.GetCSVFilesInDirErr, err,
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
			t, util.GetCSVFilesInDirErr, err,
			`Supplied dir \(\./testData/dataBadFile\) contained non-csv files in strict mode`,
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
			t, util.GetCSVFilesInDirErr, err,
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
			t, util.GetCSVFilesInDirErr, err,
			`Supplied dir \(\./testData/dataBadFile\) contained non-csv files in strict mode`,
		)
	}
}

// func TestMiscTestUploadCSVDataDirHyperparamBadDir(t *testing.T) {
// 	ctxt, cleanup := resetApp(t, context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
// 		HyperparamsDir:        "./testData/dataBadDir",
// 		HyperparamsCreateType: types.Create,
// 	})
// 	sbtest.ContainsError(t, types.DataDirErr, err)
// 	sbtest.ContainsError(t, types.DirInDataDirErr, err)
// }
//
// func TestMiscTestUploadCSVDataDirHyperparamBadFile(t *testing.T) {
// 	ctxt, cleanup := resetApp(t, context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
// 		HyperparamsDir:        "./testData/dataBadFile",
// 		HyperparamsCreateType: types.Create,
// 	})
// 	sbtest.ContainsError(t, types.DataDirErr, err)
// 	sbtest.ContainsError(t, types.UnknownFileErr, err)
// }
//
// func TestMiscTestUploadCSVDataDirHyperparamBadHyperparamType(t *testing.T) {
// 	ctxt, cleanup := resetApp(t, context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
// 		HyperparamsDir:        "./testData/dataBadHyperparamType",
// 		HyperparamsCreateType: types.Create,
// 	})
// 	sbtest.ContainsError(t, types.DataDirErr, err)
// 	sbtest.ContainsError(t, types.UnknownFileErr, err)
// }
//
// func TestMiscTestUploadCSVDataDir(t *testing.T) {
// 	ctxt, cleanup := resetApp(t, context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
// 		ClientDir:             "./testData/clientData",
// 		ClientCreateType:      types.Create,
// 		ExerciseDir:           "./testData/exerciseData",
// 		ExerciseCreateType:    types.Create,
// 		HyperparamsDir:        "./testData/hyperparamData",
// 		HyperparamsCreateType: types.Create,
// 		WorkoutDir:            "./testData/workoutData",
// 		WorkoutCreateType:     types.Create,
// 		BarPathCalcHyperparams: &types.BarPathCalcHyperparams{
// 			MinNumSamples: 5,
// 			ApproxErr:     types.SecondOrder,
// 		},
// 		BarPathTrackerHyperparams: &types.BarPathTrackerHyperparams{},
// 		Opts:                      sbcsv.Opts{TimeFormat: "1/2/2006"},
// 	})
// 	sbtest.Nil(t, err)
//
// 	numClients, err := ReadNumClients(ctxt)
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numClients, 2)
//
// 	numExercises, err := ReadNumExercises(ctxt)
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numExercises, int64(len(migrations.ExerciseSetupData)+3))
//
// 	numHyperparams, err := ReadNumHyperparams(ctxt)
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numHyperparams, int64(len(migrations.HyperparamsSetupData)+4))
//
// 	numHyperparams, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numHyperparams, 3)
// 	numHyperparams, err = ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numHyperparams, 3)
//
// 	numWorkouts, err := ReadClientNumWorkouts(ctxt, "two@gmail.com")
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numWorkouts, 3)
//
// 	numPhysData, err := ReadClientTotalNumPhysEntries(ctxt, "two@gmail.com")
// 	sbtest.Nil(t, err)
// 	sbtest.Eq(t, numPhysData, 1)
// }
