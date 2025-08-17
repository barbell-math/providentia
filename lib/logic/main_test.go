package logic

import (
	"context"
	"fmt"

	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

func resetApp(ctxtIn context.Context) (context.Context, func()) {
	testCtxt, stateCleanup, err := ParseState(
		ctxtIn,
		[]string{"-conf", "../../bs/testSetup.toml"},
	)
	defer stateCleanup()
	if err != nil {
		panic(err)
	}
	dropDatabase(testCtxt)
	addDatabase(testCtxt)
	return initTestState(testCtxt)
}

func dropDatabase(ctxt context.Context) {
	state, ok := StateFromContext(ctxt)
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
	state, ok := StateFromContext(ctxt)
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

func initTestState(testCtxt context.Context) (context.Context, func()) {
	testSetupState, ok := StateFromContext(testCtxt)
	if ok != true {
		panic("Could not find state in context!")
	}

	// Note - this is sort of how setup would need to occur in an application
	// as well.
	// :AppSetup - Notice how cancelation can be derived from the testCtxt
	testSetupState.Log.Info("Setting up tests database...")
	provCtxt, provCleanup, err := ParseState(
		testCtxt,
		[]string{"-conf", "../../bs/testsDB.toml"},
	)
	if err != nil {
		panic(err)
	}
	provState, ok := StateFromContext(provCtxt)
	if ok != true {
		panic("Could not find state in context!")
	}

	if err := RunMigrations(provCtxt); err != nil {
		panic(err)
	}
	pollerCtxt, cancel := context.WithCancel(testCtxt)
	go sbjobqueue.Poll(pollerCtxt, provState.PhysicsJobQueue) //, state.VideoJobQueue)

	testSetupState.Log.Info("Setting up tests database... done.")

	return provCtxt, func() {
		cancel()
		provCleanup()
	}
}
