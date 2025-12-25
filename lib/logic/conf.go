package logic

import (
	"context"
	"flag"
	"fmt"
	"math"
	"runtime"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
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
			BatchSize: 1e3,
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
		CSVLoaderJobQueue: sbjobqueue.Opts{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
		GPJobQueue: sbjobqueue.Opts{
			QueueLen:       10,
			MaxNumWorkers:  uint32(runtime.NumCPU()),
			MaxJobsPerPoll: 1,
		},
		ClientCSVFileChunks: sbcsv.ChunkFileOpts{
			PredictedAvgRowSizeInBytes: 100,
			MinChunkRows:               50000,
			MaxChunkRows:               math.MaxInt,
			RequestedNumChunks:         runtime.NumCPU(),
		},
		ExerciseCSVFileChunks: sbcsv.ChunkFileOpts{
			PredictedAvgRowSizeInBytes: 100,
			MinChunkRows:               50000,
			MaxChunkRows:               math.MaxInt,
			RequestedNumChunks:         runtime.NumCPU(),
		},
		HyperparamCSVFileChunks: sbcsv.ChunkFileOpts{
			NumRowSamples:      2,
			MinChunkRows:       50000,
			MaxChunkRows:       math.MaxInt,
			RequestedNumChunks: runtime.NumCPU(),
		},
		WorkoutCSVFileChunks: sbcsv.ChunkFileOpts{
			NumRowSamples: 2,
			// MinChunkRows is set kinda low so that other job queues (physics,
			// and video mainly) have a greater chance of being filled up. Can
			// help loading data with sparse physics data, which will be the
			// more typical use case, but won't speed up loading data with dense
			// physics data. - TODO - what???
			MinChunkRows:       1e2,
			MaxChunkRows:       math.MaxInt,
			RequestedNumChunks: runtime.NumCPU(),
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
//   - <longArgStart>.PhysicsData.TimeDeltaEps
//   - <longArgStart>.PhysicsJobQueue.QueueLen
//   - <longArgStart>.PhysicsJobQueue.MaxNumWorkers
//   - <longArgStart>.PhysicsJobQueue.MaxJobsPerPoll
//   - <longArgStart>.VideoJobQueue.QueueLen
//   - <longArgStart>.VideoJobQueue.MaxNumWorkers
//   - <longArgStart>.VideoJobQueue.MaxJobsPerPoll
//   - <longArgStart>.CSVLoaderJobQueue.QueueLen
//   - <longArgStart>.CSVLoaderJobQueue.MaxNumWorkers
//   - <longArgStart>.CSVLoaderJobQueue.MaxJobsPerPoll
//   - <longArgStart>.GPJobQueue.QueueLen
//   - <longArgStart>.GPJobQueue.MaxNumWorkers
//   - <longArgStart>.GPJobQueue.MaxJobsPerPoll
//   - <longArgStart>.ClientCSVFileChunks.PredictedAvgRowSizeInBytes
//   - <longArgStart>.ClientCSVFileChunks.NumRowSamples
//   - <longArgStart>.ClientCSVFileChunks.MinChunkRows
//   - <longArgStart>.ClientCSVFileChunks.MaxChunkRows
//   - <longArgStart>.ClientCSVFileChunks.RequestedNumChunks
//   - <longArgStart>.ExerciseCSVFileChunks.PredictedAvgRowSizeInBytes
//   - <longArgStart>.ExerciseCSVFileChunks.NumRowSamples
//   - <longArgStart>.ExerciseCSVFileChunks.MinChunkRows
//   - <longArgStart>.ExerciseCSVFileChunks.MaxChunkRows
//   - <longArgStart>.ExerciseCSVFileChunks.RequestedNumChunks
//   - <longArgStart>.HyperparamCSVFileChunks.PredictedAvgRowSizeInBytes
//   - <longArgStart>.HyperparamCSVFileChunks.NumRowSamples
//   - <longArgStart>.HyperparamCSVFileChunks.MinChunkRows
//   - <longArgStart>.HyperparamCSVFileChunks.MaxChunkRows
//   - <longArgStart>.HyperparamCSVFileChunks.RequestedNumChunks
//   - <longArgStart>.WorkoutCSVFileChunks.PredictedAvgRowSizeInBytes
//   - <longArgStart>.WorkoutCSVFileChunks.NumRowSamples
//   - <longArgStart>.WorkoutCSVFileChunks.MinChunkRows
//   - <longArgStart>.WorkoutCSVFileChunks.MaxChunkRows
//   - <longArgStart>.WorkoutCSVFileChunks.RequestedNumChunks
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

	jobQueueArguments(
		fs, startStr, "Physics",
		&val.PhysicsJobQueue, &_default.PhysicsJobQueue,
	)
	jobQueueArguments(
		fs, startStr, "Video",
		&val.VideoJobQueue, &_default.VideoJobQueue,
	)
	jobQueueArguments(
		fs, startStr, "CSVLoader",
		&val.CSVLoaderJobQueue, &_default.CSVLoaderJobQueue,
	)
	jobQueueArguments(
		fs, startStr, "GPJobQueue",
		&val.GPJobQueue, &_default.GPJobQueue,
	)

	csvFileChunks(
		fs, startStr, "Client",
		&val.ClientCSVFileChunks, &_default.ClientCSVFileChunks,
	)
	csvFileChunks(
		fs, startStr, "Exercise",
		&val.ExerciseCSVFileChunks, &_default.ExerciseCSVFileChunks,
	)
	csvFileChunks(
		fs, startStr, "Hyperparam",
		&val.HyperparamCSVFileChunks, &_default.HyperparamCSVFileChunks,
	)
	csvFileChunks(
		fs, startStr, "Workout",
		&val.WorkoutCSVFileChunks, &_default.WorkoutCSVFileChunks,
	)
}

func jobQueueArguments(
	fs *flag.FlagSet,
	startStr func(names ...string) string,
	queueName string,
	val *sbjobqueue.Opts,
	_default *sbjobqueue.Opts,
) {
	fs.Func(
		startStr(fmt.Sprintf("%sJobQueue", queueName), "QueueLen"),
		fmt.Sprintf(
			"The maximum queue length for the %s job queue. Once reached adding to the queue will be a blocking operation",
			queueName,
		),
		sbargp.Uint(
			&val.QueueLen,
			_default.QueueLen,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sJobQueue", queueName), "MaxNumWorkers"),
		fmt.Sprintf(
			"The maximum number of workers for the %s job queue",
			queueName,
		),
		sbargp.Uint(
			&val.MaxNumWorkers,
			_default.MaxNumWorkers,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sJobQueue", queueName), "MaxJobsPerPoll"),
		fmt.Sprintf(
			"The maximum number of jobs the %s job queue can run each time it is polled",
			queueName,
		),
		sbargp.Uint(
			&val.MaxJobsPerPoll,
			_default.MaxJobsPerPoll,
			10,
		),
	)
}

func csvFileChunks(
	fs *flag.FlagSet,
	startStr func(names ...string) string,
	fileType string,
	val *sbcsv.ChunkFileOpts,
	_default *sbcsv.ChunkFileOpts,
) {
	fs.Func(
		startStr(fmt.Sprintf("%sCSVFileChunks", fileType), "PredictedAvgRowSizeInBytes"),
		"Used to approximate row count for each chunk. If >0 the supplied value is used, if <=0 `NumRowSamples` random rows are used to determine the predicted average row size",
		sbargp.Int(
			&val.PredictedAvgRowSizeInBytes,
			_default.PredictedAvgRowSizeInBytes,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sCSVFileChunks", fileType), "NumRowSamples"),
		"If `PredictedAvgRowSizeInBytes` <=0 `NumRowSamples` will be used to calculate the predicted avg row size",
		sbargp.Int(
			&val.NumRowSamples,
			_default.NumRowSamples,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sCSVFileChunks", fileType), "MinChunkRows"),
		"The minimum allowed number of rows in any chunk",
		sbargp.Int(
			&val.MinChunkRows,
			_default.MinChunkRows,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sCSVFileChunks", fileType), "MaxChunkRows"),
		"The maximum allowed number of rows in any chunk",
		sbargp.Int(
			&val.MaxChunkRows,
			_default.MaxChunkRows,
			10,
		),
	)
	fs.Func(
		startStr(fmt.Sprintf("%sCSVFileChunks", fileType), "RequestedNumChunks"),
		"The requested number of chunks to split the file into. The actual number of chunks may vary depending on the `MinChunkRows` and `MaxChunkRows` values. If >0 the supplied value is used, if <=0 it will be set to [runtime.NumCPU]",
		sbargp.Int(
			&val.RequestedNumChunks,
			_default.RequestedNumChunks,
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
	state.ClientCSVFileChunks = c.ClientCSVFileChunks
	state.ExerciseCSVFileChunks = c.ExerciseCSVFileChunks
	state.HyperparamCSVFileChunks = c.HyperparamCSVFileChunks
	state.WorkoutCSVFileChunks = c.WorkoutCSVFileChunks

	if err = state.ClientCSVFileChunks.Validate(); err != nil {
		return
	}
	if err = state.ExerciseCSVFileChunks.Validate(); err != nil {
		return
	}
	if err = state.HyperparamCSVFileChunks.Validate(); err != nil {
		return
	}
	if err = state.WorkoutCSVFileChunks.Validate(); err != nil {
		return
	}

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
	if state.CSVLoaderJobQueue, err = sbjobqueue.NewJobQueue[types.CSVLoaderJob](
		&c.CSVLoaderJobQueue,
	); err != nil {
		return
	}
	if state.GPJobQueue, err = sbjobqueue.NewJobQueue[types.GeneralPurposeJob](
		&c.GPJobQueue,
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
	poolConf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &dal.PgxLogAdapter{Logger: state.Log},
		LogLevel: tracelog.LogLevelDebug,
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
