package provlib

import (
	"context"
	"errors"
	"iter"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/jackc/pgx/v5"
)

var (
	InvalidCtxtErr = errors.New("Invalid context: missing state struct")
)

func dbOpInfo(
	ctxt context.Context,
) (state *State, queries *dal.Queries, _defer func(err *error), err error) {
	var ok bool
	state, ok = FromContext(ctxt)
	if !ok {
		err = InvalidCtxtErr
		return
	}

	var tx pgx.Tx
	tx, err = state.DB.Begin(ctxt)
	if err != nil {
		return
	}

	q := dal.New(state.DB)
	queries = q.WithTx(tx)
	_defer = func(err *error) {
		if *err != nil {
			tx.Rollback(ctxt)
		} else {
			tx.Commit(ctxt)
		}
	}

	return
}

func batchIndexes[S ~[]E, E any](s S, step int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := 0; i < len(s); i += step {
			end := i + step
			if !yield(i, min(len(s), end)) {
				break
			}
		}
	}
}
