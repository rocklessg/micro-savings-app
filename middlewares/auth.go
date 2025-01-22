package middlewares

import (
	"net/http"
	"micro-savings-app/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks if the request is authenticated
// Verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! Missing Authorization header"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := services.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! Invalid token"})
			c.Abort()
			return
		}

		// Pass the user ID to the next handler
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}