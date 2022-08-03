package store

import (
	"context"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/oauth2.v3/models"
)

// custom mongodb oauth client store for go-oauth2
// non exported client model as it is only required for internal implementation
// implements oauth2.ClientStore with additional operations for CRUD

type OAuthClientStorer interface {
	oauth2.ClientStore
	Set(info oauth2.ClientInfo) error
	RemoveByID(id string) error
}

type MongoClientStore struct {
	dbclient   *mongo.Client
	database   string
	collection string
}

type client struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Secret string             `bson:"secret"`
	Domain string             `bson:"domain"`
	UserID string             `bson:"user_id"`
}

func NewMongoClientStore(dbclient *mongo.Client, dbname string) *MongoClientStore {
	return &MongoClientStore{dbclient: dbclient, database: dbname, collection: "oauth_client"}
}

// Save a client
func (cs *MongoClientStore) Set(info oauth2.ClientInfo) error {
	coll := cs.dbclient.Database(cs.database).Collection(cs.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := client{
		//ID:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		UserID: info.GetUserID(),
	}
	_, err := coll.InsertOne(ctx, data)
	return err
}

// GetByID according to the ID for the client information
func (cs *MongoClientStore) GetByID(id string) (info oauth2.ClientInfo, err error) {
	coll := cs.dbclient.Database(cs.database).Collection(cs.collection)
	res := models.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = coll.FindOne(ctx, bson.D{{"user_id", id}}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrorNoResult
		}
		return nil, err
	}
	return &res, nil
}

// Delete a client by the id
func (ts *MongoClientStore) RemoveByID(id string) error {
	coll := ts.dbclient.Database(ts.database).Collection(ts.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.DeleteOne(ctx, bson.D{{"user_id", id}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return err
}
