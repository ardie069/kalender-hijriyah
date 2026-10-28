package hijri

import (
	"net/http"
	"strconv"

	"github.com/ardie069/kalender-hijriyah/internal/validators"
	"github.com/gin-gonic/gin"
)

func (c *Controller) GetYearlyCalendar(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	if yearStr == "" {
		ctx.JSON(400, gin.H{"error": "Parameter 'year' wajib disertakan."})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Parameter 'year' harus berupa angka."})
		return
	}
	if year < 1 {
		ctx.JSON(400, gin.H{"error": "Parameter 'year' harus lebih besar dari 0."})
		return
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

	method := ctx.DefaultQuery("method", "KHGT")

	result := c.CalendarSvc.GetYearlyCalendar(year, lat, lon, method)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (c *Controller) GetGregorianMonth(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	monthStr := ctx.Query("month")
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")

	if yearStr == "" || monthStr == "" || latStr == "" || lonStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, month, lat, dan lon wajib disertakan."})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Parameter 'year' harus berupa angka."})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Parameter 'month' harus berupa angka."})
		return
	}
	if month < 1 || month > 12 {
		ctx.JSON(400, gin.H{"error": "Parameter 'month' harus berada di rentang 1 s/d 12."})
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

	result := c.DateSvc.GetGregorianMonthInfo(year, month, lat, lon)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}
