package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
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
func CreateClients(ctxt context.Context, clients ...types.Client) (opErr error) {
	if len(clients) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			return ops.CreateClients(ctxt, state, queries, clients...)
		},
	})
}

// Checks that the supplied clients are present in the database and adds them if
// they are not present. In order for the supplied clients to be be considered
// already present the first name, last name, and email fields must all match.
// Any newly created clients must satisfy the uniqueness constraints outlined by
// [CreateClients].
//
// This function will be slower than [CreateClients], so if you are working with
// large amounts of data and are ok with erroring on duplicated clients consider
// using [CreateClients].
//
// The context must have a [types.State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func EnsureClientsExist(
	ctxt context.Context,
	clients ...types.Client,
) (opErr error) {
	if len(clients) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.EnsureClientsExist(ctxt, state, queries, clients...)
		},
	})
}

// Adds the clients supplied in the csv files to the database. Has the same
// behavior as [CreateClients] other than getting the clients from csv files.
// The csv files are expected to have column names on the first row and the
// following columns must be present as identified by the column name on the
// first row. More columns may be present, they will be ignored.
//   - FirstName (string): the first name of the client
//   - LastName (string): the last name of the client
//   - Email (string): the email of the client
//
// The `ReuseRecord` field on opts will be set to true before loading the csv
// file. All other options are left alone.
//
// The context must have a [types.State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateClientsFromCSV(
	ctxt context.Context,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.CreateClientsFromCSV(
				ctxt, state, queries, opts, files...,
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
// they do not exist an error will be returned. The order of the returned
// clients will match the order of the supplied client emails.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadClientsByEmail(
	ctxt context.Context,
	emails ...string,
) (res []types.Client, opErr error) {
	if len(emails) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadClientsByEmail(ctxt, state, queries, emails...)
			return
		},
	})
	return
}

// Gets the client data associated with the supplied emails if they exist. If a
// client exists it will be put in the returned slice and the found flag will be
// set to true. If a client does not exist the value in the slice will be a zero
// initialized client and the found flag will be set to false. No error will be
// returned if a client does not exist. The order of the returned clients will
// match the order of the supplied client emails.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func FindClientsByEmail(
	ctxt context.Context,
	emails ...string,
) (res []types.Found[types.Client], opErr error) {
	if len(emails) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.FindClientsByEmail(ctxt, state, queries, emails...)
			return
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
func UpdateClients(ctxt context.Context, clients ...types.Client) (opErr error) {
	if len(clients) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.UpdateClients(ctxt, state, queries, clients...)
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
