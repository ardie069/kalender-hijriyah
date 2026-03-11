package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/models"
	"github.com/ardie069/kalender-hijriyah/internal/prayer"
	"github.com/ardie069/kalender-hijriyah/internal/services"
	"github.com/gin-gonic/gin"
)

type PrayerHandler struct {
	PrayerService *prayer.PrayerService
	HijriService  *services.HijriService
}

func NewPrayerHandler(prayerService *prayer.PrayerService, hijriService *services.HijriService) *PrayerHandler {
	return &PrayerHandler{
		PrayerService: prayerService,
		HijriService:  hijriService,
	}
}

func (h *PrayerHandler) GetPrayerTimes(c *gin.Context) {
	// 1. Ambil Query Params
	latStr := c.DefaultQuery("lat", "-7.98")
	lonStr := c.DefaultQuery("lon", "112.62")
	dateStr := c.Query("date") // Format: 2006-01-02

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	var targetDate time.Time
	if dateStr != "" {
		targetDate, _ = time.Parse("2006-01-02", dateStr)
	} else {
		targetDate = time.Now()
	}

	// 2. Tentukan Timezone Berdasarkan Lokasi
	loc := h.HijriService.GetLocation(lat, lon)

	// 3. Panggil Service (Pake logic SPICE yang udah kita bikin)
	times := h.PrayerService.GetPrayerTimes(targetDate, lat, lon)

	// 4. Ambil Info Timezone & Hijriah buat melengkapi response
	localDate, tzName := h.HijriService.GetLocalTimeInfo(targetDate, lat, lon)
	hijri := h.HijriService.ResolveDynamicHijriDate("MABIMS", targetDate, lat, lon)

	// 5. Susun JSON
	resp := models.PrayerResponse{}
	resp.Location.Latitude = lat
	resp.Location.Longitude = lon
	resp.Location.Timezone = tzName
	resp.Date.Gregorian = localDate
	resp.Date.Hijri = fmt.Sprintf("%d %s %d", hijri.Day, hijri.MonthName, hijri.Year)

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
