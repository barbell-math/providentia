//go:build !prov_gpu

package simplifiednegativespace

// #cgo CXXFLAGS: -O3 -march=native -std=c++23 -I../../../_deps/eigen
// #cgo LDFLAGS: -lstdc++
// #include "cpu.h"
import "C"
import (
	"unsafe"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/barbell-math/providentia/lib/types"
)

//go:generate go run ./structGen/structGen.go

type (
	Opts struct {
		Alpha    float64
		Beta     float64
		Gamma    float64
		MaxIters uint64
	}
)

// TODO - TEST THIS SHIT

func ModelStates(
	clientID int64,
	data []dal.ClientTrainingLogDataDateRangeAscendingRow,
	startCalcsIdx int64,
	opts Opts,
) []dal.BulkCreateModelStatesParams {
	if len(data) == 0 {
		return []dal.BulkCreateModelStatesParams{}
	}

	modelStates := make([]dal.BulkCreateModelStatesParams, len(data))
	var dataPntr *C.trainingLog_t
	var modelStatesPntr *C.modelState_t
	if len(data) > 0 {
		dataPntr = (*C.trainingLog_t)(unsafe.Pointer(&data[0]))
	}
	if len(modelStates) > 0 {
		modelStatesPntr = (*C.modelState_t)(unsafe.Pointer(&modelStates[0]))
	}

	C.calcModelStates(
		C.long(clientID),
		C.int32_t(types.SimplifiedNegativeSpace),
		dataPntr, C.long(len(data)),
		C.int64_t(startCalcsIdx),
		modelStatesPntr, C.long(len(modelStates)),
		(*C.opts_t)(unsafe.Pointer(&opts)),
	)
	return modelStates
}
