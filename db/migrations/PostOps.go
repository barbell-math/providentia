package migrations

import (
	"embed"

	sbsqlm "github.com/barbell-math/smoothbrain-sqlmigrate"
)

//go:embed *.sql
var SqlMigrations embed.FS
var PostOps = map[sbsqlm.Migration]sbsqlm.PostMigrationOp{}
