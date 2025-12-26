package jobs

import (
	"context"
	"iter"
	"maps"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/jackc/pgx/v5"
)

func resolveCreateFunc[T dal.AvailableTypes](
	_type types.CreateFuncType,
	create dal.CreateFunc[T],
	ensure dal.CreateFunc[T],
) dal.CreateFunc[T] {
	if _type == types.EnsureExists {
		return ensure
	}
	return create
}

func getFilesInDirFunc(dir string) iter.Seq2[string, error] {
	return util.FilesWithExtInDir(dir, ".csv", util.FilesWithExtInDirOpts{})
}

func hyperparamFilterFunc(extension string) func(v *string, e *error) bool {
	return func(v *string, e *error) bool {
		if *e != nil {
			return true
		}

		split := strings.SplitN(*v, ".", 3)
		if len(split) != 3 {
			*e = sberr.Wrap(
				types.UnknownFileInDataDirErr,
				"Invalid hyperparam file. File name must have the following format: <file name>.<hyperparam type>.csv\nGot: %s",
				*v,
			)
			return true
		}
		if split[1] == extension {
			return true
		}
		if _, ok := types.HyperparamFileNames[split[1]]; !ok {
			*e = sberr.Wrap(
				types.UnknownFileInDataDirErr,
				"Unknown hyperparam type: %s, must be one of %v. File: %s",
				split[1], maps.Keys(types.HyperparamFileNames), *v,
			)
			return true
		}
		return false
	}
}

func BulkUploadData(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *types.BulkUploadDataOpts,
) error {
	if opts.WorkoutDir != "" && opts.BarPathCalcHyperparams == nil {
		return sberr.Wrap(
			types.BulkDataUploadErr,
			"When a workout data dir is supplied BarPathCalcHyperparams must not be nil",
		)
	}
	if opts.WorkoutDir != "" && opts.BarPathTrackerHyperparams == nil {
		return sberr.Wrap(
			types.BulkDataUploadErr,
			"When a workout data dir is supplied BarPathTrackerHyperparams must not be nil",
		)
	}

	batch, _ := sbjobqueue.BatchWithContext(ctxt)

	if err := UploadFromCSV(ctxt, state, tx, &CSVLoaderOpts[types.Client]{
		Opts:  &opts.Opts,
		Files: getFilesInDirFunc(opts.ClientDir),
		Batch: batch,
		Creator: resolveCreateFunc(
			opts.ClientCreateType, dal.CreateClients, dal.EnsureClientsExist,
		),
	}); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	if err := UploadFromCSV(ctxt, state, tx, &CSVLoaderOpts[types.Exercise]{
		Opts:  &opts.Opts,
		Files: getFilesInDirFunc(opts.ExerciseDir),
		Batch: batch,
		Creator: resolveCreateFunc(
			opts.ExerciseCreateType, dal.CreateExercises, dal.EnsureExercisesExist,
		),
	}); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	if err := UploadFromCSV(
		ctxt, state, tx, &CSVLoaderOpts[types.BarPathCalcHyperparams]{
			Opts: &opts.Opts,
			Files: util.FilterSeq2Err(
				getFilesInDirFunc(opts.HyperparamsDir),
				hyperparamFilterFunc(types.BarPathCalcFileExt),
			),
			Batch: batch,
			Creator: resolveCreateFunc(
				opts.HyperparamsCreateType,
				dal.CreateHyperparams[types.BarPathCalcHyperparams],
				dal.EnsureHyperparamsExist[types.BarPathCalcHyperparams],
			),
		},
	); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	if err := UploadFromCSV(
		ctxt, state, tx, &CSVLoaderOpts[types.BarPathTrackerHyperparams]{
			Opts: &opts.Opts,
			Files: util.FilterSeq2Err(
				getFilesInDirFunc(opts.HyperparamsDir),
				hyperparamFilterFunc(types.BarPathTrackerFileExt),
			),
			Batch: batch,
			Creator: resolveCreateFunc(
				opts.HyperparamsCreateType,
				dal.CreateHyperparams[types.BarPathTrackerHyperparams],
				dal.EnsureHyperparamsExist[types.BarPathTrackerHyperparams],
			),
		},
	); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	if err := batch.Wait(); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	if err := UploadWorkoutsFromCSV(ctxt, state, tx, &CSVWorkoutLoaderOpts{
		Opts:                      &opts.Opts,
		Files:                     getFilesInDirFunc(opts.WorkoutDir),
		Batch:                     batch,
		BarPathCalcHyperparams:    opts.BarPathCalcHyperparams,
		BarPathTrackerHyperparams: opts.BarPathTrackerHyperparams,
	}); err != nil {
		return sberr.AppendError(types.BulkDataUploadErr, err)
	}

	return batch.Wait()
}
