package handlers

import (
	"context"
	"micro-savings-app/database"
	"micro-savings-app/models"
	"micro-savings-app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
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
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the user document
	newUser := models.User{
		Name:              request.Name,
		Email:             request.Email,
		PasswordHash:      string(hashedPassword),
		SavingsBalance:    0,
		InvestmentBalance: 0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	_, err = usersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login user handles user login returns JWT token
func Login(c *gin.Context) {
    var request struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Fetch the user document from the database
    usersCollection := database.GetCollection("users")
    var user models.User
    err := usersCollection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Compare the password with the hash (verify the password)
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate JWT
	token, err := services.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}