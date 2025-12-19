// All logic that works with the library's exposed types is exposed through this
// package.
package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
)

type (
	opCalls struct {
		preOp  func(state *types.State, queries *dal.SyncQueries) (err error)
		op     func(state *types.State, queries *dal.SyncQueries) (err error)
		postOp func(state *types.State, queries *dal.SyncQueries, opSucceeded bool)
	}
)

func runOp(ctxt context.Context, calls opCalls) (opErr error) {
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
	syncQueries := dal.NewSyncQueries(queries)

	defer func() {
		if opErr != nil {
			tx.Rollback(ctxt)
		} else {
			tx.Commit(ctxt)
		}
		if calls.postOp != nil {
			calls.postOp(state, syncQueries, opErr == nil)
		}
	}()

	errRes := make(chan error)
	if calls.preOp != nil {
		go func() { errRes <- calls.preOp(state, syncQueries) }()

		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
		case err := <-errRes:
			opErr = err
			return
		}
	}

	if calls.op != nil {
		go func() { errRes <- calls.op(state, syncQueries) }()

		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
		case err := <-errRes:
			opErr = err
			return
		}
	}

	return
}
