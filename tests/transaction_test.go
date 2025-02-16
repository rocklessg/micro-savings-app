package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"micro-savings-app/handlers"
	"micro-savings-app/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupUserForTransaction() primitive.ObjectID {
	// Insert a test user with a balance
	userCollection := database.GetTestCollection("users")
	user := bson.M{
		"email":           "test@example.com",
		"savings_balance": 1000.0,
	}
	res, _ := userCollection.InsertOne(context.Background(), user)
	return res.InsertedID.(primitive.ObjectID)
}

func TestDeposit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := setupUserForTransaction()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"amount": 500}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/transactions/deposit", bytes.NewBuffer([]byte(requestBody)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", userID.Hex()) // Simulate authentication

	handlers.Deposit(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var user map[string]interface{}
	userCollection := database.GetTestCollection("users")
	_ = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	assert.Equal(t, 1500.0, user["savings_balance"])
}

func TestWithdraw(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := setupUserForTransaction()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody := `{"amount": 500}`
	c.Request, _ = http.NewRequest("POST", "/transactions/withdraw", bytes.NewBuffer([]byte(requestBody)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", userID.Hex()) // Simulate authentication

	handlers.Withdraw(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var user map[string]interface{}
	userCollection := database.GetTestCollection("users")
	_ = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	assert.Equal(t, 500.0, user["savings_balance"])
}
