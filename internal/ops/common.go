package ops

import (
	"context"
	"iter"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

type (
	createFunc[
		T types.Client | types.Exercise | types.RawWorkout | types.Hyperparams,
	] func(
		ctxt context.Context,
		state *types.State,
		queries *dal.SyncQueries,
		values ...T,
	) (opErr error)

	workoutCreateFunc func(
		ctxt context.Context,
		state *types.State,
		queries *dal.SyncQueries,
		barPathCalcParams *types.BarPathCalcHyperparams,
		barTrackerCalcParams *types.BarPathTrackerHyperparams,
		values ...types.RawWorkout,
	) (opErr error)

	hyperparamCreators struct {
		barPathCalc    createFunc[types.BarPathCalcHyperparams]
		barPathTracker createFunc[types.BarPathTrackerHyperparams]
	}
)

func NewHyperparamCreators(createType types.CreateFuncType) hyperparamCreators {
	creators := hyperparamCreators{
		barPathCalc:    CreateHyperparams[types.BarPathCalcHyperparams],
		barPathTracker: CreateHyperparams[types.BarPathTrackerHyperparams],
	}
	if createType == types.EnsureExists {
		creators = hyperparamCreators{
			barPathCalc:    EnsureHyperparamsExist[types.BarPathCalcHyperparams],
			barPathTracker: EnsureHyperparamsExist[types.BarPathTrackerHyperparams],
		}
	}
	return creators
}

func batchIndexes[S ~[]E, E any](s S, step int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := 0; i < len(s); i += step {
			end := i + step
			if !yield(i, min(len(s), end)) {
				break
			}
		}
	}
}

func flatten[S ~[]E, E any](s ...S) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i := range len(s) {
			for j := range len(s[i]) {
				if !yield(j, s[i][j]) {
					break
				}
			}
		}
	}
}
