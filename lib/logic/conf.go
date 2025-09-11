package logic

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Returns a [types.Conf] struct with sensible default values. Can be passed to
// [ParseConfig] as the `_default` parameter.
func ConfDefaults() *types.Conf {
	return &types.Conf{
		Logging: sbargp.LoggingConf{
			Verbosity:       0,
			SaveTo:          "",
			Name:            "providentia",
			MaxNumLogs:      1,
			MaxLogSizeBytes: 1e6, // 1 MB
		},
		DB: sbargp.DBConf{
			Host: "localhost",
			Port: 5432,
		},
		Global: types.GlobalConf{
			BatchSize:             1e3,
			PerRequestIdCacheSize: 1e2,
		},
		BarPathCalc: types.BarPathCalcConf{
			MinNumSamples:   100,
			TimeDeltaEps:    1e-6,
			ApproxErr:       types.FourthOrder,
			NearZeroFilter:  0.1,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
		},
		BarPathTracker: types.BarPathTrackerConf{
			MinLength:   5,
			MinFileSize: 5e7, // 50MB
			MaxFileSize: 5e8, // 500MB
		},
		PhysicsJobQueue: sbjobqueue.Opts{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
		VideoJobQueue: sbjobqueue.Opts{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
	}
}

// Returns a list of required arguments for the default conf configuration.
// Depending on the defaults you choose to set the list of required args may
// change.
func ConfDefaultRequiredArgs() []string {
	return []string{
		"DB.User", "DB.PswdEnvVar", "DB.Name",
	}
}

// Adds cmd line parsing arguments to the supplied flag set so that the
// configuration options for the library can be parsed from the cmd line. All
// cmd line flags will be prepended with the `longArgStart` name and will have
// a default value from the `_default` struct. The following flags will be
// added:
//   - <longArgStart>.Logging.SaveTo
//   - l  (Same as <longArgStart>.SaveTo)
//   - <longArgStart>.Logging.Name
//   - <longArgStart>.Logging.MaxNumLogs
//   - <longArgStart>.Logging.MaxLogSizeBytes
//   - <longArgStart>.Logging.Verbose
//   - v  (same as <longArgStart>.Verbose)
//   - <longArgStart>.DB.User
//   - <longArgStart>.DB.PswdEnvVar
//   - <longArgStart>.DB.Host
//   - <longArgStart>.DB.Port
//   - <longArgStart>.DB.Name
//   - <longArgStart>.Global.BatchSize
//   - <longArgStart>.Global.PerRequestIdCacheSize
//   - <longArgStart>.PhysicsData.MinNumSamples
//   - <longArgStart>.PhysicsData.TimeDeltaEps
//   - <longArgStart>.PhysicsJobQueue.QueueLen
//   - <longArgStart>.PhysicsJobQueue.MaxNumWorkers
//   - <longArgStart>.PhysicsJobQueue.MaxJobsPerPoll
//   - <longArgStart>.VideoJobQueue.QueueLen
//   - <longArgStart>.VideoJobQueue.MaxNumWorkers
//   - <longArgStart>.VideoJobQueue.MaxJobsPerPoll
//   - <longArgStart>.BarPathCalc.MinNumSamples
//   - <longArgStart>.BarPathCalc.TimeDeltaEps
//   - <longArgStart>.BarPathCalc.ApproxErr
//   - <longArgStart>.BarPathCalc.NearZeroFilter
//   - <longArgStart>.BarPathCalc.SmootherWeight1
//   - <longArgStart>.BarPathCalc.SmootherWeight2
//   - <longArgStart>.BarPathCalc.SmootherWeight3
//   - <longArgStart>.BarPathCalc.SmootherWeight4
//   - <longArgStart>.BarPathCalc.SmootherWeight5
func ConfParser(
	fs *flag.FlagSet,
	val *types.Conf,
	longArgStart string,
	_default *types.Conf,
) {
	startStr := func(names ...string) string {
		fmtedNames := strings.Join(names, ".")
		if len(longArgStart) > 0 {
			return fmt.Sprintf("%s.%s", longArgStart, fmtedNames)
		} else {
			return fmtedNames
		}
	}

	sbargp.Logging(fs, &val.Logging, startStr("Logging"), _default.Logging)
	sbargp.DB(fs, &val.DB, startStr("DB"), _default.DB)

	fs.Func(
		startStr("Global", "BatchSize"),
		"The batch size the library functions will work with. Smaller will use less memory but may be slightly slower",
		sbargp.Uint(
			&val.Global.BatchSize,
			_default.Global.BatchSize,
			10,
		),
	)
	fs.Func(
		startStr("Global", "PerRequestIdCacheSize"),
		"The maximum allowed cache size for each requests id caches. Smaller numbers will use less memory at the potential expense of more netowrk trips.",
		sbargp.Uint(
			&val.Global.PerRequestIdCacheSize,
			_default.Global.PerRequestIdCacheSize,
			10,
		),
	)

	fs.Func(
		startStr("PhysicsJobQueue", "QueueLen"),
		"The maximum queue length for the physics job queue. Once reached adding to the queue will be a blocking operation",
		sbargp.Uint(
			&val.PhysicsJobQueue.QueueLen,
			_default.PhysicsJobQueue.QueueLen,
			10,
		),
	)
	fs.Func(
		startStr("PhysicsJobQueue", "MaxNumWorkers"),
		"The maximum number of workers for the physics job queue",
		sbargp.Uint(
			&val.PhysicsJobQueue.MaxNumWorkers,
			_default.PhysicsJobQueue.MaxNumWorkers,
			10,
		),
	)
	fs.Func(
		startStr("PhysicsJobQueue", "MaxJobsPerPoll"),
		"The maximum number of jobs the physics job queue can run each time it is polled",
		sbargp.Uint(
			&val.PhysicsJobQueue.MaxJobsPerPoll,
			_default.PhysicsJobQueue.MaxJobsPerPoll,
			10,
		),
	)

	fs.Func(
		startStr("VideoJobQueue", "QueueLen"),
		"The maximum queue length for the video job queue. Once reached adding to the queue will be a blocking operation",
		sbargp.Uint(
			&val.VideoJobQueue.QueueLen,
			_default.VideoJobQueue.QueueLen,
			10,
		),
	)
	fs.Func(
		startStr("VideoJobQueue", "MaxNumWorkers"),
		"The maximum number of workers for the video job queue",
		sbargp.Uint(
			&val.VideoJobQueue.MaxNumWorkers,
			_default.VideoJobQueue.MaxNumWorkers,
			10,
		),
	)
	fs.Func(
		startStr("VideoJobQueue", "MaxJobsPerPoll"),
		"The maximum number of jobs the video job queue can run each time it is polled",
		sbargp.Uint(
			&val.VideoJobQueue.MaxJobsPerPoll,
			_default.VideoJobQueue.MaxJobsPerPoll,
			10,
		),
	)

	fs.Func(
		startStr("BarPathCalc", "MinNumSamples"),
		"The minimum number of samples that should be present in physics data",
		sbargp.Uint(
			&val.BarPathCalc.MinNumSamples,
			_default.BarPathCalc.MinNumSamples,
			10,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "TimeDeltaEps"),
		"The maximum acceptable variance between time sample deltas",
		sbargp.Float(
			&val.BarPathCalc.TimeDeltaEps,
			_default.BarPathCalc.TimeDeltaEps,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "ApproxErr"),
		fmt.Sprintf(
			"The accuracy of the approximation error. One of: %v",
			types.ApproximationErrorNames(),
		),
		sbargp.FromTextUnmarshaler(
			&val.BarPathCalc.ApproxErr,
			_default.BarPathCalc.ApproxErr,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "NearZeroFilter"),
		"How close to 0 the vertical bar position can be for it to be considered 0",
		sbargp.Float(
			&val.BarPathCalc.NearZeroFilter,
			_default.BarPathCalc.NearZeroFilter,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "SmootherWeight1"),
		"The weight of the second-left value in the weighted average smoother function",
		sbargp.Float(
			&val.BarPathCalc.SmootherWeight1,
			_default.BarPathCalc.SmootherWeight1,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "SmootherWeight2"),
		"The weight of the first-left value in the weighted average smoother function",
		sbargp.Float(
			&val.BarPathCalc.SmootherWeight2,
			_default.BarPathCalc.SmootherWeight2,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "SmootherWeight3"),
		"The weight of the central value in the weighted average smoother function",
		sbargp.Float(
			&val.BarPathCalc.SmootherWeight3,
			_default.BarPathCalc.SmootherWeight3,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "SmootherWeight4"),
		"The weight of the first-right value in the weighted average smoother function",
		sbargp.Float(
			&val.BarPathCalc.SmootherWeight4,
			_default.BarPathCalc.SmootherWeight4,
		),
	)
	fs.Func(
		startStr("BarPathCalc", "SmootherWeight5"),
		"The weight of the second-right value in the weighted average smoother function",
		sbargp.Float(
			&val.BarPathCalc.SmootherWeight5,
			_default.BarPathCalc.SmootherWeight5,
		),
	)

	fs.Func(
		startStr("BarPathTracker", "MinLength"),
		"The minimum length of a video provided for bar path analysis in seconds",
		sbargp.Float(
			&val.BarPathTracker.MinLength,
			_default.BarPathTracker.MinLength,
		),
	)
	fs.Func(
		startStr("BarPathTracker", "MinFileSize"),
		"The minimum file size of a video provided for bar path analysis in bytes",
		sbargp.Uint(
			&val.BarPathTracker.MinFileSize,
			_default.BarPathTracker.MinFileSize,
			10,
		),
	)
	fs.Func(
		startStr("BarPathTracker", "MaxFileSize"),
		"The maximum file size of a video provided for bar path analysis in bytes",
		sbargp.Uint(
			&val.BarPathTracker.MaxFileSize,
			_default.BarPathTracker.MaxFileSize,
			10,
		),
	)
}

// Takes the supplied [types.Conf] struct and translates it into a [types.State]
// struct that can be passed to all of the other library calls. The [types.Conf]
// struct holds configuration values for database connections, job queues,
// logging, etc. This function creates database connections, job queues, and
// logging structs with those configuration values.
func ConfToState(
	ctxt context.Context,
	c *types.Conf,
) (state *types.State, err error) {
	var s types.State
	state = &s

	state.Global = c.Global
	state.BarPathCalc = c.BarPathCalc
	state.BarPathTracker = c.BarPathTracker

	if state.PhysicsJobQueue, err = sbjobqueue.NewJobQueue[types.PhysicsJob](
		&c.PhysicsJobQueue,
	); err != nil {
		return
	}
	if state.VideoJobQueue, err = sbjobqueue.NewJobQueue[types.VideoJob](
		&c.VideoJobQueue,
	); err != nil {
		return
	}

	if state.Log, err = sblog.New(sblog.Opts{
		CurVerbosityLevel: uint(c.Logging.Verbosity),
		RotateWriterOpts: sblog.RotateWriterOpts{
			LogDir:          string(c.Logging.SaveTo),
			LogName:         c.Logging.Name,
			MaxNumLogs:      uint(c.Logging.MaxNumLogs),
			MaxLogSizeBytes: uint64(c.Logging.MaxLogSizeBytes),
		},
	}); err != nil {
		return
	}

	var poolConf *pgxpool.Config
	if poolConf, err = pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.PswdEnvVar, c.DB.Name,
	)); err != nil {
		return
	}
	if state.DB, err = pgxpool.NewWithConfig(ctxt, poolConf); err != nil {
		return
	}
	if err = state.DB.Ping(ctxt); err != nil {
		return
	}
	if err = ValidateState(&s); err != nil {
		return
	}
	return
}
