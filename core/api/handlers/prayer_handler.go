package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/models"
	"github.com/ardie069/kalender-hijriyah/core/prayer"
	"github.com/ardie069/kalender-hijriyah/core/services"
	"github.com/gin-gonic/gin"
)

type PrayerHandler struct {
	PrayerService *services.PrayerService
	HijriService  *services.HijriService
}

func NewPrayerHandler(prayerService *services.PrayerService, hijriService *services.HijriService) *PrayerHandler {
	return &PrayerHandler{
		PrayerService: prayerService,
		HijriService:  hijriService,
	}
}

func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	// 1. Ambil Query Params
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	if latStr == "" || lonStr == "" {
		c.JSON(400, gin.H{"error": "Parameter 'lat' dan 'lon' wajib disertakan."})
		return
	}
	dateStr := c.Query("date")              // Format: 2006-01-02
	methodStr := c.DefaultQuery("method", "KEMENAG")
	madhabStr := c.DefaultQuery("madhab", "shafii")
	highLatStr := c.DefaultQuery("high_lat", "ANGLE_BASED")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parameter 'lat' harus berupa angka desimal."})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Parameter 'lon' harus berupa angka desimal."})
		return
	}

	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD."})
			return
		}
	} else {
		targetDate = time.Now()
	}

	// 2. Validasi dan parsing metode
	if !prayer.IsValidMethod(methodStr) {
		c.JSON(400, gin.H{
			"error":           "Metode tidak dikenali.",
			"valid_methods":   []string{"KEMENAG", "MWL", "ISNA", "EGYPTIAN", "UMM_AL_QURA", "KARACHI", "TEHRAN", "JAKIM", "MUIS"},
		})
		return
	}

	cfg := prayer.GetPreset(methodStr)
	cfg.Madhab = prayer.ParseMadhab(madhabStr)
	cfg.HighLatMethod = prayer.ParseHighLat(highLatStr)

	// 3. Tentukan Timezone Berdasarkan Lokasi
	loc := h.HijriService.GetLocation(lat, lon)

	// 4. Panggil Service
	times, err := h.PrayerService.GetPrayerTimes(targetDate, lat, lon, cfg)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Gagal menghitung waktu sholat: %v", err)})
		return
	}

	// 5. Ambil Info Timezone & Hijriah buat melengkapi response
	localDate, tzName := h.HijriService.GetLocalTimeInfo(targetDate, lat, lon)

	// Petakan metode waktu sholat ke kriteria Kalender Hijriyah yang sesuai
	hijriMethod := "MABIMS" // Default regional
	switch methodStr {
	case "UMM_AL_QURA":
		hijriMethod = "UMM_AL_QURA"
	case "KEMENAG", "JAKIM", "MUIS":
		hijriMethod = "MABIMS"
	}
	hijriTargetDate := h.HijriService.GetHijriTargetDate(targetDate, lat, lon)
	hijri := h.HijriService.ResolveDynamicHijriDate(hijriMethod, hijriTargetDate, lat, lon)

	// 6. Susun JSON
	resp := models.PrayerResponse{}
	resp.Location.Latitude = lat
	resp.Location.Longitude = lon
	resp.Location.Timezone = tzName
	resp.Date.Gregorian = localDate
	resp.Date.Hijri = fmt.Sprintf("%d %s %d", hijri.Day, hijri.MonthName, hijri.Year)

	resp.Method.Name = cfg.Name
	resp.Method.Madhab = prayer.MadhabName(cfg.Madhab)
	if math.Abs(lat) >= cfg.HighLatThreshold && cfg.HighLatMethod != prayer.HIGH_LAT_NONE {
		resp.Method.HighLat = prayer.HighLatName(cfg.HighLatMethod)
	}

	// Format waktu ke string HH:mm dalam waktu LOKAL
	resp.Times.Fajr = times.Fajr.In(loc).Format("15:04")
	resp.Times.Sunrise = times.Sunrise.In(loc).Format("15:04")
	resp.Times.Dhuhr = times.Dhuhr.In(loc).Format("15:04")
	resp.Times.Asr = times.Asr.In(loc).Format("15:04")
	resp.Times.Maghrib = times.Maghrib.In(loc).Format("15:04")
	resp.Times.Isha = times.Isha.In(loc).Format("15:04")
	resp.Times.Midnight = times.Midnight.In(loc).Format("15:04")
	resp.Times.ThirdNight = times.ThirdNight.In(loc).Format("15:04")

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
