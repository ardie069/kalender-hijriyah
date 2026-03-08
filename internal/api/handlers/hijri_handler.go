package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/services"
	"github.com/gin-gonic/gin"
)

type HijriHandler struct {
	Service *services.HijriService
	Adapter *astronomy.Adapter
}

func NewHijriHandler(service *services.HijriService, adapter *astronomy.Adapter) *HijriHandler {
	return &HijriHandler{
		Service: service,
		Adapter: adapter,
	}
}

// Ping for health check
func (h *HijriHandler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Ready",
		"message": "Engine NASA lu udah panas, Bang! 🔥",
		"uptime":  time.Since(time.Now()).String(),
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
		// Kita coba format RFC3339 dulu, kalau gagal fallback ke 2006-01-02
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

	lat, _ := strconv.ParseFloat(ctx.DefaultQuery("lat", "-7.98"), 64)
	lon, _ := strconv.ParseFloat(ctx.DefaultQuery("lon", "112.62"), 64)

	result := h.Service.GetFullCalendarInfo(targetDate, lat, lon)
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
	result := h.Service.GetTabularOnly(targetDate)

	ctx.JSON(200, gin.H{
		"status": "success",
		"query":  dateStr,
		"data":   result,
	})
}

// GetTelemetry - Telemetry Hilal
func (h *HijriHandler) GetTelemetry(ctx *gin.Context) {
	lat, _ := strconv.ParseFloat(ctx.DefaultQuery("lat", "-7.98"), 64)
	lon, _ := strconv.ParseFloat(ctx.DefaultQuery("lon", "112.62"), 64)

	tel, err := h.Adapter.GetMoonTelemetry(time.Now(), lat, lon)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Gagal dapet data langit"})
		return
	}
	ctx.JSON(200, tel)
}
