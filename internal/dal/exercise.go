package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

type (
	ReadExerciseByNameOpts struct {
		Names     []string
		Exercises *[]types.Exercise
	}

	FindExerciseByNameOpts struct {
		Names     []string
		Exercises *[]types.Found[types.Exercise]
	}
)

const (
	exerciseTableName = "exercise"
)

func CreateExercisesWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	exercises []types.IdWrapper[types.Exercise],
) error {
	return genericCreateWithId(
		ctxt, state, tx, &genericCreateOpts[types.IdWrapper[types.Exercise]]{
			TableName: exerciseTableName,
			Columns:   []string{"id", "name", "kind_id", "focus_id"},
			Data:      exercises,
			ValueGetter: func(v *types.IdWrapper[types.Exercise], res *[]any) error {
				*res = util.SliceClamp(*res, 4)
				(*res)[0] = v.Id
				(*res)[1] = v.Val.Name
				(*res)[2] = v.Val.KindId
				(*res)[3] = v.Val.FocusId
				return nil
			},
			Err: types.CouldNotCreateAllExercisesErr,
		},
	)
}

func CreateExercises(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	exercises []types.Exercise,
) error {
	return genericCreate(
		ctxt, state, tx, &genericCreateOpts[types.Exercise]{
			TableName: exerciseTableName,
			Columns:   []string{"name", "kind_id", "focus_id"},
			Data:      exercises,
			ValueGetter: func(v *types.Exercise, res *[]any) error {
				*res = util.SliceClamp(*res, 3)
				(*res)[0] = v.Name
				(*res)[1] = v.KindId
				(*res)[2] = v.FocusId
				return nil
			},
			Err: types.CouldNotCreateAllExercisesErr,
		},
	)
}

func EnsureExercisesExist(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	exercises []types.Exercise,
) error {
	return genericEnsureExists(
		ctxt, state, tx, &genericCreateOpts[types.Exercise]{
			TableName: exerciseTableName,
			Columns:   []string{"name", "kind_id", "focus_id"},
			Data:      exercises,
			ValueGetter: func(v *types.Exercise, res *[]any) error {
				*res = make([]any, 3)
				(*res)[0] = v.Name
				(*res)[1] = v.KindId
				(*res)[2] = v.FocusId
				return nil
			},
			Err: types.CouldNotCreateAllExercisesErr,
		},
	)
}

func ReadNumExercises(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	num *int64,
) error {
	return genericReadTotalNum(
		ctxt, state, tx, &genericReadTotalNumOpts{
			TableName: exerciseTableName,
			Res:       num,
		},
	)
}

func ReadExercisesByName(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadExerciseByNameOpts,
) error {
	return genericReadByUniqueId(
		ctxt, state, tx, &genericReadByUniqueIdOpts[string, types.Exercise]{
			TableName:  exerciseTableName,
			Columns:    []string{"name", "kind_id", "focus_id"},
			UniqueCol:  "name",
			IdsSqlType: "TEXT",
			Ids:        opts.Names,
			Res:        opts.Exercises,
			Err:        types.CouldNotReadAllExercisesErr,
		},
	)
}

func FindExercisesByName(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindExerciseByNameOpts,
) error {
	return genericFindByUniqueId(
		ctxt, state, tx, &genericFindByUniqueIdOpts[
			string, types.Found[types.Exercise], types.Exercise,
		]{
			TableName:  exerciseTableName,
			Columns:    []string{"name", "kind_id", "focus_id"},
			UniqueCol:  "name",
			IdsSqlType: "TEXT",
			Ids:        opts.Names,
			Res:        opts.Exercises,
			SetScanValues: func(v *types.Exercise, res []any) {
				res[0] = &v.Name
				res[1] = &v.KindId
				res[2] = &v.FocusId
			},
			Err: types.CouldNotReadAllExercisesErr,
		},
	)
}

func DeleteExercises(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	emails []string,
) error {
	return genericDeleteByUniqueId(
		ctxt, state, tx, &genericDeleteByUniqueIdOpts[string]{
			Ids:       emails,
			TableName: exerciseTableName,
			UniqueCol: "name",
			Err:       types.CouldNotDeleteAllExercisesErr,
		},
	)
}
