package barpathtracker

// #cgo CXXFLAGS: -O3 -Wall -Werror -march=native -std=c++23
// #cgo CXXFLAGS: -I../../../_deps/ffmpeg/include
// #cgo LDFLAGS: -lstdc++
// #cgo LDFLAGS: -L../../../_deps/ffmpeg/lib
// #cgo LDFLAGS: -lavformat -lavcodec -lavutil -lavdevice -lswscale -lswresample
// #cgo LDFLAGS: -lpthread -pthread
// #cgo LDFLAGS: -lz -lm -ldl -llzma
// #cgo LDFLAGS: -lva -lva-drm -lva-x11
// #cgo LDFLAGS: -ldrm -lX11
// #include "cpu.h"
import "C"
import "code.barbellmath.net/barbell-math/providentia/lib/types"

//go:generate go-enum --marshal --names --values --nocase --noprefix

type (
	// ENUM(
	//	NoBarPathTrackerErr
	//	VAAPINotSupportedErr
	//	CouldNotAllocateAVPacketErr
	//	CouldNotOpenVideoFileErr
	//	CouldNotFindInputStreamInfoErr
	//	CouldNotFindVideoStreamErr
	//	DecoderDoesNotSupportVAAPIErr
	//	AVCodecParametersToCtxtErr
	//	CouldNotCreateHwDeviceErr
	//	CouldNotOpenCodecForStreamErr
	// )
	BarPathTrackerErrCode int64

	CData struct{}
)

func Calc(
	rawData *types.PhysicsData,
) error {
	err := C.CalcBarPathTrackerData()

	switch BarPathTrackerErrCode(err) {
	// case TimeSeriesNotIncreasingErr:
	// 	return sberr.Wrap(
	// 		types.TimeSeriesDecreaseErr,
	// 		"Time samples must be increasing",
	// 	)
	// case TimeSeriesNotMonotonicErr:
	// 	return sberr.Wrap(
	// 		types.TimeSeriesNotMonotonicErr,
	// 		"Adjacent time samples must all have the same delta (within %f variance)",
	// 		barPathCalcParams.TimeDeltaEps,
	// 	)
	// case InvalidApproximationErrErr:
	// 	return types.ErrInvalidApproximationError
	}

	return nil
}
