package types

import (
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

type (
	// Global settings that configure many parts of providentia's behavior.
	GlobalConf struct {
		BatchSize uint
	}

	// Holds all configuration data for the library. Used to define the state of
	// the library and provides associated utility functions for cmd line
	// argument parsing such as [logic.ConfParser] and [logic.ConfDefaults].
	Conf struct {
		Global GlobalConf

		Logging sbargp.LoggingConf
		DB      sbargp.DBConf

		// Configuration that is used when setting up the physics job queue. The
		// physics job queue is responsible for taking position time series data
		// that represents the bar path and calculating all other values such as
		// velocity, acceleration, etc
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		PhysicsJobQueue sbjobqueue.Opts
		// Configuration that is used when setting up the video job queue. The
		// video job queue is responsible for taking position a video and
		// producing time series data that represents the bar path.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		VideoJobQueue sbjobqueue.Opts
		// Configuration that is used when setting up the csv loader job queue.
		// The csv loader job queue is responsible for loading data from csv
		// files on disk, verifying that data, and uploading it to the database.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		CSVLoaderJobQueue sbjobqueue.Opts

		// Configuration for how client csv files get chunked up to allow the
		// chunks to be processed in parallel by the CSVLoaderJobQueue.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-csv#ChunkFileOpts
		ClientCSVFileChunks sbcsv.ChunkFileOpts
		// Configuration for how exercise csv files get chunked up to allow the
		// chunks to be processed in parallel by the CSVLoaderJobQueue.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-csv#ChunkFileOpts
		ExerciseCSVFileChunks sbcsv.ChunkFileOpts
		// Configuration for how hyperparam csv files get chunked up to allow
		// the chunks to be processed in parallel by the CSVLoaderJobQueue.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-csv#ChunkFileOpts
		HyperparamCSVFileChunks sbcsv.ChunkFileOpts
		// Configuration for how workout csv files get chunked up to allow the
		// chunks to be processed in parallel by the CSVLoaderJobQueue.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-csv#ChunkFileOpts
		WorkoutCSVFileChunks sbcsv.ChunkFileOpts
	}
)
