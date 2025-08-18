package logic

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	// Configuraiton that is used when setting up the physics job queue. The
	// physics job queue is responsible for taking position time series data
	// that represents the bar path and calculating all other values such as
	// velocity, acceleration, etc
	PhysicsJobQueueConf sbjobqueue.Opts
	// Configuraiton that is used when setting up the video job queue. The video
	// job queue is responsible for taking position a video and producing time
	// series data that represents the bar path.
	VideoJobQueueConf sbjobqueue.Opts

	// This is the full conf struct that will be populated by parsing the CMD
	// line args and TOML file in [Parse]. This is a superset of the [Conf]
	// struct so that the [State] struct will never be able to have conflicting
	// information between the generated values such as the db conn and the
	// DBConf args.
	// All fields that are used to do things like set up secondary variables in
	// the [State] struct should not be included in the [Conf] struct but should
	// be included here.
	fullConf struct {
		Logging sbargp.LoggingConf
		DB      sbargp.DBConf

		Global      types.GlobalConf
		PhysicsData types.PhysicsDataConf

		PhysicsJobQueue PhysicsJobQueueConf
		VideoJobQueue   VideoJobQueueConf
	}

	ctxtKey struct{}
)

var (
	stateCtxtKey ctxtKey
)

// Returns a [State] from the supplied context or nil if it was not present. The
// boolean flag indicates if the [State] value was present.
func StateFromContext(ctxt context.Context) (*types.State, bool) {
	s, ok := ctxt.Value(stateCtxtKey).(*types.State)
	return s, ok
}

// Adds the supplied state value to the supplied context, returning a new
// context with the state value and a cancel function that will clean up the
// state. This function should be called when the state is no longer needed, and
// is intended to be used with a defer stmt.
//
// The [WithStateValue] function should be used when you want the same database
// connection and logger from your application code to be used in the
// providentia library. Note that you must provide a valid state or you may get
// undesired behavior from providentia lib.
func WithStateValue(
	ctxt context.Context,
	s *types.State,
) (context.Context, func(), error) {
	newCtxt, cancel := setStateValue(ctxt, s)
	return newCtxt, cancel, validateState(s)
}

// Parses a set of arguments from the supplied args (which usually come from the
// cmd line) as well as an optional TOML file to create the application state.
// The TOML file that will be parsed will be defined by the `--conf` cmd line
// argument if it is present. The returned context will have the generated state
// value and the returned function will clean up the given state. This function
// is intended to be used with a defer stmt.
//
// This function should be used when you want a different database connection
// and a different logger to be used in the providentia library than the ones in
// your application code.
func ParseState(
	ctxt context.Context,
	args []string,
) (context.Context, func(), error) {
	var err error
	var poolConf *pgxpool.Config
	var _fullConf fullConf
	state := types.State{}

	if err = sbargp.Parse(&_fullConf, args, sbargp.ParserOpts[fullConf]{
		ProgName: "providentia",
		RequiredArgs: []string{
			"DB.User", "DB.PswdEnvVar", "DB.Name",
		},
		ArgDefsSetter: func(conf *fullConf, fs *flag.FlagSet) error {
			sbargp.Logging(fs, &conf.Logging, "Logging", sbargp.LoggingConf{
				Verbosity:       0,
				SaveTo:          "",
				Name:            "providentia",
				MaxNumLogs:      1,
				MaxLogSizeBytes: 1e6, // 1 MB
			})
			sbargp.DB(fs, &conf.DB, "DB", sbargp.DBConf{
				Host: "localhost",
				Port: 5432,
			})
			fs.UintVar(
				&conf.Global.BatchSize,
				"Global.BatchSize",
				1e6,
				"The batch size the library functions will work with. Smaller will use less memory but may be slightly slower",
			)

			fs.UintVar(
				&conf.PhysicsData.MinNumSamples,
				"PhysicsData.MinNumSamples",
				20,
				"The minimum number of samples that should be present in physics data",
			)
			fs.Float64Var(
				&conf.PhysicsData.TimeDeltaEps,
				"PhysicsData.TimeDeltaEps",
				1e-6,
				"The maximum acceptable variance between time sample deltas",
			)

			fs.Uint64Var(
				&conf.PhysicsJobQueue.QueueLen,
				"PhysicsJobQueue.QueueLen",
				10,
				"The maximum queue length for the physics job queue. Once reached adding to the queue will be a blocking operation",
			)
			fs.Uint64Var(
				&conf.PhysicsJobQueue.MaxNumWorkers,
				"PhysicsJobQueue.MaxNumWorkers",
				uint64(runtime.NumCPU()),
				"The maximum number of workers for the physics job queue",
			)
			fs.Uint64Var(
				&conf.PhysicsJobQueue.MaxJobsPerPoll,
				"PhysicsJobQueue.MaxJobsPerPoll",
				1,
				"The maximum number of jobs the physics job queue can run each time it is polled",
			)

			fs.Uint64Var(
				&conf.VideoJobQueue.QueueLen,
				"VideoJobQueue.QueueLen",
				10,
				"The maximum queue length for the video job queue. Once reached adding to the queue will be a blocking operation",
			)
			fs.Uint64Var(
				&conf.VideoJobQueue.MaxNumWorkers,
				"VideoJobQueue.MaxNumWorkers",
				uint64(runtime.NumCPU()),
				"The maximum number of workers for the video job queue",
			)
			fs.Uint64Var(
				&conf.VideoJobQueue.MaxJobsPerPoll,
				"VideoJobQueue.MaxJobsPerPoll",
				1,
				"The maximum number of jobs the video job queue can start each time it is polled",
			)

			// fs.Float64Var(
			// 	&conf.SimplifiedNegativeSpaceModel.Alpha,
			// 	"SimplifiedNegativeSpaceModel.Alpha",
			// 	1,
			// 	"The value to use for alpha in the simplified negative space model",
			// )
			// fs.Float64Var(
			// 	&conf.SimplifiedNegativeSpaceModel.Beta,
			// 	"SimplifiedNegativeSpaceModel.Beta",
			// 	2,
			// 	"The value to use for beta in the simplified negative space model",
			// )
			// fs.Float64Var(
			// 	&conf.SimplifiedNegativeSpaceModel.Gamma,
			// 	"SimplifiedNegativeSpaceModel.Gamma",
			// 	2,
			// 	"The value to use for gamma in the simplified negative space model",
			// )
			// fs.Uint64Var(
			// 	&conf.SimplifiedNegativeSpaceModel.MaxIters,
			// 	"SimplifiedNegativeSpaceModel.MaxIters",
			// 	1e6,
			// 	"The maximum number of iterations that each model state can go through before exiting",
			// )
			return nil
		},
	}); err != nil {
		goto done
	}

	state.PhysicsData = _fullConf.PhysicsData
	state.Global = _fullConf.Global
	if err = validateState(&state); err != nil {
		goto done
	}

	if state.PhysicsJobQueue, err = sbjobqueue.NewJobQueue[types.PhysicsJob](
		(*sbjobqueue.Opts)(&_fullConf.PhysicsJobQueue),
	); err != nil {
		goto done
	}
	if state.VideoJobQueue, err = sbjobqueue.NewJobQueue[types.VideoJob](
		(*sbjobqueue.Opts)(&_fullConf.VideoJobQueue),
	); err != nil {
		goto done
	}

	if state.Log, err = sblog.New(sblog.Opts{
		CurVerbosityLevel: uint(_fullConf.Logging.Verbosity),
		RotateWriterOpts: sblog.RotateWriterOpts{
			LogDir:          string(_fullConf.Logging.SaveTo),
			LogName:         _fullConf.Logging.Name,
			MaxNumLogs:      uint(_fullConf.Logging.MaxNumLogs),
			MaxLogSizeBytes: uint64(_fullConf.Logging.MaxLogSizeBytes),
		},
	}); err != nil {
		goto done
	}

	if poolConf, err = pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		_fullConf.DB.Host,
		_fullConf.DB.Port,
		_fullConf.DB.User,
		_fullConf.DB.PswdEnvVar,
		_fullConf.DB.Name,
	)); err != nil {
		goto done
	}
	if state.DB, err = pgxpool.NewWithConfig(ctxt, poolConf); err != nil {
		goto done
	}
	if err = state.DB.Ping(ctxt); err != nil {
		goto done
	}

done:
	var newCtxt context.Context
	var cancelFunc func() = func() {}
	if err == nil {
		// validateState has already been called at this point
		newCtxt, cancelFunc = setStateValue(ctxt, &state)
	}
	return newCtxt, cancelFunc, err
}

func setStateValue(ctxt context.Context, s *types.State) (context.Context, func()) {
	return context.WithValue(ctxt, stateCtxtKey, s), func() {
		// Should this be a cancel as well?
		// No - cancelation will be handeled by caller not here, see :AppSetup
		// for an example
		if s.DB != nil {
			s.DB.Close()
		}
	}
}

func validateState(s *types.State) error {
	if s.PhysicsData.MinNumSamples < 2 {
		return sberr.AppendError(
			types.InvalidPhysicsDataConfErr,
			sberr.Wrap(
				types.InvalidMinNumSamplesErr,
				"Must be >=2. Got: %d", s.PhysicsData.MinNumSamples,
			),
		)
	}
	return nil
}
