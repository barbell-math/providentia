package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

func runOp[T any](
	ctxt context.Context,
	op func(ctxt context.Context, state *types.State, tx pgx.Tx, opts T) error,
	opts T,
) (opErr error) {
	state, ok := StateFromContext(ctxt)
	if !ok {
		opErr = sberr.Wrap(types.InvalidCtxtErr, "Missing State struct")
		return
	}
	if state.DB == nil {
		opErr = sberr.Wrap(
			types.InvalidCtxtErr, "State db was not setup properly",
		)
		return
	}
	if state.Log == nil {
		opErr = sberr.Wrap(
			types.InvalidCtxtErr, "State logging was not setup properly",
		)
		return
	}

	return pgx.BeginTxFunc(
		ctxt, state.DB,
		pgx.TxOptions{},
		func(tx pgx.Tx) error { return op(ctxt, state, tx, opts) },
	)
}
