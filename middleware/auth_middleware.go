package middleware

import (
	"net/http"
	"weather/utils"

	"github.com/gin-gonic/gin"
)

// AuthenticationMiddleware checks if the user has a valid JWT token
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token2"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
