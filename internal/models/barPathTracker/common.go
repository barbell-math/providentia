package barpathtracker

// #cgo CXXFLAGS: -O3 -Wall -Werror -march=native -std=c++23
// #cgo CXXFLAGS: -I../../../_deps/ffmpeg/include
// #cgo LDFLAGS: -L../../../_deps/ffmpeg/lib
// #cgo LDFLAGS: -L../../../_deps/vkSdk/lib
// #cgo LDFLAGS: -lavfilter -lavformat -lavcodec -lavutil -lavdevice -lswscale -lswresample
// #cgo LDFLAGS: -lglslang -lSPIRV -lOSDependent -lMachineIndependent -lGenericCodeGen -lglslang-default-resource-limits -lSPIRV-Tools -lshaderc_combined
// #cgo LDFLAGS: -lstdc++ -pthread -lpthread
// #cgo LDFLAGS: -lz -lm -llzma -ldrm
// #include "cpu.h"
import "C"
import (
	"fmt"
	"image"
	"image/png"
	"os"
	"unsafe"

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

//export goSaveImage
func goSaveImage(data *C.uchar, width C.int, height C.int) {
	fmt.Println("IN GO CODE: ", width, height)
	var name string
	fmt.Scanln(&name)
	size := width * height
	imgData := (*[1 << 30]uint8)(unsafe.Pointer(data))[:size:size]

	img := &image.Gray{
		Pix:    imgData,
		Stride: int(width),
		Rect:   image.Rect(0, 0, int(width), int(height)),
	}
	outfile, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	err = png.Encode(outfile, img)
	if err != nil {
		panic(err)
	}
}

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
