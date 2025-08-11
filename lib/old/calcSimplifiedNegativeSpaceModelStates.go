package provlib

import (
	"context"
	"errors"
	"slices"
	"time"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	simplifiednegativespace "code.barbellmath.net/barbell-math/providentia/internal/models/simplifiedNegativeSpace"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	StartDateMustBeBeforeEndDateErr = errors.New(
		"The supplied start date must be before the supplied end date",
	)
	CouldNotAddModelStatesErr = errors.New(
		"Could not add the requested model states",
	)
)

// Calculates the model states for the simplified negative space model over the
// supplied date range. The calculated model states will be uploaded to the
// database.
//
// If the supplied date range is outside of the date range available in the
// specified clients training log data then no error will be returned and no
// work will be done.
//
// `startDate` must be < `endDate` or a [StartDateMustBeBeforeEndDateErr] will
// be returned.
//
// If any error occurs no changes will be made to the database.
func CalcSimplifiedNegativeSpaceModelStates(
	ctxt context.Context,
	clientEmail string,
	startDate time.Time,
	endDate time.Time,
) (opErr error) {
	state, queries, cleanup, err := dbOpInfo(ctxt)
	if err != nil {
		opErr = err
		return
	}
	defer cleanup(&opErr)

	if endDate.Before(startDate) {
		opErr = sberr.Wrap(
			StartDateMustBeBeforeEndDateErr,
			"Start date: %s End date: %s", startDate, endDate,
		)
	}

	var clientID int64
	clientID, opErr = queries.GetClientIDFromEmail(ctxt, clientEmail)
	if opErr != nil {
		opErr = sberr.AppendError(CouldNotFindRequestedClientErr, opErr)
		return
	}

	allData, err := queries.ClientTrainingLogDataDateRangeAscending(
		ctxt, dal.ClientTrainingLogDataDateRangeAscendingParams{
			ClientID:      clientID,
			DatePerformed: pgtype.Date{Time: endDate, Valid: true},
		},
	)
	daysSinceStartDate := int32(endDate.Sub(startDate).Hours() / 24)
	if allData[len(allData)-1].DaysSince > daysSinceStartDate {
		// No work to be done if the startDate is beyond the available data
		return
	}
	splitIdx, _ := slices.BinarySearchFunc(
		allData, daysSinceStartDate,
		func(e dal.ClientTrainingLogDataDateRangeAscendingRow, t int32) int {
			// "a negative num if the slice element precedes the target"
			// Returns <0 when e.DaysSince is > daysSinceStartDate
			// "a positive num if the slice element follows the target"
			// Returns >0 when e.DaysSince is < daysSinceStartDate
			return int(t - e.DaysSince)
		},
	)

	// Where the real work is done...
	modelStates := simplifiednegativespace.ModelStates(
		clientID, allData, int64(splitIdx),
		state.Conf.SimplifiedNegativeSpaceModel,
	)

	var numRows int64
	numRows, opErr = queries.BulkCreateModelStates(ctxt, modelStates)
	if opErr != nil {
		opErr = sberr.AppendError(CouldNotAddModelStatesErr, opErr)
		return
	}
	state.Log.Log(
		ctxt, sblog.VLevel(2),
		"Inserted simplified negative space model states",
		"ClientID", clientID,
		"NumRows", numRows,
	)

	return
}
