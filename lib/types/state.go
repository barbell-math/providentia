package types

import (
	"log/slog"

	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	PhysicsJob        struct{} // Used to identify a physics job
	VideoJob          struct{} // Used to identify a video job
	CSVLoaderJob      struct{} // Used to identify a csv loader job
	GeneralPurposeJob struct{} // Used to identify a general purpose job

	// The state the rest of providentia will use. Almost all functions
	// available for external use from this library will require this state to
	// be available in the passed in context.
	State struct {
		Log *slog.Logger
		DB  *pgxpool.Pool

		PhysicsJobQueue   *sbjobqueue.JobQueue[PhysicsJob]
		VideoJobQueue     *sbjobqueue.JobQueue[VideoJob]
		CSVLoaderJobQueue *sbjobqueue.JobQueue[CSVLoaderJob]

		ClientCSVFileChunks     sbcsv.ChunkFileOpts
		ExerciseCSVFileChunks   sbcsv.ChunkFileOpts
		HyperparamCSVFileChunks sbcsv.ChunkFileOpts
		WorkoutCSVFileChunks    sbcsv.ChunkFileOpts

		Global GlobalConf
	}
)
