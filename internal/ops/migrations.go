package ops

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/db"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
)

func RunMigrations(ctxt context.Context, state *types.State) (opErr error) {
	if opErr = sbsqlm.Load(
		db.SqlMigrations, "migrations", db.PostOps,
	); opErr != nil {
		return
	}
	if opErr = sbsqlm.Run(ctxt, state.DB); opErr != nil {
		return
	}

	return
}
