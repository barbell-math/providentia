package migrations

import (
	"context"
	"errors"

	dal "github.com/barbell-math/providentia/db/dataAccessLayer"
	sberr "github.com/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

var (
	ExerciseFocusSetupData = []dal.BulkCreateExerciseFocusParams{
		{ID: 0, Focus: "UnknownExerciseFocus"},
		{ID: 1, Focus: "Squat"},
		{ID: 2, Focus: "Bench"},
		{ID: 3, Focus: "Deadlift"},
	}

	ExerciseKindSetupData = []dal.BulkCreateExerciseKindParams{
		{
			ID: 1, Kind: "MainCompound",
			Description: "The squat, bench, and deadlift.",
		},
		{
			ID: 2, Kind: "MainCompoundAccessory",
			Description: "Variations of the squat, bench, and deadlift.",
		},
		{
			ID: 3, Kind: "CompoundAccessory",
			Description: "Multi-joint accessories that are not part of the main compound accessory group.",
		},
		{
			ID: 4, Kind: "Accessory",
			Description: "Single joint lifts and core work.",
		},
	}

	SetupMissingDataErr = errors.New("Did not copy all setup values")
)

func postOp0(ctxt context.Context, tx pgx.Tx) error {
	q := dal.New(tx)

	cnt, err := q.BulkCreateExerciseFocus(ctxt, ExerciseFocusSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseFocusSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.ExerciseFocus table: Expected %d rows to be created, got %d",
			len(ExerciseFocusSetupData), cnt,
		)
	}

	cnt, err = q.BulkCreateExerciseKind(ctxt, ExerciseKindSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseKindSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.ExerciseKind table: Expected %d rows to be created, got %d",
			len(ExerciseFocusSetupData), cnt,
		)
	}

	return nil
}
