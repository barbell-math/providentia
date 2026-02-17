package barpathtracker

// #cgo CXXFLAGS: -O3 -Wall -Werror -march=native -std=c++23
// #cgo CXXFLAGS: -I../../../_deps/ffmpeg/include
// #cgo LDFLAGS: -lstdc++
// #cgo LDFLAGS: -L../../../_deps/ffmpeg/lib
// #cgo LDFLAGS: -L../../../_deps/vkSdk/1.4.341.1/x86_64/lib
// #cgo LDFLAGS: -lavfilter -lavformat -lavcodec -lavutil -lavdevice -lswscale -lswresample
// #cgo LDFLAGS: -lvulkan -lglslang
// #cgo LDFLAGS: -lpthread -pthread
// #cgo LDFLAGS: -lz -lm -ldl -llzma
// #cgo LDFLAGS: -ldrm
// #include "cpu.h"
import "C"
import (
	"fmt"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

//go:generate go-enum --marshal --names --values --nocase --noprefix

type (
	// ENUM(
	//	NoBarPathTrackerErr
	//	CouldNotAllocFrameErr
	//	CouldNotAllocPacketErr
	//	CouldNotAllocDecoderCtxErr
	//	CouldNotAllocFilterGraphErr
	//	CouldNotOpenVideoFileErr
	//	CouldNotFindInputStreamInfoErr
	//	CouldNotFindVideoStreamErr
	//	CouldNotOpenCodecForStreamErr
	//	CouldNotCreateBufferSourceErr
	//	CouldNotCreateBufferSinkErr
	//	CouldNotSetSinkPixFmtErr
	//	CouldNotInitBufferSinkErr
	//	CouldNotParseFilterErr
	//	CouldNotConfigureFilterErr
	//	CouldNotAddFrameToFilterGraphErr
	//	CouldNotGetFrameFromFilterGraphErr
	//
	//	VulkanNotSupportedErr
	//	DecoderDoesNotSupportVulkanErr
	//	AVCodecParametersToCtxtErr
	//	CouldNotCreateHwDeviceErr
	//	CouldNotReadFrameErr
	//	CouldNotSendPacketErr
	//	CouldNotReceiveFrameErr
	//	CouldNotTransferDataFromGPUToCPUErr
	// )
	BarPathTrackerErrCode int64

	CData struct{}
)

func Calc(
	rawData *types.PhysicsData,
) error {
	err := C.CalcBarPathTrackerData()
	fmt.Println("Back in the go code...", BarPathTrackerErrCode(err))

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
