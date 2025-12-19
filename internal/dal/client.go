package dal

import (
	"context"
	"fmt"

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
	readClientsByEmailSql = `
SELECT
	providentia.client.first_name,
	providentia.client.last_name,
	providentia.client.email
FROM providentia.client 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(email, ord)
USING (email) 
ORDER BY ord;
`

	findClientsByEmailSql = `
SELECT
	providentia.client.first_name,
	providentia.client.last_name,
	providentia.client.email,
	ord::INT8
FROM providentia.client 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(email, ord)
USING (email) 
ORDER BY ord;
`

	ensureClientsExistSql = `
INSERT INTO providentia.client (first_name, last_name, email)
VALUES ($1, $2, $3)
ON CONFLICT (first_name, last_name, email) DO NOTHING;
`

	updateClientsSql = `
UPDATE providentia.client SET first_name=$1, last_name=$2
WHERE providentia.client.email=$3;
`

	deleteClientsSql = `
DELETE FROM providentia.client WHERE email = $1 RETURNING id;
`
)

func CreateClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clients []types.Client,
) error {
	cpy := CpyFromSlice[types.Client]{
		Data: clients,
		ValueGetter: func(v *types.Client, res *[]any) error {
			if len(*res) < 3 {
				*res = make([]any, 3)
			}
			(*res)[0] = v.FirstName
			(*res)[1] = v.LastName
			(*res)[2] = v.Email
			return nil
		},
	}
	if n, err := tx.CopyFrom(
		ctxt, pgx.Identifier{"providentia", "client"},
		[]string{"first_name", "last_name", "email"},
		&cpy,
	); err != nil {
		return sberr.AppendError(types.CouldNotCreateAllClientsErr, err)
	} else if n != int64(len(clients)) {
		return sberr.Wrap(
			types.CouldNotCreateAllClientsErr,
			"Expected to create %d clients but only created %d, rolling back",
			len(clients), n,
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Created new clients",
		"NumRows", len(clients),
	)
	return nil
}

func EnsureClientsExist(
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

		created := int64(0)
		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(
				ensureClientsExistSql,
				clients[i].FirstName, clients[i].LastName, clients[i].Email,
			)
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return err
			} else {
				created += cmdTag.RowsAffected()
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Ensured clients exist",
			"NumAffectedRows/NumRows", fmt.Sprintf("%d/%d", created, end-start),
		)
	}
	return nil
}

func ReadNumClients(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	num *int64,
) error {
	row := tx.QueryRow(ctxt, "SELECT COUNT(*) FROM providentia.client;")
	state.Log.Log(ctxt, sblog.VLevel(3), "DAL: Read total num clients")
	return row.Scan(num)
}

func ReadClientsByEmail(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadClientByEmailOpts,
) error {
	if len(*opts.Clients) < len(opts.Emails) {
		*opts.Clients = make([]types.Client, len(opts.Emails))
	} else if len(*opts.Clients) > len(opts.Emails) {
		*opts.Clients = (*opts.Clients)[:len(opts.Emails)]
	}
	for start, end := range batchIndexes(opts.Emails, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(ctxt, readClientsByEmailSql, opts.Emails[start:end])
		if err != nil {
			return err
		}

		cntr := start
		for rows.Next() {
			if (*opts.Clients)[cntr], err = pgx.RowToStructByName[types.Client](
				rows,
			); err != nil {
				rows.Close()
				return err
			}
			cntr++
		}
		rows.Close()

		if cntr != end {
			return sberr.Wrap(
				types.CouldNotReadAllClientsErr,
				"Only read %d clients out of batch of %d requests",
				cntr-start, end-start,
			)
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Read clients by email",
			"NumRows", end-start,
		)
	}
	return nil
}

func FindClientsByEmail(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindClientByEmailOpts,
) error {
	if len(*opts.Clients) < len(opts.Emails) {
		*opts.Clients = make([]types.Found[types.Client], len(opts.Emails))
	} else if len(*opts.Clients) > len(opts.Emails) {
		*opts.Clients = (*opts.Clients)[:len(opts.Emails)]
	}
	for start, end := range batchIndexes(opts.Emails, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(ctxt, findClientsByEmailSql, opts.Emails[start:end])
		if err != nil {
			return err
		}

		found := int64(0)
		for rows.Next() {
			iterVal := types.Client{}
			ord := int64(0)
			if err := rows.Scan(
				&iterVal.FirstName,
				&iterVal.LastName,
				&iterVal.Email,
				&ord,
			); err != nil {
				rows.Close()
				return err
			}
			(*opts.Clients)[int64(start)+ord-1] = types.Found[types.Client]{
				Value: iterVal,
				Found: true,
			}
			found++
		}
		rows.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Found clients by email",
			"NumFound/NumRows", fmt.Sprintf("%d/%d", found, end-start),
		)
	}
	return nil
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
	// Deleting all referenced/referencing data is handled by cascade rules
	for start, end := range batchIndexes(emails, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(deleteClientsSql, emails[i])
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return err
			} else if cmdTag.RowsAffected() == 0 {
				return sberr.Wrap(
					types.CouldNotDeleteAllClientsErr,
					"Could not delete client at idx %d (Does client exist?)",
					i,
				)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Deleted clients",
			"NumRows", end-start,
		)
	}
	return nil
}
