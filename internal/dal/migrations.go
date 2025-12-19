package dal

import (
	"context"
	"embed"

	migrations "code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
	"github.com/jackc/pgx/v5"
)

//go:generate sqlc -f sqlc.yaml generate

//go:embed migrations/*.sql
var SqlMigrations embed.FS
var PostOps = map[sbsqlm.Migration]sbsqlm.PostMigrationOp{
	0: migrations.PostOp0,
}

func RunMigrations(
	ctxt context.Context,
	state *types.State,
	_ pgx.Tx,
	_ struct{},
) (opErr error) {
	if opErr = sbsqlm.Load(SqlMigrations, "migrations", PostOps); opErr != nil {
		return
	}
	if opErr = sbsqlm.Run(ctxt, state.DB); opErr != nil {
		return
	}

	return
}
