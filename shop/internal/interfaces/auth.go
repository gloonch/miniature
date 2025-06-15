package interfaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"miniature/pkg/token" // Assuming this pkg/token is a shared module
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := token.ValidateToken(tokenStr) // This will need to be adjusted if token validation logic changes
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// store user info in context
		c.Set("user_id", claims.UserID)
		// Roles might be different for shops, e.g., "OWNER", "SELLER"
		// For now, keeping it as claims.Role but this might need adjustment
		c.Set("role", claims.Role)
		c.Next()
	}
}
