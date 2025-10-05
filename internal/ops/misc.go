package ops

import (
	"context"
	"io/fs"
	"os"
	"path"

	"code.barbellmath.net/barbell-math/providentia/internal/db"
	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
)

func RunMigrations(ctxt context.Context, state *types.State) (opErr error) {
	if opErr = sbsqlm.Load(
		db.SqlMigrations, "migrations", db.PostOps,
	); opErr != nil {
		return
	}
	if opErr = sbsqlm.Run(ctxt, state.DB); opErr != nil {
		return
	}

	return
}

func UploadCSVDataDir(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	opts *types.CSVDataDirOptions,
) (opErr error) {
	batch, _ := sbjobqueue.BatchWithContext(ctxt)

	var clientFiles []fs.DirEntry
	var exerciseFiles []fs.DirEntry
	var hyperparamFiles []fs.DirEntry
	if clientFiles, opErr = os.ReadDir(opts.ClientDir); opErr != nil {
		return
	}
	if exerciseFiles, opErr = os.ReadDir(opts.ExerciseDir); opErr != nil {
		return
	}
	if hyperparamFiles, opErr = os.ReadDir(opts.HyperparamsDir); opErr != nil {
		return
	}

	clientCreator := CreateClients
	if opts.ClientCreateType == types.EnsureExists {
		clientCreator = EnsureClientsExist
	}
	for _, f := range clientFiles {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != "csv" {
			continue
		}
		fullPath := path.Join(opts.ClientDir, f.Name())
		state.Log.Log(
			ctxt, sblog.VLevel(3), "Uploading client CSV data", "File", fullPath,
		)
		state.GPJobQueue.Schedule(&jobs.GP[struct{}]{
			S: state,
			Q: queries,
			B: batch,
			F: func(
				ctxt context.Context,
				state *types.State,
				queries *dal.SyncQueries,
				_ struct{},
			) error {
				return UploadClientsFromCSV(
					ctxt, state, queries,
					clientCreator, opts.Opts, fullPath,
				)
			},
		})
	}

	exerciseCreator := CreateExercises
	if opts.ExerciseCreateType == types.EnsureExists {
		exerciseCreator = EnsureExercisesExist
	}
	for _, f := range exerciseFiles {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != "csv" {
			continue
		}
		fullPath := path.Join(opts.ExerciseDir, f.Name())
		state.Log.Log(
			ctxt, sblog.VLevel(3), "Uploading exercise CSV data", "File", fullPath,
		)
		state.GPJobQueue.Schedule(&jobs.GP[struct{}]{
			S: state,
			Q: queries,
			B: batch,
			F: func(
				ctxt context.Context,
				state *types.State,
				queries *dal.SyncQueries,
				_ struct{},
			) error {
				return UploadExercisesFromCSV(
					ctxt, state, queries,
					exerciseCreator, opts.Opts, fullPath,
				)
			},
		})
	}

	for _, f := range hyperparamFiles {
		if f.IsDir() {
			continue
		}
	}

	if opErr = batch.Wait(); opErr != nil {
		return
	}

	// TODO - Upload workouts
	// workoutCreator := CreateWorkouts
	// if opts.WorkoutCreateType == types.EnsureExists {
	// 	workoutCreator = EnsureWorkoutsExist
	// }

	return
}
