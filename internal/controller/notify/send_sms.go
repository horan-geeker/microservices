package notify

import (
	"github.com/gin-gonic/gin"
)

func (n *notifyController) SendSms(c *gin.Context) (map[string]any, error) {
	if err := n.logic.Notify().SendSmsCode(c.Request.Context(), "13571899655", "1234"); err != nil {
		return nil, err
	}
	return nil, nil
}
