package tests

import (
	"context"
	"fmt"
	"math"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestExercise(t *testing.T) {
	t.Run("failingNoWrites", exerciseFailingNoWrites)
	t.Run("duplicateName", exerciseDuplicateName)
	t.Run("createRead", exerciseCreateRead)
	t.Run("ensureRead", exerciseEnsureRead)
	t.Run("createFind", exerciseCreateFind)
	t.Run("createDeleteRead", exerciseCreateDeleteRead)
	t.Run("createCSVRead", exerciseCreateCSVRead)
	t.Run("ensureCSVRead", exerciseEnsureCSVRead)
}

func exerciseFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("missingName", exerciseMissingName(ctxt))
	t.Run("invalidKindId", exerciseInvalidKindId(ctxt))
	t.Run("invalidFocusId", exerciseInvalidFocusId(ctxt))

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData)), n)
}

func exerciseMissingName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateExercises(ctxt, types.Exercise{
			KindId:  types.MainCompound,
			FocusId: types.Bench,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllExercisesErr, err,
			`ERROR: new row for relation "exercise" violates check constraint "name_not_empty" \(SQLSTATE 23514\)`,
		)
	}
}

func exerciseInvalidFocusId(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateExercises(ctxt, types.Exercise{
			Name:    "asdf",
			KindId:  types.MainCompound,
			FocusId: types.ExerciseFocus(math.MaxInt32),
		})
		fmt.Println(err)
		sbtest.ContainsError(
			t, types.CouldNotCreateAllExercisesErr, err,
			`insert or update on table "exercise" violates foreign key constraint "exercise_focus_id_fkey" \(SQLSTATE 23503\)`,
		)
	}
}

func exerciseInvalidKindId(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateExercises(ctxt, types.Exercise{
			Name:    "asdf",
			KindId:  types.ExerciseKind(math.MaxInt32),
			FocusId: types.Bench,
		})
		fmt.Println(err)
		sbtest.ContainsError(
			t, types.CouldNotCreateAllExercisesErr, err,
			`insert or update on table "exercise" violates foreign key constraint "exercise_kind_id_fkey" \(SQLSTATE 23503\)`,
		)
	}
}

func exerciseDuplicateName(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateExercises(ctxt, types.Exercise{
		Name:    "testExercise",
		KindId:  types.MainCompound,
		FocusId: types.Bench,
	}, types.Exercise{
		Name:    "testExercise",
		KindId:  types.MainCompound,
		FocusId: types.Bench,
	})
	fmt.Println(err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllExercisesErr, err,
		`duplicate key value violates unique constraint "exercise_name_key" \(SQLSTATE 23505\)`,
	)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData)), n)
}

func exerciseCreateRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	exercises := []types.Exercise{
		{
			Name:    "testExercise",
			KindId:  types.MainCompound,
			FocusId: types.Squat,
		}, {
			Name:    "testExercise1",
			KindId:  types.MainCompound,
			FocusId: types.Bench,
		},
	}
	err := logic.CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	readExercises, err := logic.ReadExercisesByName(
		ctxt, exercises[0].Name, exercises[1].Name,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, exercises, readExercises)

	readExercises, err = logic.ReadExercisesByName(ctxt, "asdfasdf")
	sbtest.ContainsError(
		t, types.CouldNotReadAllExercisesErr, err,
		"Only read 0 entries out of batch of 1 requests",
	)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)
}

func exerciseEnsureRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	exercises := []types.Exercise{
		{
			Name:    "testExercise",
			KindId:  types.MainCompound,
			FocusId: types.Squat,
		}, {
			Name:    "testExercise1",
			KindId:  types.MainCompound,
			FocusId: types.Bench,
		},
	}
	err := logic.EnsureExercisesExist(ctxt, exercises...)
	sbtest.Nil(t, err)

	readExercises, err := logic.ReadExercisesByName(
		ctxt, exercises[0].Name, exercises[1].Name,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, exercises, readExercises)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)

	err = logic.EnsureExercisesExist(ctxt, exercises...)
	sbtest.Nil(t, err)

	readExercises, err = logic.ReadExercisesByName(
		ctxt, exercises[0].Name, exercises[1].Name,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, exercises, readExercises)

	n, err = logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)
}

func exerciseCreateFind(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	exercises := []types.Exercise{
		{
			Name:    "testExercise",
			KindId:  types.MainCompound,
			FocusId: types.Squat,
		}, {
			Name:    "testExercise1",
			KindId:  types.MainCompound,
			FocusId: types.Bench,
		},
	}
	err := logic.CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	foundExercises, err := logic.FindExercisesByName(
		ctxt, exercises[0].Name, exercises[1].Name, "asdfasdf",
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, []types.Found[types.Exercise]{{
		Found: true,
		Value: exercises[0],
	}, {
		Found: true,
		Value: exercises[1],
	}, {
		Found: false,
	}}, foundExercises)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)
}

func exerciseCreateDeleteRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	exercises := []types.Exercise{
		{
			Name:    "testExercise",
			KindId:  types.MainCompound,
			FocusId: types.Squat,
		}, {
			Name:    "testExercise1",
			KindId:  types.MainCompound,
			FocusId: types.Bench,
		}, {
			Name:    "testExercise2",
			KindId:  types.MainCompound,
			FocusId: types.Deadlift,
		},
	}
	err := logic.CreateExercises(ctxt, exercises...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+3, n)

	err = logic.DeleteExercises(ctxt, exercises[0].Name)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)

	readExercises, err := logic.ReadExercisesByName(
		ctxt, exercises[1].Name, exercises[2].Name,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, exercises[1:], readExercises)

	err = logic.DeleteExercises(ctxt, exercises[0].Name)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllExercisesErr, err,
		`Could not delete entry with id 'testExercise' \(Does id exist\?\)`,
	)

	n, err = logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+2, n)
}

func exerciseCreateCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateExercisesFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/exerciseData/exercises.csv",
	)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+3, n)

	client, err := logic.ReadExercisesByName(ctxt, "one", "two", "three")
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, client, []types.Exercise{{
		Name:    "one",
		KindId:  types.MainCompoundAccessory,
		FocusId: types.UnknownExerciseFocus,
	}, {
		Name:    "two",
		KindId:  types.MainCompound,
		FocusId: types.Squat,
	}, {
		Name:    "three",
		KindId:  types.Accessory,
		FocusId: types.Deadlift,
	}})

	_, err = logic.ReadExercisesByName(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotReadAllExercisesErr, err)

	err = logic.CreateExercisesFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/exerciseData/exercises.csv",
	)
	sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllExercisesErr, err,
		`duplicate key value violates unique constraint "exercise_name_key" \(SQLSTATE 23505\)`,
	)

	n, err = logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+3, n)
}

func exerciseEnsureCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.EnsureExercisesExistFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/exerciseData/exercises.csv",
	)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+3, n)

	client, err := logic.ReadExercisesByName(ctxt, "one", "two", "three")
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, client, []types.Exercise{{
		Name:    "one",
		KindId:  types.MainCompoundAccessory,
		FocusId: types.UnknownExerciseFocus,
	}, {
		Name:    "two",
		KindId:  types.MainCompound,
		FocusId: types.Squat,
	}, {
		Name:    "three",
		KindId:  types.Accessory,
		FocusId: types.Deadlift,
	}})

	_, err = logic.ReadExercisesByName(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotReadAllExercisesErr, err)

	err = logic.EnsureExercisesExistFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/exerciseData/exercises.csv",
	)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumExercises(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.ExerciseSetupData))+3, n)
}
