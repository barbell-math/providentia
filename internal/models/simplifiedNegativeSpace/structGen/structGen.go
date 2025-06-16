package main

import (
	"reflect"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	simplifiednegativespace "github.com/barbell-math/providentia/internal/models/simplifiedNegativeSpace"
	structgen "github.com/barbell-math/smoothbrain-cgostructgen"
)

func main() {
	sg := structgen.New(structgen.Opts{
		ExitOnErr: true,
		StructRename: map[string]string{
			reflect.TypeFor[dal.ClientTrainingLogDataDateRangeAscendingRow]().Name(): "trainingLog",
			reflect.TypeFor[dal.BulkCreateModelStatesParams]().Name():                "modelState",
			"Opts": "opts",
		},
	})
	structgen.GenerateFor[simplifiednegativespace.Opts](sg)
	structgen.GenerateFor[dal.ClientTrainingLogDataDateRangeAscendingRow](sg)
	structgen.GenerateFor[dal.BulkCreateModelStatesParams](sg)
	sg.WriteTo("./cgoStructs.h", "CGO_STRUCTS")
}
