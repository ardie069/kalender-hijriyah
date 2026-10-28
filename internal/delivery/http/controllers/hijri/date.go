package hijri

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/validators"
	"github.com/gin-gonic/gin"
)

func (c *Controller) GetHijriDate(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	var targetDate time.Time

	if dateStr == "" {
		targetDate = time.Now().UTC()
	} else {
		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			parsedDate, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD atau RFC3339."})
				return
			}
			now := time.Now()
			targetDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
		} else {
			targetDate = parsedDate.UTC()
		}
	}

	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")
	if latStr == "" || lonStr == "" {
		ctx.JSON(400, gin.H{"error": "Parameter 'lat' dan 'lon' wajib disertakan."})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Parameter 'lat' harus berupa angka desimal."})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Parameter 'lon' harus berupa angka desimal."})
		return
	}
	if err := validators.ValidateCoords(lat, lon); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result := c.DateSvc.GetFullCalendarInfo(targetDate, lat, lon)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (c *Controller) SearchDate(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(400, gin.H{"error": "Masukkan tanggal terlebih dahulu dengan format YYYY-MM-DD!"})
		return
	}

	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Format tanggal salah. Kudu YYYY-MM-DD ya!"})
		return
	}

	result := c.DateSvc.GetTabularOnly(targetDate)

	ctx.JSON(200, gin.H{
		"status": "success",
		"query":  dateStr,
		"data":   result,
	})
}
