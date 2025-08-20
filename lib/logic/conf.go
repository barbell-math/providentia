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
func ConfValDefaults() *types.Conf {
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
			BatchSize: 1e3,
		},
		PhysicsJobQueue: types.PhysicsJobQueueConf{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
		VideoJobQueue: types.VideoJobQueueConf{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
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
//   - <longArgStart>.PhysicsData.MinNumSamples
//   - <longArgStart>.PhysicsData.TimeDeltaEps
//   - <longArgStart>.PhysicsJobQueue.QueueLen
//   - <longArgStart>.PhysicsJobQueue.MaxNumWorkers
//   - <longArgStart>.PhysicsJobQueue.MaxJobsPerPoll
//   - <longArgStart>.VideoJobQueue.QueueLen
//   - <longArgStart>.VideoJobQueue.MaxNumWorkers
//   - <longArgStart>.VideoJobQueue.MaxJobsPerPoll
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
	fs.UintVar(
		&val.Global.BatchSize,
		startStr("Global", "BatchSize"),
		_default.Global.BatchSize,
		"The batch size the library functions will work with. Smaller will use less memory but may be slightly slower",
	)

	fs.Func(
		startStr("PhysicsData", "MinNumSamples"),
		"The minimum number of samples that should be present in physics data",
		sbargp.Uint(
			&val.PhysicsData.MinNumSamples,
			_default.PhysicsData.MinNumSamples,
			10,
		),
	)
	fs.Func(
		startStr("PhysicsData", "TimeDeltaEps"),
		"The maximum acceptable variance between time sample deltas",
		sbargp.Float(
			&val.PhysicsData.TimeDeltaEps,
			_default.PhysicsData.TimeDeltaEps,
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

	state.PhysicsData = c.PhysicsData
	state.Global = c.Global
	if err = validateState(state); err != nil {
		return
	}

	if state.PhysicsJobQueue, err = sbjobqueue.NewJobQueue[types.PhysicsJob](
		(*sbjobqueue.Opts)(&c.PhysicsJobQueue),
	); err != nil {
		return
	}
	if state.VideoJobQueue, err = sbjobqueue.NewJobQueue[types.VideoJob](
		(*sbjobqueue.Opts)(&c.VideoJobQueue),
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
	return
}
