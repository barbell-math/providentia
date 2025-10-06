package logic

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/db/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestMiscTestUploadCSVDataDirClientBadDir(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		ClientDir:        "./testData/dataBadDir",
		ClientCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.DirInDataDirErr, err)
}

func TestMiscTestUploadCSVDataDirClientBadFile(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		ClientDir:        "./testData/dataBadFile",
		ClientCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.UnknownFileErr, err)
}

func TestMiscTestUploadCSVDataDirExerciseBadDir(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		ExerciseDir:        "./testData/dataBadDir",
		ExerciseCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.DirInDataDirErr, err)
}

func TestMiscTestUploadCSVDataDirExerciseBadFile(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		ExerciseDir:        "./testData/dataBadFile",
		ExerciseCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.UnknownFileErr, err)
}

func TestMiscTestUploadCSVDataDirHyperparamBadDir(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		HyperparamsDir:        "./testData/dataBadDir",
		HyperparamsCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.DirInDataDirErr, err)
}

func TestMiscTestUploadCSVDataDirHyperparamBadFile(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		HyperparamsDir:        "./testData/dataBadFile",
		HyperparamsCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.UnknownFileErr, err)
}

func TestMiscTestUploadCSVDataDirHyperparamBadHyperparamType(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		HyperparamsDir:        "./testData/dataBadHyperparamType",
		HyperparamsCreateType: types.Create,
	})
	sbtest.ContainsError(t, types.DataDirErr, err)
	sbtest.ContainsError(t, types.UnknownFileErr, err)
}

func TestMiscTestUploadCSVDataDir(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := UploadCSVDataDir(ctxt, &types.CSVDataDirOptions{
		ClientDir:                 "./testData/clientData",
		ClientCreateType:          types.Create,
		ExerciseDir:               "./testData/exerciseData",
		ExerciseCreateType:        types.Create,
		HyperparamsDir:            "./testData/hyperparamData",
		HyperparamsCreateType:     types.Create,
		WorkoutDir:                "./testData/workoutData",
		WorkoutCreateType:         types.Create,
		BarPathCalcHyperparams:    &types.BarPathCalcHyperparams{},
		BarPathTrackerHyperparams: &types.BarPathTrackerHyperparams{},
		Opts:                      sbcsv.Opts{TimeFormat: "1/2/2006"},
	})
	sbtest.Nil(t, err)

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numClients, 2)

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numExercises, int64(len(migrations.ExerciseSetupData)+3))

	numHyperparams, err := ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, int64(len(migrations.HyperparamsSetupData)+4))

	numHyperparams, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, 3)
	numHyperparams, err = ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numHyperparams, 3)

	numWorkouts, err := ReadClientNumWorkouts(ctxt, "two@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, numWorkouts, 3)
}
