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
