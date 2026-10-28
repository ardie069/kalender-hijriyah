package hijri

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetVisibilityMap(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	method := ctx.DefaultQuery("method", "KHGT")

	var targetDate time.Time
	if dateStr == "" {
		targetDate = time.Now().UTC()
	} else {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD."})
			return
		}
		targetDate = parsedDate
	}

	result, err := c.VisSvc.GenerateVisibilityGrid(targetDate, method)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat peta visibilitas: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}
