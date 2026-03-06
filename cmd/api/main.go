package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Engine NASA (Pake NewEphemerisManager biar aman)
	// Fungsi ini bakal nge-loop dan mastiin semua kernel ke-load tanpa skip error
	manager, err := astronomy.NewEphemerisManager(
		"data/naif0012.tls",
		"data/de440.bsp",
		"data/pck00011.tpc",
	)
	if err != nil {
		log.Fatalf("❌ NASA Engine Failure: %v", err)
	}

	// 2. Setup Layers (Gunakan Factory Function biar konsisten)
	adapter := astronomy.GetAdapter(manager)
	logic := calendar.NewLogic(adapter, manager)
	service := &services.HijriService{
		Astro: adapter, // Pastiin 'Astro' di struct HijriService pake huruf Gede
		Cal:   logic,   // Pastiin 'Cal' di struct HijriService pake huruf Gede
	}

	fmt.Println("🚀 Kalender Hijriyah Engine: STANDBY ON PORT 8080")
	fmt.Println("🌍 Coordinate System: NASA SPICE Topocentric")

	r := gin.Default()

	// --- ROUTES ---

	// Health Check - Penting buat monitoring Docker nanti
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "Ready",
			"message": "Engine NASA lu udah panas, Bang! 🔥",
			"uptime":  time.Since(time.Now()).String(), // Tambahan info dikit
		})
	})

	// Endpoint Utama: Prediksi Kalender 4 Metode
	r.GET("/api/v1/calendar/hijri", func(ctx *gin.Context) {
		dateStr := ctx.Query("date") // Jangan pake DefaultQuery dulu
		var targetDate time.Time

		if dateStr == "" {
			// AMBIL WAKTU SEKARANG (Real-time)
			targetDate = time.Now().UTC()
		} else {
			// PARSING TANGGAL DARI USER (YYYY-MM-DD)
			// Kita tambahin jam saat ini biar telemetry-nya realtime, bukan 00:00
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Format salah"})
				return
			}
			now := time.Now()
			targetDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
		}

		lat, _ := strconv.ParseFloat(ctx.DefaultQuery("lat", "-7.98"), 64)
		lon, _ := strconv.ParseFloat(ctx.DefaultQuery("lon", "112.62"), 64)

		result := service.GetFullCalendarInfo(targetDate, lat, lon)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
	})

	// Endpoint buat Telemetry (Debug visual di Flutter nanti)
	r.GET("/api/v1/moon/telemetry", func(ctx *gin.Context) {
		lat, _ := strconv.ParseFloat(ctx.DefaultQuery("lat", "-7.98"), 64)
		lon, _ := strconv.ParseFloat(ctx.DefaultQuery("lon", "112.62"), 64)

		tel, err := adapter.GetMoonTelemetry(time.Now(), lat, lon)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Gagal dapet data langit"})
			return
		}
		ctx.JSON(200, tel)
	})

	r.Run(":8080")
}
