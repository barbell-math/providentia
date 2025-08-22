package ops

import (
	"context"
	"net/mail"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

func CreateClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	clients ...dal.BulkCreateClientsParams,
) (opErr error) {
	for start, end := range batchIndexes(clients, int(state.Global.BatchSize)) {
		for i := start; i < end; i++ {
			iterCd := clients[i]
			if iterCd.FirstName == "" {
				opErr = sberr.Wrap(
					types.InvalidClientErr, "First name must not be empty",
				)
				return
			}
			if iterCd.LastName == "" {
				opErr = sberr.Wrap(
					types.InvalidClientErr, "Last name must not be empty",
				)
				return
			}
			if iterCd.Email == "" {
				opErr = sberr.Wrap(
					types.InvalidClientErr, "Email must not be empty",
				)
				return
			}
			if _, err := mail.ParseAddress(iterCd.Email); err != nil {
				opErr = sberr.Wrap(
					types.InvalidClientErr, "Invalid email: %s", err,
				)
				return
			}
		}

		var numRows int64
		// The buffered writer is not used because it would create a copy of the
		// clients, which is unnecessary in this case
		numRows, opErr = queries.BulkCreateClients(ctxt, clients[start:end])
		if opErr != nil {
			opErr = sberr.AppendError(types.CouldNotAddClientsErr, opErr)
			return
		}
		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Added new clients",
			"NumRows", numRows,
		)
	}

	return
}

func ReadNumClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
) (res int64, opErr error) {
	res, opErr = queries.GetNumClients(ctxt)
	if opErr != nil {
		opErr = sberr.AppendError(types.CouldNotGetNumClientsErr, opErr)
		return
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read num clients",
	)
	return
}

func ReadClientsByEmail(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	emails ...string,
) (res []types.Client, opErr error) {
	res = make([]types.Client, len(emails))

	// Note: the client cache is not updated because that would require
	// returning an types.IdWrapper rather than just a types.Client struct and
	// that would require copying all of the returned results one at a time
	// rather than in chunks with a copy command.
	for start, end := range batchIndexes(emails, int(state.Global.BatchSize)) {
		var rawData []dal.GetClientsByEmailRow
		rawData, opErr = queries.GetClientsByEmail(ctxt, emails[start:end])
		if opErr != nil {
			opErr = sberr.AppendError(types.CouldNotFindRequestedClientErr, opErr)
			return
		}
		if len(rawData) != end-start {
			opErr = types.CouldNotFindRequestedClientErr
			return
		}

		_ = types.Client(dal.GetClientsByEmailRow{})
		copy(
			res[start:end],
			*(*[]types.Client)(unsafe.Pointer(&rawData)),
		)

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Read clients from email",
			"Num", len(rawData),
		)
	}

	return
}

func UpdateClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	clients ...dal.UpdateClientByEmailParams,
) (opErr error) {
	cntr := 0
	for _, c := range clients {
		// Note: the client cache does not need to be updated because the email
		// (and hence id in the database) does not change.
		opErr = queries.UpdateClientByEmail(ctxt, c)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotUpdateRequestedClientErr, opErr,
			)
			return
		}
		cntr++
	}
	if cntr != len(clients) {
		opErr = sberr.AppendError(
			types.CouldNotUpdateRequestedClientErr,
			types.CouldNotFindRequestedClientErr,
		)
		return
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Updated clients",
		"Num", cntr,
	)

	return
}

func DeleteClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.Queries,
	emails ...string,
) (opErr error) {
	// TODO - delete all referenced training log data, video data, model data

	for _, e := range emails {
		state.ClientCache.Invalidate(e)
	}

	var count int64
	count, opErr = queries.DeleteClientsByEmail(ctxt, emails)
	if opErr != nil {
		opErr = sberr.AppendError(types.CouldNotDeleteRequestedClientErr, opErr)
		return
	}
	if count != int64(len(emails)) {
		opErr = sberr.AppendError(
			types.CouldNotDeleteRequestedClientErr,
			types.CouldNotFindRequestedClientErr,
		)
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Deleted clients",
		"Num", count,
	)

	return
}
