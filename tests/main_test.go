package tests

import (
	"context"
	"log"
	"micro-savings-app/database"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// Load .env file before running tests
	err := godotenv.Load("../.env") // Adjust path if needed
	if err != nil {
		panic("Failed to load .env file")
	}

	// Connect to MongoDB test database
	clientOptions := options.Client().ApplyURI(os.Getenv("TEST_MONGO_URI")) // Ensure TEST_MONGO_URI is set in .env
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	database.ConnectTestDB() // Custom function to set test DB

	exitCode := m.Run()
	//os.Exit(m.Run()) // Run tests

	// Clean up
	client.Disconnect(context.Background())

	os.Exit(exitCode)
}