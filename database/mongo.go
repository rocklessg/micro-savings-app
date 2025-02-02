package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	MongoClient = client
	log.Println("Connected to MongoDB!")
	return nil
}

func DisconnectMongoDB() error {
	if MongoClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := MongoClient.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("Disconnected from MongoDB!")
	return nil
}

// GetCollection returns a reference to a specific collection
func GetCollection(collectionName string) *mongo.Collection {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        panic("DB_NAME is not set in the environment variables")
    }
    return MongoClient.Database(dbName).Collection(collectionName)
}
