package services

import (
	"context"
	"errors"
	"time"

	"micro-savings-app/models"
	"micro-savings-app/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserByID fetches a user from the database by ID
func GetUserByID(userID string) (*models.User, error) {
	collection := database.GetCollection("users") // Ensure the correct collection name

	// Convert string ID to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Define a timeout context for the query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query the database
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
