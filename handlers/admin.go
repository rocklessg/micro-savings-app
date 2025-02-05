package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"micro-savings-app/database"
	"micro-savings-app/models"
	"micro-savings-app/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//var usersCollection = database.GetCollection("users")

func RegisterAdmin(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		SecretKey   string `json:"secret_key" binding:"required"` // Admin secret Key
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email is already registered
	usersCollection := database.GetCollection("users")
	var existingUser models.User
	err := usersCollection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered. Kindly contact admin."})
		return
	}

	// verify admin secret key
	if request.SecretKey != os.Getenv("ADMIN_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the admin user document
	admin := models.User{
		Name:              request.Name,
		Email:             request.Email,
		PasswordHash:      string(hashedPassword),
		SavingsBalance:    0,
		InvestmentBalance: 0,
		IsAdmin: 		   true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	_, err = usersCollection.InsertOne(context.Background(), admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin user registered successfully"})
}

// This is used to make an existing user an admin
func MakeAdmin(c *gin.Context) {
	var request struct {
		UserId   string `json:"user_id" binding:"required"`
		SecretKey   string `json:"secret_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the user from the database
	user, err := services.GetUserByID(request.UserId)
	if err != nil || user == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	// Check if the user is an admin
	if user.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden! user is already an admin."})
		c.Abort()
		return
	}

	// verify admin secret key
	if request.SecretKey != os.Getenv("ADMIN_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
	}

	// update user to admin
    var usersCollection = database.GetCollection("users")

	_, err = usersCollection.UpdateOne(context.Background(),
				bson.M{"_id": user.ID},
				bson.M{"$set": bson.M{
					"is_admin": true,
					"updated_at": time.Now(),
				}})
				
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an admin user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin user created successfully"})
}

// This is used to remove admin rights from a user
func RemoveAdmin(c *gin.Context) {
	var request struct {
		UserId   string `json:"user_id" binding:"required"`
	}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Fetch the user from the database
		user, err := services.GetUserByID(request.UserId)
		if err != nil || user == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if the user is not an admin
		if !user.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not an admin."})
			c.Abort()
			return		
		}

		// update. remove user from being an admin
        usersCollection := database.GetCollection("users")

	_, err = usersCollection.UpdateOne(context.Background(),
				bson.M{"_id": user.ID},
				bson.M{"$set": bson.M{"is_admin": false}})
				
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove a user from being an admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User removed from an admin successfully"})
}

// This function handles admin dashboard statistics
func AdminDashboard(c *gin.Context) {
	// fetch stats from the database
	usersCollection := database.GetCollection("users")
	totalUsers, err := usersCollection.CountDocuments(context.Background(), bson.M{"is_admin": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total users"})
		return
	}

	transactionsCollection := database.GetCollection("transactions")
	totalDeposits, err := transactionsCollection.CountDocuments(context.Background(), bson.M{"type": "deposit"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total deposits"})
		return
	}

	totalWithdrawals, err := transactionsCollection.CountDocuments(context.Background(), bson.M{"type": "withdrawal"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total withdrawals"})
		return
	}

	// return the stats
	c.JSON(http.StatusOK, gin.H{
		"total_users":       totalUsers,
		"total_deposits":    totalDeposits,
		"total_withdrawals": totalWithdrawals,
	})
}

// This is admin get user by ID
func AdminGetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("user_id") // fetch the user id from the url (query param)
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required"})
			c.Abort()
			return
		}

		user, err := services.GetUserByID(userId)
		if err != nil || user == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden! User not found"})
			c.Abort()
			return
		}
		
		// Return the user details
		c.JSON(http.StatusOK, user)
	}
}