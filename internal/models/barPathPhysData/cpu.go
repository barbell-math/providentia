package barpathphysdata

// #cgo CXXFLAGS: -O3 -march=native -std=c++23 -I../../../_deps/eigen  -I../../glue
// #cgo LDFLAGS: -lstdc++
// #include "cpu.h"
import "C"
import (
	"errors"
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
)

var (
	InvalidRawDataIdxErr = errors.New("Invalid raw data index")
	InvalidRawDataLenErr = errors.New("Invalid raw data length")
)

func Calc(
	state *types.State,
	rawData *dal.CreatePhysicsDataParams,
	idx int,
) error {
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
	if len(rawData.Velocity) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed velocity range: [0, %d)",
			idx, len(rawData.Velocity),
		)
	}
	if len(rawData.Acceleration) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed acceleration range: [0, %d)",
			idx, len(rawData.Acceleration),
		)
	}
	if len(rawData.Jerk) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed jerk range: [0, %d)",
			idx, len(rawData.Jerk),
		)
	}
	if len(rawData.Work) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed work range: [0, %d)",
			idx, len(rawData.Work),
		)
	}
	if len(rawData.Impulse) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed impulse range: [0, %d)",
			idx, len(rawData.Impulse),
		)
	}
	if len(rawData.Force) <= idx {
		return sberr.Wrap(
			InvalidRawDataIdxErr,
			"Supplied index (%d) outside allowed force range: [0, %d)",
			idx, len(rawData.Force),
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
	if len(rawData.Velocity[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected velocity slice of len %d, got len %d",
			expLen, len(rawData.Velocity[idx]),
		)
	}
	if len(rawData.Acceleration[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected acceleration slice of len %d, got len %d",
			expLen, len(rawData.Acceleration[idx]),
		)
	}
	if len(rawData.Jerk[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected jerk slice of len %d, got len %d",
			expLen, len(rawData.Jerk[idx]),
		)
	}
	if len(rawData.Work[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected work slice of len %d, got len %d",
			expLen, len(rawData.Work[idx]),
		)
	}
	if len(rawData.Impulse[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected impulse slice of len %d, got len %d",
			expLen, len(rawData.Impulse[idx]),
		)
	}
	if len(rawData.Force[idx]) != expLen {
		return sberr.Wrap(
			InvalidRawDataLenErr,
			"Expected force slice of len %d, got len %d",
			expLen, len(rawData.Force[idx]),
		)
	}
	// Note:
	// Checks for monotonically increasing time series data are done in the
	// [C.calcBarPathPhysData] func because those checks can be performance
	// intensive operations.

	err := C.calcBarPathPhysData(
		C.int64_t(len(rawData.Time[idx])),
		(*C.double)(unsafe.SliceData(rawData.Time[idx])),
		(*C.posVec2_t)(unsafe.Pointer(&rawData.Position[idx][0])),
		(*C.velVec2_t)(unsafe.Pointer(&rawData.Velocity[idx][0])),
		(*C.accVec2_t)(unsafe.Pointer(&rawData.Acceleration[idx][0])),
		(*C.jerkVec2_t)(unsafe.Pointer(&rawData.Jerk[idx][0])),
		(*C.workVec2_t)(unsafe.Pointer(&rawData.Work[idx][0])),
		(*C.impulseVec2_t)(unsafe.Pointer(&rawData.Impulse[idx][0])),
		(*C.forceVec2_t)(unsafe.Pointer(&rawData.Force[idx][0])),
		(*C.barPathCalcConf_t)(unsafe.Pointer(&state.BarPathCalc)),
		(*C.physDataConf_t)(unsafe.Pointer(&state.PhysicsData)),
	)

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
