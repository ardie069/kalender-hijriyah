package hijri

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Ready",
		"message": "Engine NASA lu udah panas, Bang! 🔥",
		"uptime":  time.Since(c.startTime).String(),
	})
}
