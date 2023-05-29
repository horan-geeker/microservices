package system

import "github.com/gin-gonic/gin"

// Health .
func Health(c *gin.Context) (map[string]any, int, error) {
	return map[string]interface{}{
		"status": "UP",
	}, 0, nil
}
