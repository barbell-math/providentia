package logic

import (
	"context"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Adds the supplied clients to the database. The supplied first name, last
// name, and email for each client must not be an empty string. Emails must not
// be duplicated, including the set of client emails that are already in the
// database.
//
// The context must have a [types.State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateClients(
	ctxt context.Context,
	clients ...types.Client,
) (opErr error) {
	if len(clients) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			_ = dal.BulkCreateClientsParams(types.Client{})
			return ops.CreateClients(
				ctxt, state, queries,
				*(*[]dal.BulkCreateClientsParams)(unsafe.Pointer(&clients))...,
			)
		},
	})
}

// Gets the total number of clients in the database.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumClients(ctxt context.Context) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadNumClients(ctxt, state, queries)
			return err
		},
	})
	return
}

// Gets the client data associated with the supplied emails if they exist. If
// they do not exist an error will be returned.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadClientsByEmail(
	ctxt context.Context,
	emails ...string,
) (res []types.Client, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadClientsByEmail(ctxt, state, queries, emails...)
			return err
		},
	})
	return
}

// Updates the supplied clients, as identified by their email, with the data
// from the supplied structs. Emails cannot be updated due to their uniqueness
// constraint. If a client is supplied with an email that does not exist in the
// database an error will be returned.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func UpdateClients(
	ctxt context.Context,
	clients ...types.Client,
) (opErr error) {
	if len(clients) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			_ = dal.UpdateClientByEmailParams(types.Client{})
			return ops.UpdateClients(
				ctxt, state, queries,
				*(*[]dal.UpdateClientByEmailParams)(unsafe.Pointer(&clients))...,
			)
		},
	})
}

// Deletes the supplied clients, as identified by their email. All data
// associated with the client will be deleted.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func DeleteClients(ctxt context.Context, emails ...string) (opErr error) {
	if len(emails) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.DeleteClients(ctxt, state, queries, emails...)
		},
	})
}
