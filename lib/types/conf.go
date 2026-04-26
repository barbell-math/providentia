package types

import (
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbargparse"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbjobqueue"
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
		// physics job queue is responsible for calculating bar path position
		// time series data from videos as well as calculating all other physics
		// values such as velocity, acceleration, etc from bar path position
		// time series data.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		PhysicsJobQueue sbjobqueue.Opts
		// Configuration that is used when setting up the csv loader job queue.
		// The csv loader job queue is responsible for loading data from csv
		// files on disk, verifying that data, and uploading it to the database.
		// Refer to: http://code.barbellmath.net/carmichaeljr/smoothbrain/src/branch/main/sbjobqueue/DOCS.md#Opts
		CSVLoaderJobQueue sbjobqueue.Opts
	}
)
