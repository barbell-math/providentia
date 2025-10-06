package logic

import (
	"context"
	"testing"
	"time"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestCancelation(t *testing.T) {
	timeoutCtxt, cancel := context.WithTimeout(
		context.Background(), 500*time.Millisecond,
	)
	defer cancel()

	ctxt, cleanup := resetApp(t,timeoutCtxt)
	t.Cleanup(cleanup)

	err := runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			time.Sleep(5 * time.Second)
			return nil
		},
	})
	sbtest.ContainsError(t, context.DeadlineExceeded, err)
}
