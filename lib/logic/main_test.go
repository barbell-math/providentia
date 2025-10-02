package logic

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbargp "code.barbellmath.net/barbell-math/smoothbrain-argparse"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TestDBPool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	dbPswd, ok := os.LookupEnv("DB_PSWD")
	if !ok {
		panic("Set DB_PSWD env var")
	}

	var err error
	var poolConf *pgxpool.Config
	if poolConf, err = pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "postgres", dbPswd, "postgres",
	)); err != nil {
		panic(err)
	}
	if TestDBPool, err = pgxpool.NewWithConfig(
		context.Background(), poolConf,
	); err != nil {
		panic(err)
	}
	if err = TestDBPool.Ping(context.Background()); err != nil {
		panic(err)
	}
	m.Run()
	TestDBPool.Close()
}

func resetApp(ctxtIn context.Context) (context.Context, func()) {
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS provlib_tests WITH (FORCE);")
	_, err := TestDBPool.Exec(ctxtIn, sql)
	if err != nil {
		panic(err)
	}

	sql = fmt.Sprintf("CREATE DATABASE provlib_tests;")
	_, err = TestDBPool.Exec(ctxtIn, sql)
	if err != nil {
		panic(err)
	}

	return testAppMain(ctxtIn, "-conf", "../../bs/testsDB.toml")
}

// This is meant to represent the main function of a separate application.
//
// Normally no arguments would be passed in and no parameters would be returned
// but because this test application sits beneath the testing framework it has
// to accept parameters and return values.
//
// Every time you see `testCtxt` it could be replaced with context.Background()
// in a real application.
//
// `args` would be gathered from os.Args.
func testAppMain(
	testCtxt context.Context,
	args ...string,
) (context.Context, func()) {
	var conf types.Conf
	if err := sbargp.Parse(&conf, args, sbargp.ParserOpts[types.Conf]{
		ProgName: "testApp",
		RequiredArgs: []string{
			"DB.User", "DB.PswdEnvVar", "DB.Name",
		},
		ArgDefsSetter: func(conf *types.Conf, fs *flag.FlagSet) error {
			ConfParser(fs, conf, "", ConfDefaults())
			return nil
		},
	}); err != nil {
		panic(err)
	}

	// Normally this would be derived from context.Background()
	appLifetime, appCancel := context.WithCancel(testCtxt)
	state, err := ConfToState(appLifetime, &conf)
	if err != nil {
		panic(err)
	}
	go sbjobqueue.Poll(
		appLifetime,
		state.PhysicsJobQueue, state.VideoJobQueue, state.CSVLoaderJobQueue,
	)

	// Notice how cancellation can be derived from a parent context. This allows
	// separate lib calls to be canceled separately while still allowing a
	// single parent context to cancel all lib calls at once.
	provLifetime := WithStateValue(appLifetime, state)
	if err := RunMigrations(provLifetime); err != nil {
		panic(err)
	}

	// Normally there would be a defer function call rather than returning a
	// function
	return provLifetime, func() {
		appCancel()
		CleanupState(state)
	}
}
