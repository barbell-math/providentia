package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Returns a [types.BarPathVariant] initialized with a video path as the data
// source.
func BarPathVideo(videoPath string) types.BarPathVariant {
	return types.BarPathVariant{
		Flag:      types.VideoBarPathData,
		VideoPath: videoPath,
	}
}

// Returns a [types.BarPathVariant] initialized with time series data as the data
// source.
func BarPathTimeSeriesData(data types.RawTimeSeriesData) types.BarPathVariant {
	return types.BarPathVariant{
		Flag:       types.TimeSeriesBarPathData,
		TimeSeries: data,
	}
}

// Calculates the physics data for the supplied exercise using the supplied
// raw data. The `Weight`, `Sets`, and `Reps` fields of exercise data must be
// populated with accurate values. The `PhysicsData` field will be populated
// with the results. The length of the raw data must match the number of sets.
//
// If an error occurs the state of the `PhysicsData` field in the supplied
// `exerciseData` struct will not be deterministic and should not be used. All
// other fields will remain untouched.
func CalcPhysicsData(
	ctxt context.Context,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	exerciseData *types.ExerciseData,
	rawData ...types.BarPathVariant,
) (opErr error) {
	if exerciseData == nil || len(rawData) == 0 {
		return
	}
	return runOp(ctxt, jobs.RunPhysicsJobs, jobs.PhysicsOpts{
		BarPathCalcParams:    barPathCalcParams,
		BarTrackerCalcParams: barTrackerCalcParams,
		RawData:              rawData,
		ExerciseData:         exerciseData,
	})
}
