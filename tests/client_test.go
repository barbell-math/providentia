package tests

import (
	"context"
	"fmt"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestClient(t *testing.T) {
	t.Run("failingNoWrites", clientFailingNoWrites)
	t.Run("duplicateEmail", clientDuplicateEmail)
	t.Run("createRead", clientCreateRead)
	t.Run("ensureRead", clientEnsureRead)
	t.Run("createFind", clientCreateFind)
	t.Run("createUpdateRead", clientCreateUpdateRead)
	t.Run("createDeleteRead", clientCreateDeleteRead)
	t.Run("createCSVRead", clientCreateCSVRead)
	t.Run("ensureCSVRead", clientEnsureCSVRead)
}

func clientFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("missingFirstName", clientMissingFirstName(ctxt))
	t.Run("missingLastName", clientMissingLastName(ctxt))
	t.Run("missingEmail", clientMissingEmail(ctxt))
	t.Run("invalidEmail", clientInvalidEmail(ctxt))

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, n)
}

func clientMissingFirstName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateClients(ctxt, types.Client{
			LastName: "LName",
			Email:    "email@email.com",
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllClientsErr, err,
			`new row for relation \"client\" violates check constraint \"first_name_not_empty\" \(SQLSTATE 23514\)`,
		)
	}
}

func clientMissingLastName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateClients(ctxt, types.Client{
			FirstName: "FName",
			Email:     "email@email.com",
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllClientsErr, err,
			`new row for relation \"client\" violates check constraint \"last_name_not_empty\" \(SQLSTATE 23514\)`,
		)
	}
}

func clientMissingEmail(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateClients(ctxt, types.Client{
			FirstName: "FName",
			LastName:  "LName",
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllClientsErr, err,
			`new row for relation \"client\" violates check constraint \"email_not_empty\" \(SQLSTATE 23514\)`,
		)
	}
}

func clientInvalidEmail(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateClients(ctxt, types.Client{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "asdfasdf",
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllClientsErr, err,
			`new row for relation \"client\" violates check constraint \"valid_email_format\" \(SQLSTATE 23514\)`,
		)
	}
}

func clientDuplicateEmail(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	}, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.ContainsError(
		t, types.CouldNotCreateAllClientsErr, err,
		`duplicate key value violates unique constraint "client_email_key" \(SQLSTATE 23505\)`,
	)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, n)
}

func clientCreateRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	clients := []types.Client{
		{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email1@email.com",
		},
	}
	err := logic.CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	readClients, err := logic.ReadClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients, readClients)

	readClients, err = logic.ReadClientsByEmail(ctxt, "asdfasdf")
	sbtest.ContainsError(
		t, types.CouldNotReadAllClientsErr, err,
		"Only read 0 clients out of batch of 1 requests",
	)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientEnsureRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	clients := []types.Client{
		{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email1@email.com",
		},
	}
	err := logic.EnsureClientsExist(ctxt, clients...)
	sbtest.Nil(t, err)

	readClients, err := logic.ReadClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients, readClients)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)

	err = logic.EnsureClientsExist(ctxt, clients...)
	sbtest.Nil(t, err)

	readClients, err = logic.ReadClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients, readClients)

	n, err = logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientCreateFind(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	clients := []types.Client{
		{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email1@email.com",
		},
	}
	err := logic.EnsureClientsExist(ctxt, clients...)
	sbtest.Nil(t, err)

	foundClients, err := logic.FindClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email, "asdfasdf",
	)
	fmt.Println(foundClients)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, []types.Found[types.Client]{{
		Found: true,
		Value: clients[0],
	}, {
		Found: true,
		Value: clients[1],
	}, {
		Found: false,
	}}, foundClients)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientCreateUpdateRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	clients := []types.Client{
		{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email1@email.com",
		},
	}
	err := logic.CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	clients[0].FirstName = "fname"
	clients[0].LastName = "lname"
	clients[1].FirstName = "fname"
	clients[1].LastName = "lname"
	err = logic.UpdateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	readClients, err := logic.ReadClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients, readClients)

	err = logic.UpdateClients(ctxt, types.Client{
		Email:     "email2@email.com",
		FirstName: "asdf",
		LastName:  "asdf",
	})
	sbtest.ContainsError(
		t, types.CouldNotUpdateAllClientsErr, err,
		`Could not update client at idx 0 \(Does client exist\?\)`,
	)

	readClients, err = logic.ReadClientsByEmail(
		ctxt, clients[0].Email, clients[1].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients, readClients)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientCreateDeleteRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	clients := []types.Client{
		{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email1@email.com",
		}, {
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email2@email.com",
		},
	}
	err := logic.CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, n)

	err = logic.DeleteClients(ctxt, clients[0].Email)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)

	readClients, err := logic.ReadClientsByEmail(
		ctxt, clients[1].Email, clients[2].Email,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, clients[1:], readClients)

	err = logic.DeleteClients(ctxt, clients[0].Email)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllClientsErr, err,
		`Could not delete client at idx 0 \(Does client exist\?\)`,
	)

	n, err = logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientCreateCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClientsFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/clientData/clients.csv",
	)
	sbtest.Nil(t, err)

	numClients, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, numClients)

	client, err := logic.ReadClientsByEmail(ctxt, "one@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, client[0], types.Client{
		FirstName: "OneFN",
		LastName:  "OneLN",
		Email:     "one@gmail.com",
	})

	client, err = logic.ReadClientsByEmail(ctxt, "two@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, client[0], types.Client{
		FirstName: "TwoFN",
		LastName:  "TwoLN",
		Email:     "two@gmail.com",
	})

	_, err = logic.ReadClientsByEmail(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotReadAllClientsErr, err)

	err = logic.CreateClientsFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/clientData/clients.csv",
	)
	sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllClientsErr, err,
		`duplicate key value violates unique constraint "client_email_key" \(SQLSTATE 23505\)`,
	)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}

func clientEnsureCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.EnsureClientsExistFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/clientData/clients.csv",
	)
	sbtest.Nil(t, err)

	numClients, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, numClients)

	client, err := logic.ReadClientsByEmail(ctxt, "one@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, client[0], types.Client{
		FirstName: "OneFN",
		LastName:  "OneLN",
		Email:     "one@gmail.com",
	})

	client, err = logic.ReadClientsByEmail(ctxt, "two@gmail.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, client[0], types.Client{
		FirstName: "TwoFN",
		LastName:  "TwoLN",
		Email:     "two@gmail.com",
	})

	_, err = logic.ReadClientsByEmail(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotReadAllClientsErr, err)

	err = logic.EnsureClientsExistFromCSV(
		ctxt, &sbcsv.Opts{}, "./testData/clientData/clients.csv",
	)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, n)
}
