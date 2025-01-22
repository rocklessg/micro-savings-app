package main

import (
	"os"
	"micro-savings-app/database"
	"micro-savings-app/handlers"

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

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}