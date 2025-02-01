package middlewares

import (
	"micro-savings-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user ID is present in the context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Fetch the user from the database
		user, err := services.GetUserByID(userID.(string))
		if err != nil || user == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden! User not found"})
			c.Abort()
			return
		}

		// Check if the user is an admin
		if !user.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden! Admin access required"})
			c.Abort()
			return
		}
		// Allow the request to proceed
		c.Next()
	}
}