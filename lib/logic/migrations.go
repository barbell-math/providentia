package logic

import (
	"context"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/barbell-math/providentia/internal/ops"
	"github.com/barbell-math/providentia/lib/types"
)

func RunMigrations(ctxt context.Context) (opErr error) {
	return runOp(ctxt, func(state *types.State, queries *dal.Queries) (err error) {
		return ops.RunMigrations(ctxt, state)
	})
}
