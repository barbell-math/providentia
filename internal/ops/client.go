package ops

import (
	"context"
	"net/mail"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

func CreateClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	clients ...types.Client,
) (opErr error) {
	for start, end := range batchIndexes(clients, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		for i := start; i < end; i++ {
			iterCd := clients[i]
			if iterCd.FirstName == "" {
				opErr = sberr.AppendError(
					types.InvalidClientErr, types.MissingFirstNameErr,
				)
				return
			}
			if iterCd.LastName == "" {
				opErr = sberr.AppendError(
					types.InvalidClientErr, types.MissingLastNameErr,
				)
				return
			}
			if iterCd.Email == "" {
				opErr = sberr.AppendError(
					types.InvalidClientErr, types.MissingEmailErr,
				)
				return
			}
			if _, err := mail.ParseAddress(iterCd.Email); err != nil {
				opErr = sberr.Wrap(
					sberr.AppendError(types.InvalidClientErr, err),
					"Invalid email: %s", iterCd.Email,
				)
				return
			}
		}

		var numRows int64
		chunk := clients[start:end]
		_ = dal.BulkCreateClientsParams(types.Client{})
		numRows, opErr = dal.Query1x2(
			dal.Q.BulkCreateClients, queries, ctxt,
			*(*[]dal.BulkCreateClientsParams)(unsafe.Pointer(&chunk)),
		)
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

func CreateClientsFromCSV(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	files ...string,
) (opErr error) {
	clients := []types.Client{}

	idxs := [3]int{}
	requestedCols := [3]string{"FirstName", "LastName", "Email"}

	for _, file := range files {
		if opErr = sbcsv.LoadCSVFile(file, &sbcsv.Opts{
			RequestedCols: requestedCols[:],
			Idxs:          idxs[:],
			Op: func(rowIdx int, row []string, idxs []int) error {
				clients = append(clients, types.Client{
					FirstName: row[idxs[0]],
					LastName:  row[idxs[1]],
					Email:     row[idxs[2]],
				})
				return nil
			},
		}); opErr != nil {
			return opErr
		}
	}

	return CreateClients(ctxt, state, queries, clients...)
}

func ReadNumClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	res, opErr = dal.Query0x2(dal.Q.GetNumClients, queries, ctxt)
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
	queries *dal.SyncQueries,
	emails ...string,
) (res []types.Client, opErr error) {
	res = make([]types.Client, len(emails))
	for start, end := range batchIndexes(emails, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var rawData []dal.GetClientsByEmailRow
		rawData, opErr = dal.Query1x2(
			dal.Q.GetClientsByEmail, queries, ctxt, emails[start:end],
		)
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
	queries *dal.SyncQueries,
	clients ...types.Client,
) (opErr error) {
	cntr := 0
	for _, c := range clients {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		_ = dal.UpdateClientByEmailParams(types.Client{})
		opErr = dal.Query1x1(
			dal.Q.UpdateClientByEmail, queries, ctxt,
			*(*dal.UpdateClientByEmailParams)(unsafe.Pointer(&c)),
		)
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

	state.Log.Log(ctxt, sblog.VLevel(3), "Updated clients", "Num", cntr)
	return
}

func DeleteClients(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	emails ...string,
) (opErr error) {
	// Deleting all referenced/referencing data is handled by cascade rules

	var count int64
	count, opErr = dal.Query1x2(dal.Q.DeleteClientsByEmail, queries, ctxt, emails)
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
