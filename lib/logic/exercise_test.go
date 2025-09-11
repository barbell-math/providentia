package logic

import (
	"context"
	"fmt"
	"math"
	"testing"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/db/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestExercisesTypeConversions(t *testing.T) {
	structsEquivalent[dal.BulkCreateExercisesParams, types.Exercise](t)
}

func TestExercise(t *testing.T) {
	t.Run("failingNoWrites", exerciseFailingNoWrites)
	t.Run("duplicateName", exerciseDuplicateName)
	t.Run("transactionRollback", exerciseTransactionRollback)
	t.Run("addGet", exerciseAddGet)
	t.Run("addUpdateGet", exerciseAddUpdateGet)
	t.Run("addDeleteGet", exerciseAddDeleteGet)
}

func exerciseFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)
	t.Run("missingName", exerciseMissingName(ctxt))
	t.Run("invalidFocusID", exerciseInvalidFocusID(ctxt))
	t.Run("invalidKindID", exerciseInvalidKindID(ctxt))

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData)), numExercises)
}

func exerciseMissingName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateExercises(ctxt, types.Exercise{
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		})
		sbtest.ContainsError(t, types.InvalidExerciseErr, err)
	}
}

func exerciseInvalidFocusID(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateExercises(ctxt, types.Exercise{
			Name:    "asdf",
			KindID:  types.MainCompound,
			FocusID: types.ExerciseFocus(math.MaxInt32),
		})
		sbtest.ContainsError(t, types.InvalidExerciseErr, err)
	}
}

func exerciseInvalidKindID(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateExercises(ctxt, types.Exercise{
			Name:    "asdf",
			KindID:  types.ExerciseKind(math.MaxInt32),
			FocusID: types.Bench,
		})
		sbtest.ContainsError(t, types.InvalidExerciseErr, err)
	}
}

func exerciseDuplicateName(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	exercises := make([]types.Exercise, 13)
	for i := range len(exercises) {
		exercises[i] = types.Exercise{
			Name:    fmt.Sprintf("testExercise%d", i),
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		}
	}
	exercises[len(exercises)-1].Name = fmt.Sprintf(
		"testExercise%d", len(exercises)-2,
	)

	err := CreateExercises(ctxt, exercises...)
	sbtest.ContainsError(t, types.CouldNotAddExercisesErr, err)

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData)), numExercises)
}

func exerciseTransactionRollback(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	exercises := make([]types.Exercise, 13)
	for i := range len(exercises) {
		exercises[i] = types.Exercise{
			Name:    fmt.Sprintf("testExercise%d", i),
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		}
	}

	err := CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)
	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)

	for i := 0; i < 5; i++ {
		exercises[i].Name = fmt.Sprintf("testExercise%d", i+len(exercises))
	}

	err = CreateExercises(ctxt, exercises...)
	sbtest.ContainsError(t, types.CouldNotAddExercisesErr, err)
	numExercises, err = ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)
}

func exerciseAddGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	exercises := make([]types.Exercise, 13)
	for i := range len(exercises) {
		exercises[i] = types.Exercise{
			Name:    fmt.Sprintf("testExercise%d", i),
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		}
	}

	err := CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)

	for i := range len(exercises) {
		res, err := ReadExercisesByName(ctxt, exercises[i].Name)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, exercises[i].Name, res[0].Name)
		sbtest.Eq(t, exercises[i].KindID, res[0].KindID)
		sbtest.Eq(t, exercises[i].FocusID, res[0].FocusID)
	}

	_, err = ReadExercisesByName(ctxt, "badExercise")
	sbtest.ContainsError(t, types.CouldNotFindRequestedExerciseErr, err)
}

func exerciseAddUpdateGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	exercises := make([]types.Exercise, 13)
	for i := range len(exercises) {
		exercises[i] = types.Exercise{
			Name:    fmt.Sprintf("testExercise%d", i),
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		}
	}

	err := CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)

	for i := range len(exercises) {
		exercises[i].KindID = types.Accessory
		exercises[i].FocusID = types.Squat
	}
	err = UpdateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	err = UpdateExercises(ctxt, types.Exercise{Name: "badExercise"})
	sbtest.Nil(t, err)
	numExercises, err = ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)

	for i := range len(exercises) {
		res, err := ReadExercisesByName(ctxt, exercises[i].Name)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, exercises[i].Name, res[0].Name)
		sbtest.Eq(t, exercises[i].KindID, res[0].KindID)
		sbtest.Eq(t, exercises[i].FocusID, res[0].FocusID)
	}

	_, err = ReadExercisesByName(ctxt, "badExercise")
	sbtest.ContainsError(t, types.CouldNotFindRequestedExerciseErr, err)
}

func exerciseAddDeleteGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	exercises := make([]types.Exercise, 13)
	for i := range len(exercises) {
		exercises[i] = types.Exercise{
			Name:    fmt.Sprintf("testExercise%d", i),
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		}
	}

	err := CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	numExercises, err := ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+13, numExercises)

	names := [5]string{}
	for i := range len(names) {
		names[i] = exercises[i].Name
	}
	err = DeleteExercises(ctxt, names[:]...)
	sbtest.Nil(t, err)
	numExercises, err = ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(
		len(migrations.ExerciseSetupData)+len(exercises)-len(names),
	), numExercises)

	err = DeleteExercises(ctxt, "badExercise")
	sbtest.ContainsError(t, types.CouldNotDeleteRequestedExerciseErr, err)
	sbtest.ContainsError(t, types.CouldNotFindRequestedExerciseErr, err)
	numExercises, err = ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(
		len(migrations.ExerciseSetupData)+len(exercises)-len(names),
	), numExercises)

	for i := range numExercises - int64(len(migrations.ExerciseSetupData)) {
		offset := int(i) + len(names)
		res, err := ReadExercisesByName(ctxt, exercises[offset].Name)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, exercises[offset].Name, res[0].Name)
		sbtest.Eq(t, exercises[offset].KindID, res[0].KindID)
		sbtest.Eq(t, exercises[offset].FocusID, res[0].FocusID)
	}

	_, err = ReadExercisesByName(ctxt, "badExercise")
	sbtest.ContainsError(t, types.CouldNotFindRequestedExerciseErr, err)
}
