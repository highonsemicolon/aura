package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPermission(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"allowed": false})
}
