package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	ReadClientByEmailOpts struct {
		Emails  []string
		Clients *[]types.Client
	}

	FindClientByEmailOpts struct {
		Emails  []string
		Clients *[]types.Found[types.Client]
	}
)

const (
	clientTableName = "client"

	updateClientsSql = `
UPDATE providentia.client SET first_name=$1, last_name=$2
WHERE providentia.client.email=$3;
`
)

func CreateClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clients []types.Client,
) error {
	return genericCreate(
		ctxt, state, tx, &genericCreateOpts[types.Client]{
			TableName: clientTableName,
			Columns:   []string{"first_name", "last_name", "email"},
			Data:      clients,
			ValueGetter: func(v *types.Client, res *[]any) error {
				if len(*res) < 3 {
					*res = make([]any, 3)
				}
				(*res)[0] = v.FirstName
				(*res)[1] = v.LastName
				(*res)[2] = v.Email
				return nil
			},
			Err: types.CouldNotCreateAllClientsErr,
		},
	)
}

func EnsureClientsExist(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clients []types.Client,
) error {
	return genericEnsureExists(
		ctxt, state, tx, &genericCreateOpts[types.Client]{
			TableName: clientTableName,
			Columns:   []string{"first_name", "last_name", "email"},
			Data:      clients,
			ValueGetter: func(v *types.Client, res *[]any) error {
				*res = make([]any, 3)
				(*res)[0] = v.FirstName
				(*res)[1] = v.LastName
				(*res)[2] = v.Email
				return nil
			},
			Err: types.CouldNotCreateAllClientsErr,
		},
	)
}

func ReadNumClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	num *int64,
) error {
	return genericReadTotalNum(
		ctxt, state, tx, &genericReadTotalNumOpts{
			TableName: clientTableName,
			Res:       num,
		},
	)
}

func ReadClientsByEmail(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadClientByEmailOpts,
) error {
	return genericReadByUniqueId(
		ctxt, state, tx, &genericReadByUniqueIdOpts[string, types.Client]{
			TableName:  clientTableName,
			Columns:    []string{"first_name", "last_name", "email"},
			UniqueCol:  "email",
			IdsSqlType: "TEXT",
			Ids:        opts.Emails,
			Res:        opts.Clients,
			Err:        types.CouldNotReadAllClientsErr,
		},
	)
}

func FindClientsByEmail(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindClientByEmailOpts,
) error {
	return genericFindByUniqueId(
		ctxt, state, tx, &genericFindByUniqueIdOpts[
			string, types.Found[types.Client], types.Client,
		]{
			TableName:  clientTableName,
			Columns:    []string{"first_name", "last_name", "email"},
			UniqueCol:  "email",
			IdsSqlType: "TEXT",
			Ids:        opts.Emails,
			Res:        opts.Clients,
			SetScanValues: func(v *types.Client, res []any) {
				res[0] = &v.FirstName
				res[1] = &v.LastName
				res[2] = &v.Email
			},
			Err: types.CouldNotReadAllClientsErr,
		},
	)
}

func UpdateClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clients []types.Client,
) error {
	for start, end := range batchIndexes(clients, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(
				updateClientsSql,
				clients[i].FirstName, clients[i].LastName, clients[i].Email,
			)
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return err
			} else if cmdTag.RowsAffected() == 0 {
				return sberr.Wrap(
					types.CouldNotUpdateAllClientsErr,
					"Could not update client at idx %d (Does client exist?)",
					i,
				)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Updated clients",
			"NumRows", end-start,
		)
	}
	return nil
}

func DeleteClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	emails []string,
) error {
	return genericDeleteByUniqueId(
		ctxt, state, tx, &genericDeleteByUniqueIdOpts[string]{
			Ids:       emails,
			TableName: clientTableName,
			UniqueCol: "email",
			Err:       types.CouldNotDeleteAllClientsErr,
		},
	)
}
