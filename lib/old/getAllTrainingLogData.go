package provlib

import (
	"context"
	"unsafe"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	sblog "github.com/barbell-math/smoothbrain-logging"
)

type (
	AllClientsTrainingLogData dal.GetAllClientsTrainingLogDataRow
)

// Gets up to `limit` training log entries across all clients. The training
// log entries will be sorted by client, with each clients data sorted by date
// in descending order. If `limit` is <=0 no work will be done.
//
// The context must have a [State] variable.
//
// No changes will be made to the database.
func GetAllClientsTrainingLogData(
	ctxt context.Context,
	limit int32,
) (tlData []AllClientsTrainingLogData, opErr error) {
	tlData = []AllClientsTrainingLogData{}
	if limit <= 0 {
		return
	}

	state, queries, cleanup, err := dbOpInfo(ctxt)
	if err != nil {
		opErr = err
		return
	}
	defer cleanup(&opErr)

	var rawData []dal.GetAllClientsTrainingLogDataRow
	rawData, opErr = queries.GetAllClientsTrainingLogData(ctxt, limit)
	state.Log.Log(
		ctxt, sblog.VLevel(2),
		"Retrieved training log data",
		"NumRows", len(tlData),
	)
	tlData = *(*[]AllClientsTrainingLogData)(unsafe.Pointer(&rawData))

	return
}
