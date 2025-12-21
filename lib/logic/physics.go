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

func CalcPhysicsData(
	ctxt context.Context,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	exerciseData []types.ExerciseData,
	rawData ...types.BarPathVariant,
) (opErr error) {
	if len(exerciseData) == 0 {
		return
	}
	return runOp(ctxt, jobs.RunPhysicsJobs, jobs.PhysicsOpts{})
}
