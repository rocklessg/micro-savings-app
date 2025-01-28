package main

import (
	"micro-savings-app/database"
	"micro-savings-app/handlers"
	"micro-savings-app/jobs"
	"micro-savings-app/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	// Connect to MongoDB
	database.ConnectDB(os.Getenv("MONGO_URI"))

	// Create a new Gin router
	router := gin.Default()

	// Register the user routes
	router.POST("/users/register", handlers.RegisterUser)
	router.POST("/users/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())	
	protected.POST("/transactions/deposit", handlers.Deposit)
	protected.POST("/transactions/withdraw", handlers.Withdraw)

	// Set up the cron job
	c := cron.New()
	_, err = c.AddFunc("@daily", func() {
		jobs.AllocateIdleBalances(database.MongoClient.Database(os.Getenv("DB_NAME")))
	})
	if err != nil {
		panic("Failed to add cron job: " + err.Error())
	}
	c.Start()

	// Ensure cron stops when the app shuts down
	defer func() {
		c.Stop()
		database.DisconnectMongoDB() // Ensure MongoDB connection is closed
	}()
	
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}