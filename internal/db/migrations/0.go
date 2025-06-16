package migrations

import (
	"context"
	"errors"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	sberr "github.com/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

var (
	SetupMissingDataErr = errors.New("Did not copy all setup values")
)

func PostOp0(ctxt context.Context, tx pgx.Tx) error {
	q := dal.New(tx)

	cnt, err := q.BulkCreateExerciseFocusWithID(ctxt, ExerciseFocusSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseFocusSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.ExerciseFocus table: Expected %d rows to be created, got %d",
			len(ExerciseFocusSetupData), cnt,
		)
	}

	cnt, err = q.BulkCreateExerciseKindWithID(ctxt, ExerciseKindSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseKindSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.ExerciseKind table: Expected %d rows to be created, got %d",
			len(ExerciseKindSetupData), cnt,
		)
	}

	cnt, err = q.BulkCreateExerciseWithID(ctxt, ExerciseSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.Exercise table: Expected %d rows to be created, got %d",
			len(ExerciseSetupData), cnt,
		)
	}

	cnt, err = q.BulkCreateVideoDataWithID(ctxt, VideoDataSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(VideoDataSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.VideoData table: Expected %d rows to be created, got %d",
			len(VideoDataSetupData), cnt,
		)
	}

	cnt, err = q.BulkCreateModels(ctxt, ModelSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ModelSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.Model table: Expected %d rows to be created, got %d",
			len(ModelSetupData), cnt,
		)
	}

	return nil
}
