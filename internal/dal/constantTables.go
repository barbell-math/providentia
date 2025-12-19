package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	CreateExerciseFocusWithIDOpts struct {
		types.ExerciseFocus `db:"id"`
		Desc                string `db:"focus"`
	}

	CreateExerciseKindWithIDOpts struct {
		types.ExerciseKind `db:"id"`
		Name               string `db:"kind"`
		Desc               string `db:"description"`
	}
)

const (
	updateExerciseFocusSerialCountSql = `
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise_focus', 'id'),
	(SELECT MAX(id) FROM providentia.exercise_focus) + 1
);
`

	updateExerciseKindSerialCountSql = `
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise_kind', 'id'),
	(SELECT MAX(id) FROM providentia.exercise_kind) + 1
);
`
)

func CreateExerciseFocusWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []CreateExerciseFocusWithIDOpts,
) error {
	cpy := CpyFromSlice[CreateExerciseFocusWithIDOpts]{
		Data: data,
		ValueGetter: func(v *CreateExerciseFocusWithIDOpts, res *[]any) error {
			if len(*res) < 2 {
				*res = make([]any, 2)
			}
			(*res)[0] = v.ExerciseFocus
			(*res)[1] = v.Desc
			return nil
		},
	}
	if n, err := tx.CopyFrom(
		ctxt, pgx.Identifier{"providentia", "exercise_focus"},
		[]string{"id", "focus"},
		&cpy,
	); err != nil {
		return sberr.AppendError(
			types.CouldNotCreateAllExerciseFocusEntriesErr, err,
		)
	} else if n != int64(len(data)) {
		return sberr.Wrap(
			types.CouldNotCreateAllExerciseFocusEntriesErr,
			"Expected to create %d entries but only created %d, rolling back",
			len(data), n,
		)
	}

	if _, err := tx.Exec(ctxt, updateExerciseFocusSerialCountSql); err != nil {
		return sberr.Wrap(
			types.CouldNotCreateAllExerciseFocusEntriesErr,
			"Failed to update serial index",
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Created new exercise focus entries",
		"NumRows", len(data),
	)
	return nil
}

func CreateExerciseKindWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []CreateExerciseKindWithIDOpts,
) error {
	cpy := CpyFromSlice[CreateExerciseKindWithIDOpts]{
		Data: data,
		ValueGetter: func(v *CreateExerciseKindWithIDOpts, res *[]any) error {
			if len(*res) < 3 {
				*res = make([]any, 3)
			}
			(*res)[0] = v.ExerciseKind
			(*res)[1] = v.Name
			(*res)[2] = v.Desc
			return nil
		},
	}
	if n, err := tx.CopyFrom(
		ctxt, pgx.Identifier{"providentia", "exercise_kind"},
		[]string{"id", "kind", "description"},
		&cpy,
	); err != nil {
		return sberr.AppendError(
			types.CouldNotCreateAllExerciseKindEntriesErr, err,
		)
	} else if n != int64(len(data)) {
		return sberr.Wrap(
			types.CouldNotCreateAllExerciseKindEntriesErr,
			"Expected to create %d entries but only created %d, rolling back",
			len(data), n,
		)
	}

	if _, err := tx.Exec(ctxt, updateExerciseKindSerialCountSql); err != nil {
		return sberr.Wrap(
			types.CouldNotCreateAllExerciseKindEntriesErr,
			"Failed to update serial index",
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Created new exercise kind entries",
		"NumRows", len(data),
	)
	return nil
}
