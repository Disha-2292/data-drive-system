// middleware/role_middleware.go

package middleware

import (
	"net/http"

	"github.com/Disha-2292/data-drive-system/config"
	"github.com/Disha-2292/data-drive-system/models"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from the context
		userID := c.MustGet("user_id").(uint)

		var user models.User
		if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if the user has the admin role
		if user.Role.Name != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// If role is admin, proceed
		c.Next()
	}
}

func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user from the context
		userID := c.MustGet("user_id").(uint)

		var user models.User
		if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if the user has the user role
		if user.Role.Name != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// If role is user, proceed
		c.Next()
	}
}
