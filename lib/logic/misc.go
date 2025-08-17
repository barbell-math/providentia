package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Runs all database migrations if they have not been run already. This should
// be called as part of a setup or initialization routine.
func RunMigrations(ctxt context.Context) (opErr error) {
	return runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		return ops.RunMigrations(ctxt, state)
	})
}
