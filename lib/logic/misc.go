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
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.RunMigrations(ctxt, state)
		},
	})
}

func UploadCSVDataDir(
	ctxt context.Context,
	dir string,
	opts *types.CSVDataDirOptions,
) (opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			err = ops.UploadCSVDataDir(ctxt, state, queries, dir, opts)
			return
		},
	})
	return
}
