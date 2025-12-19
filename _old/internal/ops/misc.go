package ops

import (
	"context"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/internal/db"
	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
)

var hyperparamFileNames = []string{
	"barPathCalc",
	"barPathTracker",
}

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
	if opts.WorkoutDir != "" && opts.BarPathCalcHyperparams == nil {
		opErr = sberr.Wrap(
			types.InvalidDataDirOptsErr,
			"If a workout data dir is supplied the bar path calc hyperparams must not be nil",
		)
		return
	}
	if opts.WorkoutDir != "" && opts.BarPathTrackerHyperparams == nil {
		opErr = sberr.Wrap(
			types.InvalidDataDirOptsErr,
			"If a workout data dir is supplied the bar path tracker hyperparams must not be nil",
		)
		return
	}

	batch, _ := sbjobqueue.BatchWithContext(ctxt)

	if opts.ClientDir != "" {
		var clientFiles []fs.DirEntry
		if clientFiles, opErr = os.ReadDir(opts.ClientDir); opErr != nil {
			opErr = sberr.AppendError(types.DataDirErr, opErr)
			return
		}
		clientCreator := CreateClients
		if opts.ClientCreateType == types.EnsureExists {
			clientCreator = EnsureClientsExist
		}

		for _, f := range clientFiles {
			if f.IsDir() {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(types.DirInDataDirErr, "Dir: %s", f.Name()),
				)
				return
			}
			if path.Ext(f.Name()) != ".csv" {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"All client files must end in 'csv'. File: %s", f.Name(),
					),
				)
				return
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
	}

	if opts.ExerciseDir != "" {
		var exerciseFiles []fs.DirEntry
		if exerciseFiles, opErr = os.ReadDir(opts.ExerciseDir); opErr != nil {
			return
		}
		exerciseCreator := CreateExercises
		if opts.ExerciseCreateType == types.EnsureExists {
			exerciseCreator = EnsureExercisesExist
		}

		for _, f := range exerciseFiles {
			if f.IsDir() {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(types.DirInDataDirErr, "Dir: %s", f.Name()),
				)
				return
			}
			if path.Ext(f.Name()) != ".csv" {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"All exercise files must end in 'csv'. File: %s", f.Name(),
					),
				)
				return
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
	}

	if opts.HyperparamsDir != "" {
		var hyperparamFiles []fs.DirEntry
		if hyperparamFiles, opErr = os.ReadDir(opts.HyperparamsDir); opErr != nil {
			return
		}
		creators := NewHyperparamCreators(opts.HyperparamsCreateType)

		for _, f := range hyperparamFiles {
			if f.IsDir() {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(types.DirInDataDirErr, "Dir: %s", f.Name()),
				)
				return
			}
			split := strings.SplitN(f.Name(), ".", 3)
			if len(split) != 3 {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"Invalid hyperparam file. File name must have the following format: <file name>.<hyperparam type>.csv\nGot: %s",
						f.Name(),
					),
				)
				return
			}
			if split[2] != "csv" {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"All hyperparam files must end in 'csv'. File: %s",
						f.Name(),
					),
				)
				return
			}
			if !slices.Contains(hyperparamFileNames, split[1]) {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"Unknown hyperparam type: %s, must be one of %v. File: %s",
						split[1], hyperparamFileNames, f.Name(),
					),
				)
				return
			}

			fullPath := path.Join(opts.HyperparamsDir, f.Name())
			state.Log.Log(
				ctxt, sblog.VLevel(3), "Uploading hyperparam CSV data",
				"File", fullPath,
				"Type", split[1],
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
					switch split[1] {
					case "barPathCalc":
						return UploadHyperparamsFromCSV(
							ctxt, state, queries,
							creators.barPathCalc, opts.Opts, fullPath,
						)
					case "barPathTracker":
						return UploadHyperparamsFromCSV(
							ctxt, state, queries,
							creators.barPathTracker, opts.Opts, fullPath,
						)
					}
					return nil
				},
			})
		}
	}

	// Have to wait because the workout data may reference
	// clients/exercises/hyperparams that were just created.
	if opErr = batch.Wait(); opErr != nil {
		return
	}

	if opts.WorkoutDir != "" {
		var workoutFiles []fs.DirEntry
		if workoutFiles, opErr = os.ReadDir(opts.WorkoutDir); opErr != nil {
			return
		}
		workoutCreator := CreateWorkouts
		if opts.WorkoutCreateType == types.EnsureExists {
			workoutCreator = EnsureWorkoutsExist
		}

		for _, f := range workoutFiles {
			if f.IsDir() {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(types.DirInDataDirErr, "Dir: %s", f.Name()),
				)
				return
			}
			if path.Ext(f.Name()) != ".csv" {
				opErr = sberr.AppendError(
					types.DataDirErr,
					sberr.Wrap(
						types.UnknownFileErr,
						"All workout files must end in 'csv'. File: %s",
						f.Name(),
					),
				)
				return
			}
			fullPath := path.Join(opts.WorkoutDir, f.Name())
			state.Log.Log(
				ctxt, sblog.VLevel(3), "Uploading workout CSV data",
				"File", fullPath,
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
					return UploadWorkoutsFromCSV(
						ctxt, state, queries,
						opts.BarPathCalcHyperparams,
						opts.BarPathTrackerHyperparams,
						workoutCreator, opts.Opts, fullPath,
					)
				},
			})
		}
	}

	return batch.Wait()
}
