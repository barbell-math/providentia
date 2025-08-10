package types

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"runtime"

	sbargp "github.com/barbell-math/smoothbrain-argparse"
	sberr "github.com/barbell-math/smoothbrain-errs"
	sblog "github.com/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
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
		Physics PhysicsConf
		Global  GlobalConf
	}

	// Configuration that the rest of providentia will use. Setup with care as
	// the values have the ability to control how providentia behaves.
	GlobalConf struct {
		NumWorkers uint
		BatchSize  uint
		// SimplifiedNegativeSpaceModel simplifiednegativespace.Opts
	}

	// Configuration that is used when parsing, generating, and utilizing
	// physics data.
	PhysicsConf struct {
		MinNumSamples uint
		TimeDeltaEps  float64
	}

	// The state the rest of providentia will use. Almost all functions
	// available for external use from this library will require this state to
	// be available in the passed in context.
	State struct {
		Physics PhysicsConf
		Global  GlobalConf
		DB      *pgxpool.Pool
		Log     *slog.Logger
	}

	ctxtKey struct{}
)

var (
	stateCtxtKey ctxtKey
)

// Returns a [State] from the supplied context or nil if it was not present. The
// boolean flag indicates if the [State] value was present.
func FromContext(ctxt context.Context) (*State, bool) {
	s, ok := ctxt.Value(stateCtxtKey).(*State)
	return s, ok
}

// Adds the supplied state value to the supplied context, returning a new
// context with the state value and a function that will clean up the given
// state. This function should be called when the state is no longer needed, and
// is intended to be used with a defer stmt.
//
// This function should be used when you want the same database connection and
// logger from your application code to be used in the providentia library.
// Note that you must provide a valid state or you may get undesired behavior
// from providentia lib.
func WithValue(ctxt context.Context, s *State) (context.Context, func()) {
	return context.WithValue(ctxt, stateCtxtKey, s), func() {
		if s.DB != nil {
			s.DB.Close()
		}
	}
}

// Parses a set of arguments from the supplied args (which usually come from the
// cmd line) as well as an optional TOML file to create the application state.
// The TOML file that will be parsed will be defined by the `-conf` cmd line
// argument if it is present. The returned context will have the generated state
// value and the returned function will clean up the given state. This function
// is intended to be used with a defer stmt.
//
// This function should be used when you want a different database connection
// and a different logger to be used in the providentia library than the ones in
// your application code.
func Parse(ctxt context.Context, args []string) (context.Context, func(), error) {
	var err error
	var poolConf *pgxpool.Config
	var _fullConf fullConf
	state := State{}

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
				&conf.Global.NumWorkers,
				"Global.NumWorkers",
				uint(runtime.NumCPU()),
				"The number of worker threads a single library function can use",
			)
			fs.UintVar(
				&conf.Global.BatchSize,
				"Global.BatchSize",
				1e6,
				"The batch size the library functions will work with. Smaller will use less memory but may be slightly slower",
			)

			fs.UintVar(
				&conf.Physics.MinNumSamples,
				"Physics.MinNumSamples",
				20,
				"The minimum number of samples that should be present in physics data",
			)
			fs.Float64Var(
				&conf.Physics.TimeDeltaEps,
				"Physics.TimeDeltaEps",
				1e6,
				"The maximum acceptable variance between time sample deltas",
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

	if _fullConf.Physics.MinNumSamples < 2 {
		err = sberr.Wrap(
			InvalidMinNumSamplesErr,
			"Must be >=2. Got: %d", _fullConf.Physics.MinNumSamples,
		)
		goto done
	}

	state.Physics = _fullConf.Physics
	state.Global = _fullConf.Global

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
	var doneFunc func() = func() {}
	if err == nil {
		newCtxt, doneFunc = WithValue(ctxt, &state)
	}
	return newCtxt, doneFunc, err
}
