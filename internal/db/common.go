package db

import (
	"embed"

	migrations "code.barbellmath.net/barbell-math/providentia/internal/db/migrations"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
)

//go:generate sqlc -f sqlc.yaml generate

//go:embed migrations/*.sql
var SqlMigrations embed.FS
var PostOps = map[sbsqlm.Migration]sbsqlm.PostMigrationOp{
	0: migrations.PostOp0,
}
