package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole returns middleware that checks if the authenticated user
// has the required role. Must be used after AuthMiddleware.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Autentikasi diperlukan",
			})
			c.Abort()
			return
		}

		if userRole.(string) != role {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Anda tidak memiliki akses ke resource ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
