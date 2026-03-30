package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"github.com/ardie069/kalender-hijriyah/core/hijri"
	"github.com/gin-gonic/gin"
)

type HijriHandler struct {
	DateSvc     *hijri.DateService
	CalendarSvc *hijri.CalendarService
	Adapter     *astronomy.Adapter
	startTime   time.Time
}

func NewHijriHandler(dateSvc *hijri.DateService, calSvc *hijri.CalendarService, adapter *astronomy.Adapter) *HijriHandler {
	return &HijriHandler{
		DateSvc:     dateSvc,
		CalendarSvc: calSvc,
		Adapter:     adapter,
		startTime:   time.Now(),
	}
}

// Ping for health check
func (h *HijriHandler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Ready",
		"message": "Engine NASA lu udah panas, Bang! 🔥",
		"uptime":  time.Since(h.startTime).String(),
	})
}

// GetHijriDate - Prediksi Kalender 4 Metode
func (h *HijriHandler) GetHijriDate(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	var targetDate time.Time

	if dateStr == "" {
		// AMBIL WAKTU SEKARANG (Real-time)
		targetDate = time.Now().UTC()
	} else {
		// PARSING TANGGAL DARI USER (YYYY-MM-DD atau ISO8601)
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

	result := h.DateSvc.GetFullCalendarInfo(targetDate, lat, lon)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (h *HijriHandler) SearchDate(ctx *gin.Context) {
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

	// Langsung panggil logic Tabular (Tanpa butuh data NASA karena ini murni itungan kitab)
	result := h.DateSvc.GetTabularOnly(targetDate)

	ctx.JSON(200, gin.H{
		"status": "success",
		"query":  dateStr,
		"data":   result,
	})
}

// GetTelemetry - Telemetry Hilal
func (h *HijriHandler) GetTelemetry(ctx *gin.Context) {
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

	now := time.Now()
	tel, err := h.Adapter.GetMoonTelemetry(now, lat, lon)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Gagal dapet data langit"})
		return
	}

	// Calculate Moon Age (Hours since previous Ijtima)
	ijtima, err := h.DateSvc.Cal.FindPreviousIjtima(now)
	if err == nil {
		tel.AgeHours = now.Sub(ijtima).Hours()
	}

	ctx.JSON(200, tel)
}

func (h *HijriHandler) GetYearlyCalendar(ctx *gin.Context) {
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

	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")
	if latStr == "" || lonStr == "" {
		ctx.JSON(400, gin.H{"error": "Parameter 'lat' dan 'lon' wajib disertakan."})
		return
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	method := ctx.DefaultQuery("method", "KHGT")

	result := h.CalendarSvc.GetYearlyCalendar(year, lat, lon, method)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (h *HijriHandler) GetGregorianMonth(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	monthStr := ctx.Query("month")
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")

	if yearStr == "" || monthStr == "" || latStr == "" || lonStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, month, lat, dan lon wajib disertakan."})
		return
	}

	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	result := h.DateSvc.GetGregorianMonthInfo(year, month, lat, lon)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}
