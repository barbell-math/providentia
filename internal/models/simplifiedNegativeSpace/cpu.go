//go:build !prov_gpu

package simplifiednegativespace

// #cgo CXXFLAGS: -O3 -march=native -std=c++23 -I../../../deps/eigen
// #cgo LDFLAGS: -lstdc++
// #include "cpu.h"
import "C"
import (
	"unsafe"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	pubenums "github.com/barbell-math/providentia/lib/pubEnums"
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

func ModelStates(
	clientID int64,
	historicalData []dal.ClientTrainingLogDataDateRangeAscendingRow,
	needsCalc []dal.ClientTrainingLogDataDateRangeAscendingRow,
	opts Opts,
) []dal.BulkCreateModelStatesParams {
	if len(needsCalc) == 0 {
		return []dal.BulkCreateModelStatesParams{}
	}

	modelStates := make([]dal.BulkCreateModelStatesParams, len(needsCalc))
	var needsCalcPntr *C.trainingLog_t
	var historicalDataPntr *C.trainingLog_t
	var modelStatesPntr *C.modelState_t
	if len(historicalData) > 0 {
		historicalDataPntr = (*C.trainingLog_t)(unsafe.Pointer(&historicalData[0]))
	}
	if len(needsCalc) > 0 {
		needsCalcPntr = (*C.trainingLog_t)(unsafe.Pointer(&needsCalc[0]))
	}
	if len(modelStates) > 0 {
		modelStatesPntr = (*C.modelState_t)(unsafe.Pointer(&modelStates[0]))
	}

	C.calcModelStates(
		C.long(clientID),
		C.int32_t(pubenums.ModelIDSimplifiedNegativeSpace),
		historicalDataPntr, C.long(len(historicalData)),
		needsCalcPntr, C.long(len(needsCalc)),
		modelStatesPntr, C.long(len(modelStates)),
		(*C.opts_t)(unsafe.Pointer(&opts)),
	)
	return modelStates
}
