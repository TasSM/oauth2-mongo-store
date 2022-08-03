package util

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongo_test_host = "127.0.0.1"
	mongo_test_port = "27017"
	mongo_test_db   = "store-testdb"
	mongo_test_user = "test"
	mongo_test_pass = "test"
)

func ConnectLocalMongo() *mongo.Client {
	mongocs := fmt.Sprintf("mongodb://%s:%s@%s:%s/", mongo_test_user, mongo_test_pass, mongo_test_host, mongo_test_port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongocs))
	if err != nil {
		panic(err)
	}
	log.Printf("connected to test mongo server: %s", mongocs)
	return client
}
