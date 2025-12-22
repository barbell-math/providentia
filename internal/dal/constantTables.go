package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
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

	CreateModelsWithIDOpts struct {
		types.ModelID `db:"id"`
		Name          string `db:"kind"`
		Desc          string `db:"description"`
	}
)

func CreateExerciseFocusWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []CreateExerciseFocusWithIDOpts,
) error {
	return genericCreateWithId(
		ctxt, state, tx, &genericCreateOpts[CreateExerciseFocusWithIDOpts]{
			TableName: "exercise_focus",
			Columns:   []string{"id", "focus"},
			Data:      data,
			ValueGetter: func(v *CreateExerciseFocusWithIDOpts, res *[]any) error {
				*res = util.SliceClamp(*res, 2)
				(*res)[0] = v.ExerciseFocus
				(*res)[1] = v.Desc
				return nil
			},
			Err: types.CouldNotCreateAllExerciseFocusEntriesErr,
		},
	)
}

func CreateExerciseKindWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []CreateExerciseKindWithIDOpts,
) error {
	return genericCreateWithId(
		ctxt, state, tx, &genericCreateOpts[CreateExerciseKindWithIDOpts]{
			TableName: "exercise_kind",
			Columns:   []string{"id", "kind", "description"},
			Data:      data,
			ValueGetter: func(v *CreateExerciseKindWithIDOpts, res *[]any) error {
				*res = util.SliceClamp(*res, 3)
				(*res)[0] = v.ExerciseKind
				(*res)[1] = v.Name
				(*res)[2] = v.Desc
				return nil
			},
			Err: types.CouldNotCreateAllExerciseKindEntriesErr,
		},
	)
}

func CreateModelsWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []CreateModelsWithIDOpts,
) error {
	return genericCreateWithId(
		ctxt, state, tx, &genericCreateOpts[CreateModelsWithIDOpts]{
			TableName: "model",
			Columns:   []string{"id", "name", "description"},
			Data:      data,
			ValueGetter: func(v *CreateModelsWithIDOpts, res *[]any) error {
				*res = util.SliceClamp(*res, 3)
				(*res)[0] = v.ModelID
				(*res)[1] = v.Name
				(*res)[2] = v.Desc
				return nil
			},
			Err: types.CouldNotCreateAllModelsErr,
		},
	)
}
