package logic

import (
	"context"
	"fmt"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestClient(t *testing.T) {
	t.Run("failingNoWrites", clientFailingNoWrites)
	t.Run("duplicateEmail", clientDuplicateEmail)
	t.Run("transactionRollback", clientTransactionRollback)
	t.Run("addGet", clientAddGet)
	t.Run("addUpdateGet", clientAddUpdateGet)
	t.Run("addDeleteGet", clientAddDeleteGet)
}

func clientFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)
	t.Run("missingFirstName", clientMissingFirstName(ctxt))
	t.Run("missingLastName", clientMissingLastName(ctxt))
	t.Run("missingEmail", clientMissingEmail(ctxt))
	t.Run("invalidEmail", clientInvalidEmail(ctxt))

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, numClients)
}

func clientMissingFirstName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateClients(ctxt, types.Client{
			LastName: "LName",
			Email:    "email@email.com",
		})
		sbtest.ContainsError(t, types.InvalidClientErr, err)
	}
}

func clientMissingLastName(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateClients(ctxt, types.Client{
			FirstName: "FName",
			Email:     "email@email.com",
		})
		sbtest.ContainsError(t, types.InvalidClientErr, err)
	}
}

func clientMissingEmail(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateClients(ctxt, types.Client{
			FirstName: "FName",
			LastName:  "LName",
		})
		sbtest.ContainsError(t, types.InvalidClientErr, err)
	}
}

func clientInvalidEmail(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateClients(ctxt, types.Client{
			FirstName: "FName",
			LastName:  "LName",
			Email:     "email",
		})
		sbtest.ContainsError(t, types.InvalidClientErr, err)
	}
}

func clientDuplicateEmail(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	clients := make([]types.Client, 13)
	for i := range len(clients) {
		clients[i] = types.Client{
			FirstName: fmt.Sprintf("FName%d", i),
			LastName:  fmt.Sprintf("LName%d", i),
			Email:     fmt.Sprintf("email%d@email.com", i),
		}
	}
	clients[len(clients)-1].Email = fmt.Sprintf(
		"email%d@email.com", len(clients)-2,
	)

	err := CreateClients(ctxt, clients...)
	sbtest.ContainsError(t, types.CouldNotAddClientsErr, err)

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, numClients)
}

func clientTransactionRollback(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	clients := make([]types.Client, 13)
	for i := range len(clients) {
		clients[i] = types.Client{
			FirstName: fmt.Sprintf("FName%d", i),
			LastName:  fmt.Sprintf("LName%d", i),
			Email:     fmt.Sprintf("email%d@email.com", i),
		}
	}

	err := CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)
	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)

	for i := 0; i < 5; i++ {
		clients[i].Email = fmt.Sprintf("email%d@email.com", i+len(clients))
	}

	err = CreateClients(ctxt, clients...)
	sbtest.ContainsError(t, types.CouldNotAddClientsErr, err)
	numClients, err = ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)
}

func clientAddGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	clients := make([]types.Client, 13)
	for i := range len(clients) {
		clients[i] = types.Client{
			FirstName: fmt.Sprintf("FName%d", i),
			LastName:  fmt.Sprintf("LName%d", i),
			Email:     fmt.Sprintf("email%d@email.com", i),
		}
	}

	err := CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)

	for i := range len(clients) {
		res, err := ReadClientsByEmail(ctxt, clients[i].Email)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, clients[i].FirstName, res[0].FirstName)
		sbtest.Eq(t, clients[i].LastName, res[0].LastName)
		sbtest.Eq(t, clients[i].Email, res[0].Email)
	}

	_, err = ReadClientsByEmail(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
}

func clientAddUpdateGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	clients := make([]types.Client, 13)
	for i := range len(clients) {
		clients[i] = types.Client{
			FirstName: fmt.Sprintf("FName%d", i),
			LastName:  fmt.Sprintf("LName%d", i),
			Email:     fmt.Sprintf("email%d@email.com", i),
		}
	}

	err := CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)

	for i := range len(clients) {
		clients[i].FirstName = fmt.Sprintf("FName%d", i+1)
		clients[i].LastName = fmt.Sprintf("LName%d", i+1)
	}
	err = UpdateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	err = UpdateClients(ctxt, types.Client{Email: "bad@email.com"})
	sbtest.Nil(t, err)
	numClients, err = ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)

	for i := range len(clients) {
		res, err := ReadClientsByEmail(ctxt, clients[i].Email)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, clients[i].FirstName, res[0].FirstName)
		sbtest.Eq(t, clients[i].LastName, res[0].LastName)
		sbtest.Eq(t, clients[i].Email, res[0].Email)
	}

	_, err = ReadClientsByEmail(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
}

func clientAddDeleteGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	clients := make([]types.Client, 13)
	for i := range len(clients) {
		clients[i] = types.Client{
			FirstName: fmt.Sprintf("FName%d", i),
			LastName:  fmt.Sprintf("LName%d", i),
			Email:     fmt.Sprintf("email%d@email.com", i),
		}
	}

	err := CreateClients(ctxt, clients...)
	sbtest.Nil(t, err)

	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 13, numClients)

	emails := [5]string{}
	for i := range len(emails) {
		emails[i] = clients[i].Email
	}
	err = DeleteClients(ctxt, emails[:]...)
	sbtest.Nil(t, err)
	numClients, err = ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(clients)-len(emails)), numClients)

	err = DeleteClients(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotDeleteRequestedClientErr, err)
	sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
	numClients, err = ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(clients)-len(emails)), numClients)

	for i := range numClients {
		offset := int(i) + len(emails)
		res, err := ReadClientsByEmail(ctxt, clients[offset].Email)
		sbtest.Nil(t, err)
		sbtest.Eq(t, 1, len(res))
		sbtest.Eq(t, clients[offset].FirstName, res[0].FirstName)
		sbtest.Eq(t, clients[offset].LastName, res[0].LastName)
		sbtest.Eq(t, clients[offset].Email, res[0].Email)
	}

	_, err = ReadClientsByEmail(ctxt, "bad@email.com")
	sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
}
