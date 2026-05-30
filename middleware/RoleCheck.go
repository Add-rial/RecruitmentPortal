package middleware

import (
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		role := userRole.(string)

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}