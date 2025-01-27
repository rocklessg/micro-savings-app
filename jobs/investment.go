package jobs

import (
	"context"
	"fmt"
	"micro-savings-app/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AllocateIdleBalances(db *mongo.Database) {
	collection := db.Collection("users")
	transactionCollection := db.Collection("transactions")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the idle period (e.g., 30 days)
	// Idle period is Period when savings_balance is not used 
	// for any transaction
	idlePeriod := time.Hour * 24 * 30
	now := time.Now()

	// Find users with idle balances
	filter := bson.M{
		"savings_balance": bson.M{"$gt": 0}, // Positive balance
		"last_transaction_at": bson.M{
			"$lt": now.Add(-idlePeriod), // Last transaction older than 30 days
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("Failed to find users with idle savings_balance: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	// Process each idle user
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			fmt.Printf("Failed to decode user: %v\n", err)
			continue
		}

		// Move funds to the investment balance
		transferAmount := user.SavingsBalance
		update := bson.M{
			"$set": bson.M{
				"savings_balance": 0,
				"investment_balance": user.InvestmentBalance + transferAmount,
				"updated_at": now,
			},
		}

		_, err := collection.UpdateByID(ctx, user.ID, update)
		if err != nil {
			fmt.Printf("Failed to update user: %v\n", err)
			continue
		}

		// Log the investment allocation as a transaction
		transaction := models.Transaction{
			ID:        primitive.NewObjectID(),
			UserID:    user.ID,
			Type:      string(models.Investment),
			Amount:    transferAmount,
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = transactionCollection.InsertOne(ctx, transaction)
		if err != nil {
			fmt.Printf("Failed to log transaction: %v\n", err)
			continue
		}

		fmt.Printf("Allocated %v to investments for user %v\n", transferAmount, user.ID.Hex())
	}
}