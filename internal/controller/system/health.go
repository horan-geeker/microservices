package system

import (
	"github.com/gin-gonic/gin"
)

// Health .
func (s *SystemController) Health(c *gin.Context) (map[string]any, error) {
	return map[string]interface{}{
		"status": "UP",
	}, nil
}
