package provlib

import (
	"context"
	"errors"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	sberr "github.com/barbell-math/smoothbrain-errs"
	sblog "github.com/barbell-math/smoothbrain-logging"
)

type (
	ClientData struct {
		FirstName string
		LastName  string
		Email     string
	}
)

var (
	InvalidClientErr               = errors.New("Invalid client")
	CouldNotFindRequestedClientErr = errors.New(
		"Could not find the requested client",
	)
	CouldNotAddClientsErr = errors.New("Could not add the requested clients")
)

// Adds the supplied clients to the database. The supplied first name, last
// name, and email for each client must not be an empty string. Emails must not
// be duplicated, including the set of client emails that are already in the
// database.
//
// The context must have a [State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func AddClients(
	ctxt context.Context,
	cd ...ClientData,
) (opErr error) {
	if len(cd) == 0 {
		return
	}

	state, queries, cleanup, err := dbOpInfo(ctxt)
	if err != nil {
		opErr = err
		return
	}
	defer cleanup(&opErr)

	for start, end := range batchIndexes(cd, int(state.Conf.BatchSize)) {
		clients := make([]dal.BulkCreateClientsParams, len(cd))
		for i := start; i < end; i++ {
			iterCd := cd[i]
			if iterCd.FirstName == "" {
				opErr = sberr.Wrap(InvalidClientErr, "First name must not be empty")
				return
			}
			if iterCd.LastName == "" {
				opErr = sberr.Wrap(InvalidClientErr, "Last name must not be empty")
				return
			}
			if iterCd.Email == "" {
				opErr = sberr.Wrap(InvalidClientErr, "Email must not be empty")
				return
			}
			clients[i] = dal.BulkCreateClientsParams(iterCd)
		}

		var numRows int64
		numRows, opErr = queries.BulkCreateClients(ctxt, clients)
		if opErr != nil {
			opErr = sberr.AppendError(CouldNotAddClientsErr, opErr)
			return
		}
		state.Log.Log(
			ctxt, sblog.VLevel(2),
			"Added new clients",
			"NumRows", numRows,
		)
	}

	return
}
