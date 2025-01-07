package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserIDMiddleware extracts the user ID from the request headers
func UserIDMiddleware(c *gin.Context) {
	userID := c.GetHeader("X-User-ID") // NGINX adds this bt resolving JWT
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user ID not provided"})
		return
	}
	c.Set("userID", userID) // Pass the user ID to handlers via context
	c.Next()
}
