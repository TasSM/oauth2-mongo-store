package main

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// custom mongodb oauth client store for go-oauth2
// non exported client model as it is only required for internal implementation
// implements oauth2.TokenStore

type MongoTokenStore struct {
	dbclient   *mongo.Client
	database   string
	collection string
}

type tokenData struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	AccessID         string             `bson:"access_id"`
	CodeID           string             `bson:"code_id"`
	RefreshID        string             `bson:"refresh_id"`
	Data             []byte             `bson:"data"`
	RefreshExpiredAt time.Time          `bson:"expired_at_refresh"`
	ExpiredAt        time.Time          `bson:"expired_at"`
}

func NewMongoTokenStore(dbclient *mongo.Client, dbname string) *MongoTokenStore {
	return &MongoTokenStore{dbclient: dbclient, database: dbname, collection: "oauth_token"}
}

func (ts *MongoTokenStore) Create(info oauth2.TokenInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := ts.dbclient.Database(ts.database).Collection(ts.collection)
	token := tokenData{}
	// create raw token data
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if code := info.GetCode(); code != "" {
		token.CodeID = code
	}
	aexp := info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
	rexp := aexp
	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		if aexp.Second() > rexp.Second() {
			aexp = rexp
		}
		token.RefreshID = refresh
	}
	token.ExpiredAt = aexp
	token.RefreshExpiredAt = rexp
	token.AccessID = info.GetAccess()
	token.Data = jv
	indexes := []mongo.IndexModel{}
	indexes = append(indexes, mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "access_id", Value: bsonx.Int32(1)}},
	})
	indexes = append(indexes, mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "expired_at", Value: bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(math.Round(info.GetAccessExpiresIn().Seconds()))),
	})
	_, err = coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(ctx, token)
	return err
}

func (ts *MongoTokenStore) deleteTokenFor(key, value string) error {
	coll := ts.dbclient.Database(ts.database).Collection(ts.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.DeleteOne(ctx, bson.D{{key, value}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return err
}

func (ts *MongoTokenStore) RemoveByCode(code string) error {
	return ts.deleteTokenFor("code_id", code)
}

func (ts *MongoTokenStore) RemoveByRefresh(refresh string) error {
	return ts.deleteTokenFor("refresh_id", refresh)
}

func (ts *MongoTokenStore) RemoveByAccess(access string) error {
	return ts.deleteTokenFor("access_id", access)
}

func (ts *MongoTokenStore) getTokenDataFor(key, val string) (oauth2.TokenInfo, error) {
	var td tokenData
	var token models.Token
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := ts.dbclient.Database(ts.database).Collection(ts.collection)
	err := coll.FindOne(ctx, bson.D{{key, val}}).Decode(&td)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrorNoResult
		}
		return nil, err
	}
	err = json.Unmarshal(td.Data, &token)
	return &token, err
}

func (ts *MongoTokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	return ts.getTokenDataFor("code_id", code)
}

func (ts *MongoTokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	return ts.getTokenDataFor("access_id", access)
}

func (ts *MongoTokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return ts.getTokenDataFor("refresh_id", refresh)
}
