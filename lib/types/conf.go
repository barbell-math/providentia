package types

import (
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

type (
	// Global settings that configure many parts of providentia's behavior.
	GlobalConf struct {
		BatchSize uint
		// SimplifiedNegativeSpaceModel simplifiednegativespace.Opts
	}

	// Configuration that is used when parsing, generating, and utilizing
	// physics data.
	PhysicsDataConf struct {
		MinNumSamples uint
		TimeDeltaEps  float64
	}

	// Configuraiton that is used when setting up the physics job queue. The
	// physics job queue is responsible for taking position time series data
	// that represents the bar path and calculating all other values such as
	// velocity, acceleration, etc
	PhysicsJobQueueConf sbjobqueue.Opts
	// Configuraiton that is used when setting up the video job queue. The video
	// job queue is responsible for taking position a video and producing time
	// series data that represents the bar path.
	VideoJobQueueConf sbjobqueue.Opts

	// Holds all configuration data for the library. Used to define the state of
	// the library and provides associated utility functions for cmd line
	// argument parsing such as [logic.ConfParser] and [logic.ConfDefaults].
	Conf struct {
		Logging sbargp.LoggingConf
		DB      sbargp.DBConf

		Global      GlobalConf
		PhysicsData PhysicsDataConf

		PhysicsJobQueue PhysicsJobQueueConf
		VideoJobQueue   VideoJobQueueConf
	}
)
