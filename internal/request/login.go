package request

import "github.com/gin-gonic/gin"

// LoginParams .
type LoginParams struct {
	Context  *gin.Context
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}
