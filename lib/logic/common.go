package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

func getState(ctxt context.Context) (*types.State, error) {
	state, ok := StateFromContext(ctxt)
	if !ok {
		return nil, sberr.Wrap(types.InvalidCtxtErr, "Missing State struct")
	}
	if state.DB == nil {
		return nil, sberr.Wrap(
			types.InvalidCtxtErr, "State db was not setup properly",
		)
	}
	if state.Log == nil {
		return nil, sberr.Wrap(
			types.InvalidCtxtErr, "State logging was not setup properly",
		)
	}
	return state, nil
}

func runOp[T any](
	ctxt context.Context,
	op func(ctxt context.Context, state *types.State, tx pgx.Tx, opts T) error,
	opts T,
) error {
	state, err := getState(ctxt)
	if err != nil {
		return err
	}

	return pgx.BeginTxFunc(
		ctxt, state.DB,
		pgx.TxOptions{},
		func(tx pgx.Tx) error { return op(ctxt, state, tx, opts) },
	)
}
