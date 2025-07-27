package logic

import (
	"context"
	"fmt"

	"github.com/barbell-math/providentia/lib/types"
)

func resetDB(ctxtIn context.Context) (context.Context, func()) {
	ctxt, stateCleanup, err := types.Parse(
		ctxtIn,
		[]string{"-conf", "../../bs/testSetup.toml"},
	)
	defer stateCleanup()
	if err != nil {
		panic(err)
	}
	dropDatabase(ctxt)
	addDatabase(ctxt)
	return initDatabase(ctxt)
}

func dropDatabase(ctxt context.Context) {
	state, ok := types.FromContext(ctxt)
	if ok != true {
		panic("Could not find state in context!")
	}

	sql := fmt.Sprintf("DROP DATABASE IF EXISTS provlib_tests WITH (FORCE);")
	_, err := state.DB.Exec(ctxt, sql)
	state.Log.Info("Dropped", "sql", sql)
	if err != nil {
		panic(err)
	}
}

func addDatabase(ctxt context.Context) {
	state, ok := types.FromContext(ctxt)
	if ok != true {
		panic("Could not find state in context!")
	}

	sql := fmt.Sprintf("CREATE DATABASE provlib_tests;")
	_, err := state.DB.Exec(ctxt, sql)
	state.Log.Info("Created", "sql", sql)
	if err != nil {
		panic(err)
	}
}

func initDatabase(ctxt context.Context) (context.Context, func()) {
	state, ok := types.FromContext(ctxt)
	if ok != true {
		panic("Could not find state in context!")
	}

	state.Log.Info("Setting up tests database...")
	testCtxt, stateCleanup, err := types.Parse(
		ctxt,
		[]string{
			"-conf", "../../bs/testsDB.toml",
			"--Logging.Name", "setup",
		},
	)
	if err != nil {
		panic(err)
	}
	if err := RunMigrations(testCtxt); err != nil {
		panic(err)
	}

	state.Log.Info("Setting up tests database... done.")

	return testCtxt, stateCleanup
}
