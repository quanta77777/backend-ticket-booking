package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing or invalid role information"})
			c.Abort()
			return
		}

		actualRole, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role information type"})
			c.Abort()
			return
		}

		if actualRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
