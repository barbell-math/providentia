package barpathphysdata

// #cgo CXXFLAGS: -O3 -Wall -Werror -march=native -std=c++23 -I../../../_deps/eigen  -I../../clib
// #cgo LDFLAGS: -lstdc++
// #include "cpu.h"
import "C"
import (
	"runtime"
	"unsafe"

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

	CData struct {
		mass    types.Kilogram
		timeLen int64
		time    *types.Second
		pos     *types.Vec2[types.Meter, types.Meter]
		vel     *types.Vec2[types.MeterPerSec, types.MeterPerSec]
		acc     *types.Vec2[types.MeterPerSec2, types.MeterPerSec2]
		jerk    *types.Vec2[types.MeterPerSec3, types.MeterPerSec3]
		force   *types.Vec2[types.Newton, types.Newton]
		impulse *types.Vec2[types.NewtonSec, types.NewtonSec]
		power   *types.Watt
		work    *types.Joule

		reps     int32
		repSplit *types.Split

		minVel     *types.PointInTime[types.Second, types.MeterPerSec]
		maxVel     *types.PointInTime[types.Second, types.MeterPerSec]
		minAcc     *types.PointInTime[types.Second, types.MeterPerSec2]
		maxAcc     *types.PointInTime[types.Second, types.MeterPerSec2]
		minForce   *types.PointInTime[types.Second, types.Newton]
		maxForce   *types.PointInTime[types.Second, types.Newton]
		minImpulse *types.PointInTime[types.Second, types.NewtonSec]
		maxImpulse *types.PointInTime[types.Second, types.NewtonSec]
		avgWork    *types.Joule
		minWork    *types.PointInTime[types.Second, types.Joule]
		maxWork    *types.PointInTime[types.Second, types.Joule]
		avgPower   *types.Watt
		minPower   *types.PointInTime[types.Second, types.Watt]
		maxPower   *types.PointInTime[types.Second, types.Watt]
	}
)

func Calc(
	rawData *types.PhysicsData,
	barPathCalcParams *types.BarPathCalcHyperparams,
	weight types.Kilogram,
	expNumReps int32,
) error {
	expLen := len(rawData.Time)
	if expLen < int(barPathCalcParams.MinNumSamples) {
		return sberr.Wrap(
			types.InvalidRawDataLenErr,
			"The minimum number of samples (%d) was not provided, got %d samples",
			barPathCalcParams.MinNumSamples, expLen,
		)
	}
	if len(rawData.Position) != expLen {
		return sberr.Wrap(
			types.InvalidRawDataLenErr,
			"Expected position slice of len %d, got len %d",
			expLen, len(rawData.Position),
		)
	}
	if expNumReps <= 0 {
		return sberr.Wrap(
			types.InvalidExpNumRepsErr,
			"Must be >=0. Got: %d", expNumReps,
		)
	}

	// Note:
	// Checks for monotonically increasing time series data are done in the
	// [C.calcBarPathPhysData] func because those checks can be performance
	// intensive operations.

	rawData.BarPathCalcVersion = barPathCalcParams.Version
	// TODO - only alloc if necessary?
	rawData.Velocity = make([]types.Vec2[types.MeterPerSec, types.MeterPerSec], expLen)
	rawData.Acceleration = make([]types.Vec2[types.MeterPerSec2, types.MeterPerSec2], expLen)
	rawData.Jerk = make([]types.Vec2[types.MeterPerSec3, types.MeterPerSec3], expLen)
	rawData.Impulse = make([]types.Vec2[types.NewtonSec, types.NewtonSec], expLen)
	rawData.Force = make([]types.Vec2[types.Newton, types.Newton], expLen)
	rawData.Work = make([]types.Joule, expLen)
	rawData.Power = make([]types.Watt, expLen)
	rawData.RepSplits = make([]types.Split, expNumReps)
	rawData.MinVel = make([]types.PointInTime[types.Second, types.MeterPerSec], expNumReps)
	rawData.MaxVel = make([]types.PointInTime[types.Second, types.MeterPerSec], expNumReps)
	rawData.MinAcc = make([]types.PointInTime[types.Second, types.MeterPerSec2], expNumReps)
	rawData.MaxAcc = make([]types.PointInTime[types.Second, types.MeterPerSec2], expNumReps)
	rawData.MinForce = make([]types.PointInTime[types.Second, types.Newton], expNumReps)
	rawData.MaxForce = make([]types.PointInTime[types.Second, types.Newton], expNumReps)
	rawData.MinImpulse = make([]types.PointInTime[types.Second, types.NewtonSec], expNumReps)
	rawData.MaxImpulse = make([]types.PointInTime[types.Second, types.NewtonSec], expNumReps)
	rawData.AvgWork = make([]types.Joule, expNumReps)
	rawData.MinWork = make([]types.PointInTime[types.Second, types.Joule], expNumReps)
	rawData.MaxWork = make([]types.PointInTime[types.Second, types.Joule], expNumReps)
	rawData.AvgPower = make([]types.Watt, expNumReps)
	rawData.MinPower = make([]types.PointInTime[types.Second, types.Watt], expNumReps)
	rawData.MaxPower = make([]types.PointInTime[types.Second, types.Watt], expNumReps)

	baseData := CData{
		timeLen:    int64(len(rawData.Time)),
		mass:       weight,
		time:       &rawData.Time[0],
		pos:        &rawData.Position[0],
		vel:        &rawData.Velocity[0],
		acc:        &rawData.Acceleration[0],
		jerk:       &rawData.Jerk[0],
		force:      &rawData.Force[0],
		impulse:    &rawData.Impulse[0],
		power:      &rawData.Power[0],
		work:       &rawData.Work[0],
		reps:       expNumReps,
		repSplit:   &rawData.RepSplits[0],
		minVel:     &rawData.MinVel[0],
		maxVel:     &rawData.MaxVel[0],
		minAcc:     &rawData.MinAcc[0],
		maxAcc:     &rawData.MaxAcc[0],
		minForce:   &rawData.MinForce[0],
		maxForce:   &rawData.MaxForce[0],
		minImpulse: &rawData.MinImpulse[0],
		maxImpulse: &rawData.MaxImpulse[0],
		avgWork:    &rawData.AvgWork[0],
		minWork:    &rawData.MinWork[0],
		maxWork:    &rawData.MaxWork[0],
		avgPower:   &rawData.AvgPower[0],
		minPower:   &rawData.MinPower[0],
		maxPower:   &rawData.MaxPower[0],
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
	pinner.Pin(baseData.minVel)
	pinner.Pin(baseData.maxVel)
	pinner.Pin(baseData.minAcc)
	pinner.Pin(baseData.maxAcc)
	pinner.Pin(baseData.minForce)
	pinner.Pin(baseData.maxForce)
	pinner.Pin(baseData.minImpulse)
	pinner.Pin(baseData.maxImpulse)
	pinner.Pin(baseData.avgWork)
	pinner.Pin(baseData.minWork)
	pinner.Pin(baseData.maxWork)
	pinner.Pin(baseData.avgPower)
	pinner.Pin(baseData.minPower)
	pinner.Pin(baseData.maxPower)

	err := C.CalcBarPathPhysData(
		(*C.barPathData_t)(unsafe.Pointer(&baseData)),
		(*C.barPathCalcHyperparams_t)(unsafe.Pointer(barPathCalcParams)),
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
			"Adjacent time samples must all have the same delta (within %f variance)",
			barPathCalcParams.TimeDeltaEps,
		)
	case InvalidApproximationErrErr:
		return types.ErrInvalidApproximationError
	}

	return nil
}
