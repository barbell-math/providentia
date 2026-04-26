package types

import (
	"log/slog"

	"code.barbellmath.net/carmichaeljr/smoothbrain/sbjobqueue"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PhysicsJob        struct{} // Used to identify a physics job
	CSVLoaderJob      struct{} // Used to identify a csv loader job
	GeneralPurposeJob struct{} // Used to identify a general purpose job

	// The state the rest of providentia will use. Almost all functions
	// available for external use from this library will require this state to
	// be available in the passed in context.
	State struct {
		Log *slog.Logger
		DB  *pgxpool.Pool

		PhysicsJobQueue   *sbjobqueue.JobQueue[PhysicsJob]
		CSVLoaderJobQueue *sbjobqueue.JobQueue[CSVLoaderJob]

		Global GlobalConf
	}
)
