package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
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

// Bulk uploads the data referenced by the [types.BulkUploadData] struct.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func BulkUploadData(
	ctxt context.Context,
	opts *types.BulkUploadDataOpts,
) (opErr error) {
	return runOp(ctxt, jobs.BulkUploadData, opts)
}
