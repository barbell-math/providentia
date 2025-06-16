package db

import (
	"embed"

	migrations "github.com/barbell-math/providentia/internal/db/migrations"
	sbsqlm "github.com/barbell-math/smoothbrain-sqlmigrate"
)

//go:embed migrations/*.sql
var SqlMigrations embed.FS
var PostOps = map[sbsqlm.Migration]sbsqlm.PostMigrationOp{
	0: migrations.PostOp0,
}
