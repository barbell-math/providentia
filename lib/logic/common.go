package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

func runOp(
	ctxt context.Context,
	op func(state *types.State, queries *dal.Queries) (err error),
) (opErr error) {
	state, ok := StateFromContext(ctxt)
	if !ok {
		opErr = sberr.Wrap(types.InvalidCtxtErr, "missing State struct")
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

	var tx pgx.Tx
	tx, opErr = state.DB.Begin(ctxt)
	if opErr != nil {
		return
	}

	q := dal.New(state.DB)
	queries := q.WithTx(tx)
	defer func() {
		if opErr != nil {
			tx.Rollback(ctxt)
		} else {
			tx.Commit(ctxt)
		}
	}()

	errRes := make(chan error)
	go func() { errRes <- op(state, queries) }()

	select {
	case <-ctxt.Done():
		opErr = ctxt.Err()
	case err := <-errRes:
		opErr = err
	}
	return
}
