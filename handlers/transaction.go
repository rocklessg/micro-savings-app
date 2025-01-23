package handlers

import (
	"context"
	"micro-savings-app/database"
	"micro-savings-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Deposit handles user deposits into savings
func Deposit(c *gin.Context) {
	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from JWT claims
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! user not authenticated"})
		return
	}

	// Parse the user ID
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch the user's current balance
	usersCollection := database.GetCollection("users")
	var user models.User

	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&user)
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
		bson.M{"_id": userObjectID},
		bson.M{"$set": bson.M{"savings_balance": newBalance}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update savings balance"})
		return
	}

	// Log the transaction
	transactionsCollection := database.GetCollection("transactions")
	transaction := models.Transaction{
		ID:        primitive.NewObjectID(),
		UserID:    userObjectID,
		Type:      string(models.Deposit),
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = transactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transaction"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message":         "Deposit successful",
		"new_balance":     newBalance,
		"previous_balance": user.SavingsBalance,
	})
}

// Withdraw handles user withdrawals from savings
func Withdraw(c *gin.Context) {
	var request struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from JWT claims
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! user not authenticated"})
		return
	}

	// Parse the user ID
	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch the user's current balance
	usersCollection := database.GetCollection("users")
	var user models.User

	err = usersCollection.FindOne(context.Background(), bson.M{"_id": userObjectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Check if the user has sufficient balance
	if user.SavingsBalance < request.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Update the user's savings balance
	newBalance := user.SavingsBalance - request.Amount

	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": userObjectID},
		bson.M{"$set": bson.M{"savings_balance": newBalance}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update savings balance"})
		return
	}

	// Log the transaction
	transactionsCollection := database.GetCollection("transactions")
	transaction := models.Transaction{
		ID:        primitive.NewObjectID(),
		UserID:    userObjectID,
		Type:      string(models.Withdrawal),
		Amount:    request.Amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = transactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transaction"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{
		"message":         "Withdrawal successful",
		"withdrawal_amount": request.Amount,
		"previous_balance": user.SavingsBalance,
		"new_balance":     newBalance,
	})
}