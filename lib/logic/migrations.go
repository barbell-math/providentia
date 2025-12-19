package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
)

// Runs all necessary migrations. Previously run migrations will not run again
// unless the `smoothbrain_sqlmigrate` table was corrupted.
func RunMigrations(ctxt context.Context) (opErr error) {
	return runOp(ctxt, dal.RunMigrations, struct{}{})
}
