package provlib

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

// Gets up to `limit` training log entries for the supplied client. The training
// log entries will be sorted by date in descending order. If `limit` is <=0 no
// work will be done.
//
// The context must have a [State] variable.
//
// No changes will be made to the database.
func GetTrainingLogData(
	ctxt context.Context,
	clientEmail string,
	limit int32,
) (tlData []dal.GetClientTrainingLogDataRow, opErr error) {
	tlData = []dal.GetClientTrainingLogDataRow{}
	if limit <= 0 {
		return
	}

	state, queries, cleanup, err := dbOpInfo(ctxt)
	if err != nil {
		opErr = err
		return
	}
	defer cleanup(&opErr)

	tlData, opErr = queries.GetClientTrainingLogData(
		ctxt,
		dal.GetClientTrainingLogDataParams{
			Email: clientEmail,
			Limit: limit,
		},
	)
	state.Log.Log(
		ctxt, sblog.VLevel(2),
		"Retreived client training log data",
		"NumRows", len(tlData),
	)

	return
}
