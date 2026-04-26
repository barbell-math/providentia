package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbjobqueue"
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

// Initializes the providentia library by:
//
//   - Createing a cancelable context with the state variable
//   - Creating a job poller and running it in a separate go routine
//   - Checking if database migrations need to be run and running them
//
// This is a generic init function. If you wish to customize the job poller to
// include queues from your application or any other custom initialization logic
// you must perform the steps listed above somewhere in your code before calling
// any providentia library functions.
func Init(
	ctxt context.Context,
	state *types.State,
) (provLifetime context.Context, cleanup func(), opErr error) {
	appLifetime, appCancel := context.WithCancel(ctxt)
	provLifetime = WithStateValue(appLifetime, state)

	go sbjobqueue.Poll(
		appLifetime, state.PhysicsJobQueue, state.CSVLoaderJobQueue,
	)
	cleanup = func() {
		appCancel()
		CleanupState(state)
	}

	opErr = RunMigrations(provLifetime)
	return
}
