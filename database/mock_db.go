package database

import (
	"context"
	"log"
	"os"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestMongoClient *mongo.Client

func ConnectTestDB() {
	uri := os.Getenv("TEST_MONGO_URI") // Use a separate test DB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	TestMongoClient = client
	fmt.Println("Connected to Test MongoDB")
}

func GetTestCollection(name string) *mongo.Collection {
	return TestMongoClient.Database("test_db").Collection(name)
}

func DisconnectTestDB() {
	if err := TestMongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}