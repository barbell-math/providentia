package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

func RunMigrations(ctxt context.Context) (opErr error) {
	return runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		return ops.RunMigrations(ctxt, state)
	})
}
