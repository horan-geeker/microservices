package auth

import "github.com/gin-gonic/gin"

// Login .
func Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hi",
	})
	return
}
