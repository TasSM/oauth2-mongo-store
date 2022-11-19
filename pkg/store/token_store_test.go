package store_test

import (
	"errors"
	"testing"
	"time"

	"github.com/tassm/oauth2-mongo-store/pkg/store"
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

// create, read by different keys and then delete
func TestTokenStore_createAndRead(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	ts := store.NewMongoTokenStore(client, mongo_test_db)

	// save token
	saveToken(ts, t)

	token, err := ts.GetByCode(goodToken.Code)
	if err != nil {
		t.Errorf("failed to retrieve token by code with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)

	token, err = ts.GetByAccess(goodToken.Access)
	if err != nil {
		t.Errorf("failed to retrieve token by access with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)

	token, err = ts.GetByRefresh(goodToken.Refresh)
	if err != nil {
		t.Errorf("failed to retrieve token by refresh with error: %s", err.Error())
	}
	validateTokensAreEqual(token, goodToken, t)
	// delete the saved token
	err = ts.RemoveByCode(goodToken.Code)
	if err != nil {
		t.Errorf("failed to delete token by code with error: %s", err.Error())
	}
}

func TestTokenStore_createAndDelete(t *testing.T) {
	client, closefn := connectLocalMongo()
	defer closefn()

	// create a token store instance
	ts := store.NewMongoTokenStore(client, mongo_test_db)

	// save and delete by code
	saveToken(ts, t)
	err := ts.RemoveByCode(goodToken.Code)
	if err != nil {
		t.Errorf("failed to delete token by code with error: %s", err.Error())
	}
	_, err = ts.GetByCode(goodToken.Code)
	if !errors.Is(store.ErrorNoResult, err) {
		t.Error("retrieved token after it should have been deleted")
	}

	// save and delete by access
	saveToken(ts, t)
	err = ts.RemoveByAccess(goodToken.Access)
	if err != nil {
		t.Errorf("failed to delete token by access with error: %s", err.Error())
	}
	_, err = ts.GetByAccess(goodToken.Access)
	if !errors.Is(store.ErrorNoResult, err) {
		t.Error("retrieved token after it should have been deleted")
	}

	// save and delete by refresh
	saveToken(ts, t)
	err = ts.RemoveByRefresh(goodToken.Refresh)
	if err != nil {
		t.Errorf("failed to delete token by refresh with error: %s", err.Error())
	}
	_, err = ts.GetByRefresh(goodToken.Refresh)
	if !errors.Is(store.ErrorNoResult, err) {
		t.Error("retrieved token after it should have been deleted")
	}
}
