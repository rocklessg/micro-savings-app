package handlers

import (
	"context"
	"micro-savings-app/database"
	"micro-savings-app/models"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Deposit handles user deposits into savings
func Deposit(c *gin.Context) {
	var request struct {
		UserID string  `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the user ID
	userID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch the user's current balance
	usersCollection := database.GetCollection("users")
	var user models.User

	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Update the user's savings balance
	newBalance := user.SavingsBalance + request.Amount

	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"savings_balance": newBalance}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update savings balance"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message":         "Deposit successful",
		"new_balance":     newBalance,
		"previous_balance": user.SavingsBalance,
	})
}