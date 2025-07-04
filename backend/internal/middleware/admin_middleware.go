package middleware

import (
	"movie-api/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userModel, ok := user.(*model.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			c.Abort()
			return
		}

		if userModel.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
