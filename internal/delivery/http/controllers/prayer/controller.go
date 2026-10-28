package prayer

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/models"
	"github.com/ardie069/kalender-hijriyah/internal/prayer"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/timezone"
	"github.com/ardie069/kalender-hijriyah/internal/validators"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	PrayerSvc *prayer.Service
	DateSvc   *hijri.DateService
	TzSvc     *timezone.Service
}

func NewController(svc *prayer.Service, dateSvc *hijri.DateService, tzSvc *timezone.Service) *Controller {
	return &Controller{
		PrayerSvc: svc,
		DateSvc:   dateSvc,
		TzSvc:     tzSvc,
	}
}

func mapPrayerMethodToHijri(methodStr string) string {
	switch methodStr {
	case "UMM_AL_QURA":
		return "UMM_AL_QURA"
	case "KEMENAG", "JAKIM", "MUIS":
		return "MABIMS"
	default:
		return "MABIMS"
	}
}

func (c *Controller) GetPrayerTimes(ctx *gin.Context) {
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")
	if latStr == "" || lonStr == "" {
		ctx.JSON(400, gin.H{"error": "Parameter 'lat' dan 'lon' wajib disertakan."})
		return
	}
	dateStr := ctx.Query("date")
	methodStr := ctx.DefaultQuery("method", "KEMENAG")
	madhabStr := ctx.DefaultQuery("madhab", "shafii")
	highLatStr := ctx.DefaultQuery("high_lat", "ANGLE_BASED")

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

	var targetDate time.Time
	if dateStr != "" {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD."})
			return
		}
	} else {
		targetDate = time.Now()
	}

	if !prayer.IsValidMethod(methodStr) {
		ctx.JSON(400, gin.H{
			"error":         "Metode tidak dikenali.",
			"valid_methods": []string{"KEMENAG", "MWL", "ISNA", "EGYPTIAN", "UMM_AL_QURA", "KARACHI", "TEHRAN", "JAKIM", "MUIS"},
		})
		return
	}

	cfg := prayer.GetPreset(methodStr)
	cfg.Madhab = prayer.ParseMadhab(madhabStr)
	cfg.HighLatMethod = prayer.ParseHighLat(highLatStr)

	loc := c.TzSvc.GetLocation(lat, lon)

	times, err := c.PrayerSvc.GetPrayerTimes(targetDate, lat, lon, cfg)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Gagal menghitung waktu sholat: %v", err)})
		return
	}

	localDate, tzName := c.TzSvc.GetLocalTimeInfo(targetDate, lat, lon)

	hijriMethod := mapPrayerMethodToHijri(methodStr)
	hijriTargetDate := c.DateSvc.GetHijriTargetDate(targetDate, lat, lon)
	hijriDate := c.DateSvc.ResolveDynamicHijriDate(hijriMethod, hijriTargetDate, lat, lon)

	resp := models.PrayerResponse{}
	resp.Location.Latitude = lat
	resp.Location.Longitude = lon
	resp.Location.Timezone = tzName
	resp.Date.Gregorian = localDate
	resp.Date.Hijri = fmt.Sprintf("%d %s %d", hijriDate.Day, hijriDate.MonthName, hijriDate.Year)

	resp.Method.Name = cfg.Name
	resp.Method.Madhab = cfg.Madhab.Name()
	if math.Abs(lat) >= cfg.HighLatThreshold && cfg.HighLatMethod != nil {
		resp.Method.HighLat = cfg.HighLatMethod.Name()
	}

	resp.Times.Fajr = times.Fajr.In(loc).Format("15:04")
	resp.Times.Sunrise = times.Sunrise.In(loc).Format("15:04")
	resp.Times.Dhuhr = times.Dhuhr.In(loc).Format("15:04")
	resp.Times.Asr = times.Asr.In(loc).Format("15:04")
	resp.Times.Maghrib = times.Maghrib.In(loc).Format("15:04")
	resp.Times.Isha = times.Isha.In(loc).Format("15:04")
	resp.Times.Midnight = times.Midnight.In(loc).Format("15:04")
	resp.Times.ThirdNight = times.ThirdNight.In(loc).Format("15:04")

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
