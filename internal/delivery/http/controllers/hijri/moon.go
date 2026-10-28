package hijri

import (
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/validators"
	"github.com/gin-gonic/gin"
)

func (c *Controller) GetTelemetry(ctx *gin.Context) {
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

	now := time.Now()
	tel, err := c.Adapter.GetMoonTelemetry(now, lat, lon)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Gagal dapet data langit"})
		return
	}

	ijtima, err := c.DateSvc.Ephem.FindPreviousIjtima(now)
	if err == nil {
		tel.AgeHours = now.Sub(ijtima).Hours()
	}

	ctx.JSON(200, tel)
}
