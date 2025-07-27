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
	SerialUpdateCmdErr  = errors.New("Could not update the serial value")
)

func PostOp0(ctxt context.Context, tx pgx.Tx) error {
	q := dal.New(tx)

	cnt, err := q.BulkCreateExerciseFocusWithID(ctxt, ExerciseFocusSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseFocusSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.exercise_focus table: Expected %d rows to be created, got %d",
			len(ExerciseFocusSetupData), cnt,
		)
	}
	if err := q.UpdateExerciseFocusSerialCount(ctxt); err != nil {
		return sberr.Wrap(
			SerialUpdateCmdErr,
			"Setting up the providentia.exercise_focus table",
		)
	}

	cnt, err = q.BulkCreateExerciseKindWithID(ctxt, ExerciseKindSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseKindSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.exercise_kind table: Expected %d rows to be created, got %d",
			len(ExerciseKindSetupData), cnt,
		)
	}
	if err := q.UpdateExerciseKindSerialCount(ctxt); err != nil {
		return sberr.Wrap(
			SerialUpdateCmdErr,
			"Setting up the providentia.exercise_kind table",
		)
	}

	cnt, err = q.BulkCreateExerciseWithID(ctxt, ExerciseSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ExerciseSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.exercise table: Expected %d rows to be created, got %d",
			len(ExerciseSetupData), cnt,
		)
	}
	if err := q.UpdateExerciseSerialCount(ctxt); err != nil {
		return sberr.Wrap(
			SerialUpdateCmdErr,
			"Setting up the providentia.exercise table",
		)
	}

	cnt, err = q.BulkCreateVideoDataWithID(ctxt, VideoDataSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(VideoDataSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.video_data table: Expected %d rows to be created, got %d",
			len(VideoDataSetupData), cnt,
		)
	}
	if err := q.UpdateVideoDataSerialCount(ctxt); err != nil {
		return sberr.Wrap(
			SerialUpdateCmdErr,
			"Setting up the providentia.video_data table",
		)
	}

	cnt, err = q.BulkCreateModelsWithID(ctxt, ModelSetupData)
	if err != nil {
		return err
	} else if cnt != int64(len(ModelSetupData)) {
		return sberr.Wrap(
			SetupMissingDataErr,
			"Setting up the providentia.model table: Expected %d rows to be created, got %d",
			len(ModelSetupData), cnt,
		)
	}
	if err := q.UpdateModelSerialCount(ctxt); err != nil {
		return sberr.Wrap(
			SerialUpdateCmdErr,
			"Setting up the providentia.model table",
		)
	}

	return nil
}
