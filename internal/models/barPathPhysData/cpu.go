package barpathphysdata

// #cgo CXXFLAGS: -O3 -Wall -Werror -march=native -std=c++23 -I../../../_deps/eigen  -I../../clib
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

var (
	InvalidRawDataIdxErr = errors.New("Invalid raw data index")
	InvalidRawDataLenErr = errors.New("Invalid raw data length")
)

func InitBarPathCalcPhysicsData(
	rawData *dal.CreatePhysicsDataParams,
	barPathCalcParams *types.BarPathCalcHyperparams,
	numSets int,
) {
	*rawData = dal.CreatePhysicsDataParams{
		BarPathCalcParamsVersion: barPathCalcParams.Version,
		Time:                     make([][]types.Second, numSets),
		Position:                 make([][]types.Vec2[types.Meter, types.Meter], numSets),
		Velocity:                 make([][]types.Vec2[types.MeterPerSec, types.MeterPerSec], numSets),
		Acceleration:             make([][]types.Vec2[types.MeterPerSec2, types.MeterPerSec2], numSets),
		Jerk:                     make([][]types.Vec2[types.MeterPerSec3, types.MeterPerSec3], numSets),
		Force:                    make([][]types.Vec2[types.Newton, types.Newton], numSets),
		Impulse:                  make([][]types.Vec2[types.NewtonSec, types.NewtonSec], numSets),
		Work:                     make([][]types.Joule, numSets),
		Power:                    make([][]types.Watt, numSets),
		RepSplits:                make([][]types.Split, numSets),
		MinVel:                   make([][]types.PointInTime[types.Second, types.MeterPerSec], numSets),
		MaxVel:                   make([][]types.PointInTime[types.Second, types.MeterPerSec], numSets),
		MinAcc:                   make([][]types.PointInTime[types.Second, types.MeterPerSec2], numSets),
		MaxAcc:                   make([][]types.PointInTime[types.Second, types.MeterPerSec2], numSets),
		MinForce:                 make([][]types.PointInTime[types.Second, types.Newton], numSets),
		MaxForce:                 make([][]types.PointInTime[types.Second, types.Newton], numSets),
		MinImpulse:               make([][]types.PointInTime[types.Second, types.NewtonSec], numSets),
		MaxImpulse:               make([][]types.PointInTime[types.Second, types.NewtonSec], numSets),
		AvgWork:                  make([][]types.Joule, numSets),
		MinWork:                  make([][]types.PointInTime[types.Second, types.Joule], numSets),
		MaxWork:                  make([][]types.PointInTime[types.Second, types.Joule], numSets),
		AvgPower:                 make([][]types.Watt, numSets),
		MinPower:                 make([][]types.PointInTime[types.Second, types.Watt], numSets),
		MaxPower:                 make([][]types.PointInTime[types.Second, types.Watt], numSets),
	}
}

func Calc(
	tl *dal.BulkCreateTrainingLogsParams,
	rawData *dal.CreatePhysicsDataParams,
	barPathCalcParams *types.BarPathCalcHyperparams,
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
	if expLen < int(barPathCalcParams.MinNumSamples) {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"the minimum number of samples (%d) was not provided, got %d samples",
			barPathCalcParams.MinNumSamples, expLen,
		)
	}
	if len(rawData.Position[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected position slice of len %d, got len %d",
			expLen, len(rawData.Position[idx]),
		)
	}

	// Note:
	// Checks for monotonically increasing time series data are done in the
	// [C.calcBarPathPhysData] func because those checks can be performance
	// intensive operations.

	expReps := tl.Reps
	if ceilSets > tl.Sets && int(floorSets) == idx {
		expReps = max(int32((tl.Sets-floorSets)*float64(tl.Reps)), 1)
	}

	rawData.Velocity[idx] = make([]types.Vec2[types.MeterPerSec, types.MeterPerSec], expLen)
	rawData.Acceleration[idx] = make([]types.Vec2[types.MeterPerSec2, types.MeterPerSec2], expLen)
	rawData.Jerk[idx] = make([]types.Vec2[types.MeterPerSec3, types.MeterPerSec3], expLen)
	rawData.Impulse[idx] = make([]types.Vec2[types.NewtonSec, types.NewtonSec], expLen)
	rawData.Force[idx] = make([]types.Vec2[types.Newton, types.Newton], expLen)
	rawData.Work[idx] = make([]types.Joule, expLen)
	rawData.Power[idx] = make([]types.Watt, expLen)
	rawData.RepSplits[idx] = make([]types.Split, expReps)
	rawData.MinVel[idx] = make([]types.PointInTime[types.Second, types.MeterPerSec], expReps)
	rawData.MaxVel[idx] = make([]types.PointInTime[types.Second, types.MeterPerSec], expReps)
	rawData.MinAcc[idx] = make([]types.PointInTime[types.Second, types.MeterPerSec2], expReps)
	rawData.MaxAcc[idx] = make([]types.PointInTime[types.Second, types.MeterPerSec2], expReps)
	rawData.MinForce[idx] = make([]types.PointInTime[types.Second, types.Newton], expReps)
	rawData.MaxForce[idx] = make([]types.PointInTime[types.Second, types.Newton], expReps)
	rawData.MinImpulse[idx] = make([]types.PointInTime[types.Second, types.NewtonSec], expReps)
	rawData.MaxImpulse[idx] = make([]types.PointInTime[types.Second, types.NewtonSec], expReps)
	rawData.AvgWork[idx] = make([]types.Joule, expReps)
	rawData.MinWork[idx] = make([]types.PointInTime[types.Second, types.Joule], expReps)
	rawData.MaxWork[idx] = make([]types.PointInTime[types.Second, types.Joule], expReps)
	rawData.AvgPower[idx] = make([]types.Watt, expReps)
	rawData.MinPower[idx] = make([]types.PointInTime[types.Second, types.Watt], expReps)
	rawData.MaxPower[idx] = make([]types.PointInTime[types.Second, types.Watt], expReps)

	baseData := Data{
		timeLen:    int64(len(rawData.Time[idx])),
		mass:       tl.Weight,
		time:       &rawData.Time[idx][0],
		pos:        &rawData.Position[idx][0],
		vel:        &rawData.Velocity[idx][0],
		acc:        &rawData.Acceleration[idx][0],
		jerk:       &rawData.Jerk[idx][0],
		force:      &rawData.Force[idx][0],
		impulse:    &rawData.Impulse[idx][0],
		power:      &rawData.Power[idx][0],
		work:       &rawData.Work[idx][0],
		reps:       expReps,
		repSplit:   &rawData.RepSplits[idx][0],
		minVel:     &rawData.MinVel[idx][0],
		maxVel:     &rawData.MaxVel[idx][0],
		minAcc:     &rawData.MinAcc[idx][0],
		maxAcc:     &rawData.MaxAcc[idx][0],
		minForce:   &rawData.MinForce[idx][0],
		maxForce:   &rawData.MaxForce[idx][0],
		minImpulse: &rawData.MinImpulse[idx][0],
		maxImpulse: &rawData.MaxImpulse[idx][0],
		avgWork:    &rawData.AvgWork[idx][0],
		minWork:    &rawData.MinWork[idx][0],
		maxWork:    &rawData.MaxWork[idx][0],
		avgPower:   &rawData.AvgPower[idx][0],
		minPower:   &rawData.MinPower[idx][0],
		maxPower:   &rawData.MaxPower[idx][0],
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

	err := C.calcBarPathPhysData(
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
