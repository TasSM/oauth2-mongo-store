package store_test

import (
	"testing"
	"time"

	"github.com/TasSM/oauth2-mongo-store/pkg/store"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

var goodToken *models.Token = &models.Token{
	ClientID:      "client1",
	UserID:        "user1a",
	RedirectURI:   "http://localhost/",
	Scope:         "test",
	Code:          "123",
	CodeCreateAt:  time.Now(),
	CodeExpiresIn: time.Second * 10,
}

// compare 2 tokens
func validateTokensAreEqual(new oauth2.TokenInfo, old *models.Token, t *testing.T) {
	if new.GetClientID() != old.ClientID {
		t.Errorf("ClientId mismatch")
	}
	if new.GetUserID() != old.UserID {
		t.Errorf("UserId mismatch")
	}
	if new.GetRedirectURI() != old.RedirectURI {
		t.Errorf("RedirectUri mismatch")
	}
	if new.GetScope() != old.Scope {
		t.Errorf("Scope mismatch")
	}
	if new.GetCode() != old.Code {
		t.Errorf("Code mismatch")
	}
	if !new.GetCodeCreateAt().Equal(old.CodeCreateAt) {
		t.Errorf("CodeCreateAt mismatch")
	}
	if time.Duration(new.GetCodeExpiresIn().Seconds()) != time.Duration(old.CodeExpiresIn.Seconds()) {
		t.Errorf("CodeExpiresIn mismatch")
	}
}

func saveToken(ts *store.MongoTokenStore, t *testing.T) {
	err := ts.Create(goodToken)
	if err != nil {
		t.Errorf("create token failed with error: %s", err.Error())
	}
}

// CREATE TESTS //

func TestTokenStore_shouldCreateToken(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)
}

// READ TESTS //
func TestTokenStore_shouldGetByCode(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	token, err := tokenStore.GetByCode(goodToken.Code)
	if err != nil {
		t.Errorf("failed to retrieve token by code with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)
}

func TestTokenStore_shouldGetByAccess(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	token, err := tokenStore.GetByAccess(goodToken.Access)
	if err != nil {
		t.Errorf("failed to retrieve token by access with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)
}

func TestTokenStore_shouldGetByRefresh(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	token, err := tokenStore.GetByRefresh(goodToken.Refresh)
	if err != nil {
		t.Errorf("failed to retrieve token by refresh with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)
}

// DELETE TESTS //
func TestTokenStore_shouldDeleteByCode(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	err := tokenStore.RemoveByCode(goodToken.Code)
	if err != nil {
		t.Errorf("failed to delete token by code with error: %s", err.Error())
	}
	// check that it can't be retrieved
	_, err = tokenStore.GetByCode(goodToken.Code)
	if err != store.ErrorNoResult {
		t.Error("token was found when it should have been deleted")
	}
}

func TestTokenStore_shouldDeleteByAccess(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	err := tokenStore.RemoveByAccess(goodToken.Access)
	if err != nil {
		t.Errorf("failed to delete token by access with error: %s", err.Error())
	}
	// check that it can't be retrieved
	_, err = tokenStore.GetByAccess(goodToken.Access)
	if err != store.ErrorNoResult {
		t.Error("token was found when it should have been deleted")
	}
}

func TestTokenStore_shouldDeleteByRefresh(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	tokenStore := store.NewMongoTokenStore(client, mongo_test_db)
	saveToken(tokenStore, t)

	err := tokenStore.RemoveByRefresh(goodToken.Refresh)
	if err != nil {
		t.Errorf("failed to delete token by refresh with error: %s", err.Error())
	}
	// check that it can't be retrieved
	_, err = tokenStore.GetByRefresh(goodToken.Refresh)
	if err != store.ErrorNoResult {
		t.Error("token was found when it should have been deleted")
	}
}
