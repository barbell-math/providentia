package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Runs all necessary migrations. Previously run migrations will not run again
// unless the `smoothbrain_sqlmigrate` table was corrupted.
func RunMigrations(ctxt context.Context) (opErr error) {
	// A transaction is created by the smoothbrain_sqlmigrate lib in an internal
	// call so do not call [runOp] here. It would create an unnecessary
	// transaction.
	var state *types.State
	state, opErr = getState(ctxt)
	return migrations.RunMigrations(ctxt, state)
}
