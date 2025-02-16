package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TestMongoClient *mongo.Client

func ConnectTestDB() {

	err := godotenv.Load("../.env") // Adjust path if needed
	if err != nil {
		panic("Failed to load .env file")
	}

	uri := os.Getenv("TEST_MONGO_URI") // Use a separate test DB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	TestMongoClient = client
	fmt.Println("Connected to Test MongoDB")
}

func GetTestCollection(name string) *mongo.Collection {
	return TestMongoClient.Database("micro-savings-test_db").Collection(name)
}

func DisconnectTestDB() {
	if err := TestMongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}