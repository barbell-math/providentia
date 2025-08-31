package barpathphysdata

// #cgo CXXFLAGS: -O3 -Werror -march=native -std=c++23 -I../../../_deps/eigen  -I../../glue
// #cgo LDFLAGS: -lstdc++
// #include "cpu.h"
import "C"
import (
	"errors"
	"math"
	"runtime"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
)

//go:generate go-enum --marshal --names --values --nocase --noprefix

type (
	// ENUM(
	//	NoErr
	//	TimeSeriesNotIncreasingErr
	//	TimeSeriesNotMonotonicErr
	//	InvalidApproximationErrErr
	// )
	BarPathCalcErrCode int64

	Data struct {
		Mass    types.Kilogram
		timeLen int64
		time    *types.Second
		pos     *types.Vec2[types.Meter]
		vel     *types.Vec2[types.MeterPerSec]
		acc     *types.Vec2[types.MeterPerSec2]
		jerk    *types.Vec2[types.MeterPerSec3]
		force   *types.Vec2[types.Newton]
		impulse *types.Vec2[types.NewtonSec]
		power   *types.Watt
		work    *types.Joule

		Reps     int32
		repSplit *types.Split[types.Second]
	}
)

var (
	InvalidRawDataIdxErr = errors.New("Invalid raw data index")
	InvalidRawDataLenErr = errors.New("Invalid raw data length")
)

func Calc(
	state *types.State,
	tl *dal.BulkCreateTrainingLogsParams,
	rawData *dal.CreatePhysicsDataParams,
	idx int,
) error {
	ceilSets := math.Ceil(tl.Sets)
	floorSets := math.Floor(tl.Sets)
	if idx >= int(ceilSets) {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) >= the ceiling of the number of sets (%f)",
			idx, ceilSets,
		)
	}
	if tl.Reps <= 0 {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Supplied data must have at least 1 rep",
		)
	}

	if len(rawData.Time) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed time range: [0, %d)",
			idx, len(rawData.Time),
		)
	}
	if len(rawData.Position) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed position range: [0, %d)",
			idx, len(rawData.Position),
		)
	}

	expLen := len(rawData.Time[idx])
	if expLen < int(state.PhysicsData.MinNumSamples) {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"the minimum number of samples (%d) was not provided, got %d samples",
			state.PhysicsData.MinNumSamples, expLen,
		)
	}
	if len(rawData.Position[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected position slice of len %d, got len %d",
			expLen, len(rawData.Position[idx]),
		)
	}

	expReps := tl.Reps
	if ceilSets > tl.Sets && int(floorSets) == idx {
		expReps = int32((tl.Sets - floorSets) * float64(tl.Reps))
	}

	rawData.Velocity[idx] = make([]types.Vec2[types.MeterPerSec], expLen)
	rawData.Acceleration[idx] = make([]types.Vec2[types.MeterPerSec2], expLen)
	rawData.Jerk[idx] = make([]types.Vec2[types.MeterPerSec3], expLen)
	rawData.Impulse[idx] = make([]types.Vec2[types.NewtonSec], expLen)
	rawData.Force[idx] = make([]types.Vec2[types.Newton], expLen)
	rawData.Work[idx] = make([]types.Joule, expLen)
	rawData.Power[idx] = make([]types.Watt, expLen)
	rawData.RepSplits[idx] = make([]types.Split[types.Second], expReps)

	// Note:
	// Checks for monotonically increasing time series data are done in the
	// [C.calcBarPathPhysData] func because those checks can be performance
	// intensive operations.

	baseData := Data{
		timeLen:  int64(len(rawData.Time[idx])),
		Mass:     tl.Weight,
		time:     &rawData.Time[idx][0],
		pos:      &rawData.Position[idx][0],
		vel:      &rawData.Velocity[idx][0],
		acc:      &rawData.Acceleration[idx][0],
		jerk:     &rawData.Jerk[idx][0],
		force:    &rawData.Force[idx][0],
		impulse:  &rawData.Impulse[idx][0],
		power:    &rawData.Power[idx][0],
		work:     &rawData.Work[idx][0],
		Reps:     expReps,
		repSplit: &rawData.RepSplits[idx][0],
	}

	pinner := runtime.Pinner{}
	pinner.Pin(baseData.time)
	pinner.Pin(baseData.pos)
	pinner.Pin(baseData.vel)
	pinner.Pin(baseData.acc)
	pinner.Pin(baseData.jerk)
	pinner.Pin(baseData.force)
	pinner.Pin(baseData.impulse)
	pinner.Pin(baseData.power)
	pinner.Pin(baseData.work)
	pinner.Pin(baseData.repSplit)

	err := C.calcBarPathPhysData(
		(*C.barPathData_t)(unsafe.Pointer(&baseData)),
		(*C.barPathCalcConf_t)(unsafe.Pointer(&state.BarPathCalc)),
		(*C.physDataConf_t)(unsafe.Pointer(&state.PhysicsData)),
	)

	pinner.Unpin()

	switch BarPathCalcErrCode(err) {
	case TimeSeriesNotIncreasingErr:
		return sberr.Wrap(
			types.TimeSeriesDecreaseErr,
			"Time samples must be increasing",
		)
	case TimeSeriesNotMonotonicErr:
		return sberr.Wrap(
			types.TimeSeriesNotMonotonicErr,
			"Time samples must all have the same delta (within %f variance)",
			state.PhysicsData.TimeDeltaEps,
		)
	case InvalidApproximationErrErr:
		return types.ErrInvalidApproximationError
	}

	return nil
}
