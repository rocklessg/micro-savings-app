package main

import (
	"micro-savings-app/database"
	"micro-savings-app/handlers"
	"micro-savings-app/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	// Connect to MongoDB
	database.ConnectDB(os.Getenv("MONGODB_URI"))

	// Create a new Gin router
	router := gin.Default()

	// Register the user routes
	router.POST("/users/register", handlers.RegisterUser)
	router.POST("/users/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())	
	protected.POST("/transactions/deposit", handlers.Deposit)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}