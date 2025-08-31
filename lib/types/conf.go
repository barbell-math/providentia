package types

import (
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/maypok86/otter/v2"
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
		MinNumSamples uint64
		TimeDeltaEps  Second
	}

	BarPathCalcConf struct {
		ApproxErr       ApproximationError
		SmootherWeight1 float64
		SmootherWeight2 float64
		SmootherWeight3 float64
		SmootherWeight4 float64
		SmootherWeight5 float64
	}

	// Holds all configuration data for the library. Used to define the state of
	// the library and provides associated utility functions for cmd line
	// argument parsing such as [logic.ConfParser] and [logic.ConfDefaults].
	Conf struct {
		Logging sbargp.LoggingConf
		DB      sbargp.DBConf

		Global      GlobalConf
		PhysicsData PhysicsDataConf
		BarPathCalc BarPathCalcConf

		// Configuraiton that is used when setting up the physics job queue. The
		// physics job queue is responsible for taking position time series data
		// that represents the bar path and calculating all other values such as
		// velocity, acceleration, etc
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		PhysicsJobQueue sbjobqueue.Opts
		// Configuraiton that is used when setting up the video job queue. The
		// video job queue is responsible for taking position a video and
		// producing time series data that represents the bar path.
		// Refer to: http://code.barbellmath.net/barbell-math/smoothbrain-jobQueue#Opts
		VideoJobQueue sbjobqueue.Opts

		// Configuration that is used when setting up the in memory client
		// cache. This cache will be used to reduce database queries at expense
		// of a larger memory footprint.
		// Refer to: https://pkg.go.dev/github.com/maypok86/otter/v2#Options
		// Note: not all conf options can be set through cmd line parameters.
		ClientCache otter.Options[string, IdWrapper[int64, Client]]
		// Configuration that is used when setting up the in memory exercise
		// cache. This cache will be used to reduce database queries at expense
		// of a larger memory footprint.
		// Refer to: https://pkg.go.dev/github.com/maypok86/otter/v2#Options
		// Note: not all conf options can be set through cmd line parameters.
		ExerciseCache otter.Options[string, IdWrapper[int32, Exercise]]
	}
)
