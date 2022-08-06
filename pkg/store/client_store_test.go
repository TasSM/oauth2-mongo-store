package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/TasSM/oauth2-mongo-store/pkg/store"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

var goodClient *models.Client = &models.Client{
	ID:     "testId",
	Secret: "testSecret",
	Domain: "localhost",
	UserID: "test",
}

// compare 2 clients
func validateClientsAreEqual(new oauth2.ClientInfo, old *models.Client, t *testing.T) {
	// don't compare the ID as it is DB internal
	if new.GetSecret() != old.Secret {
		t.Errorf("Secret mismatch")
	}
	if new.GetDomain() != old.Domain {
		t.Errorf("Domain mismatch")
	}
	if new.GetUserID() != old.UserID {
		t.Errorf("UserId mismatch")
	}
}

func TestClientStore_SetReadDelete(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// set client
	cs := store.NewMongoClientStore(client, mongo_test_db)
	err := cs.Set(goodClient)
	if err != nil {
		t.Errorf("create client failed with error: %s", err.Error())
	}

	// read client
	retrievedClient, err := cs.GetByID(context.TODO(), goodClient.UserID)
	if err != nil {
		t.Errorf("failed to retrieve client by ClientID with error: %s", err.Error())
	}
	if retrievedClient == nil {
		t.Errorf("retrieved client is nil!")
	}
	validateClientsAreEqual(retrievedClient, goodClient, t)

	// delete client
	err = cs.RemoveByID(goodClient.UserID)
	if err != nil {
		t.Errorf("failed to delete client with error: %s", err.Error())
	}
	_, err = cs.GetByID(context.TODO(), goodClient.UserID)
	if !errors.Is(store.ErrorNoResult, err) {
		t.Error("should not be able to retrieve client`")
	}
}
