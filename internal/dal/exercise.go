package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

func CreateExercisesWithID(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	exercises []types.IdWrapper[types.Exercise],
) error {
	cpy := CpyFromSlice[types.IdWrapper[types.Exercise]]{
		Data: exercises,
		ValueGetter: func(v *types.IdWrapper[types.Exercise], res *[]any) error {
			if len(*res) < 4 {
				*res = make([]any, 4)
			}
			(*res)[0] = v.Id
			(*res)[1] = v.Val.Name
			(*res)[2] = v.Val.KindId
			(*res)[3] = v.Val.FocusId
			return nil
		},
	}
	if n, err := tx.CopyFrom(
		ctxt, pgx.Identifier{"providentia", "exercise"},
		[]string{"id", "name", "kind_id", "focus_id"},
		&cpy,
	); err != nil {
		return sberr.AppendError(types.CouldNotCreateAllExercisesErr, err)
	} else if n != int64(len(exercises)) {
		return sberr.Wrap(
			types.CouldNotCreateAllExercisesErr,
			"Expected to create %d exercises but only created %d, rolling back",
			len(exercises), n,
		)
	}

	if _, err := tx.Exec(ctxt, updateExerciseKindSerialCountSql); err != nil {
		return sberr.Wrap(
			types.CouldNotCreateAllExercisesErr,
			"Failed to update serial index",
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Created new exercises with IDs",
		"NumRows", len(exercises),
	)
	return nil
}

func CreateExercises(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clients []types.Client,
) error {
	// TODO - impl
	// cpy := CpyFromSlice[types.Client]{
	// 	Data: clients,
	// 	ValueGetter: func(v *types.Client, res *[]any) error {
	// 		if len(*res) < 3 {
	// 			*res = make([]any, 3)
	// 		}
	// 		(*res)[0] = v.FirstName
	// 		(*res)[1] = v.LastName
	// 		(*res)[2] = v.Email
	// 		return nil
	// 	},
	// }
	// if n, err := tx.CopyFrom(
	// 	ctxt, pgx.Identifier{"providentia", "client"},
	// 	[]string{"first_name", "last_name", "email"},
	// 	&cpy,
	// ); err != nil {
	// 	return sberr.AppendError(types.CouldNotCreateAllClientsErr, err)
	// } else if n != int64(len(clients)) {
	// 	return sberr.Wrap(
	// 		types.CouldNotCreateAllClientsErr,
	// 		"Expected to create %d clients but only created %d, rolling back",
	// 		len(clients), n,
	// 	)
	// }

	// state.Log.Log(
	// 	ctxt, sblog.VLevel(3),
	// 	"DAL: Created new clients",
	// 	"NumRows", len(clients),
	// )
	// return nil
	return nil
}
