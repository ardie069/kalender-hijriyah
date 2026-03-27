package hijri

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// ResolveDynamicHijriDate resolves the Hijri date for a specific method on the target day.
func (s *DateService) ResolveDynamicHijriDate(m string, targetDay time.Time, lat, lon float64) models.HijriDate {
	// Pijakan dasar: Gunakan tabular untuk tahu kita di bulan dan tahun berapa secara kasar.
	// Tabular cukup akurat dalam margin +/- 2 hari.
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	baseH := calendar.GetTabularHijri(stableNoon)

	// 1. Tentukan Ijtima yang mendasari bulan ini.
	// FindPreviousIjtima akan mencari dalam rentang [targetDay-30, targetDay].
	prevIjtima, _ := s.Cal.FindPreviousIjtima(targetDay)

	// 1. Kita ambil tanggal dasar Ijtima (H-0)
	sunsetCheckDate := time.Date(prevIjtima.Year(), prevIjtima.Month(), prevIjtima.Day(), 12, 0, 0, 0, time.UTC)
	var sunsetCheck time.Time
	var isNewMonth bool

	// 2. Tentukan Tanggal Evaluasi (sunsetCheckDate) berdasarkan Kriteria
	switch m {
	case "UMM_AL_QURA":
		meccaLat, meccaLon := 21.4225, 39.8262
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		}
		moonsetMecca, _ := s.Astro.GetMoonset(sunsetCheck, meccaLat, meccaLon)
		isNewMonth = calendar.IsUmmAlQura(prevIjtima, sunsetCheck, moonsetMecca)

	case "MABIMS":
		sabangLat, sabangLon := 5.89, 95.32
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		}
		moonsetCheck, _ := s.Astro.GetMoonset(sunsetCheck, sabangLat, sabangLon)
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, sabangLat, sabangLon)
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck, sunsetCheck, moonsetCheck, prevIjtima)
		isNewMonth = resLocal.IsNewMonth

	case "KHGT":
		// KHGT tidak bergantung pada satu sunset spesifik, melainkan global.
		khgtRes := s.Cal.ScanGlobalKHGT(sunsetCheckDate, prevIjtima)
		isNewMonth = khgtRes.IsGlobalValid
		if !isNewMonth {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			khgtRes2 := s.Cal.ScanGlobalKHGT(sunsetCheckDate, prevIjtima)
			isNewMonth = khgtRes2.IsGlobalValid
		}

	default:
		// Metode Lokal (Wujudul Hilal, dll)
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, lat, lon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, lat, lon)
		}
		moonsetCheck, _ := s.Astro.GetMoonset(sunsetCheck, lat, lon)
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, lat, lon)
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck, sunsetCheck, moonsetCheck, prevIjtima)
		isNewMonth = resLocal.IsNewMonth
	}

	// 3. Konklusi: Kapan tanggal 1-nya?
	var monthStartDate time.Time
	if isNewMonth {
		monthStartDate = sunsetCheckDate.AddDate(0, 0, 1)
	} else {
		monthStartDate = sunsetCheckDate.AddDate(0, 0, 2)
	}

	// 4. Hitung Days Elapsed
	targetDayNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	monthStartNoon := time.Date(monthStartDate.Year(), monthStartDate.Month(), monthStartDate.Day(), 12, 0, 0, 0, time.UTC)

	daysElapsed := int(math.Round(targetDayNoon.Sub(monthStartNoon).Hours() / 24.0))

	hDay := daysElapsed + 1

	evalMonth := baseH.Month
	evalYear := baseH.Year

	// Jika hDay < 1, berarti kita masih di bulan sebelumnya
	if hDay < 1 {
		evalMonth--
		if evalMonth < 1 {
			evalMonth = 12
			evalYear--
		}
		hDay += 30
	} else if hDay > 30 {
		hDay -= 30
		evalMonth++
		if evalMonth > 12 {
			evalMonth = 1
			evalYear++
		}
	} else if hDay > 29 && !isNewMonth {
		if baseH.Day == 1 {
			evalMonth--
			if evalMonth < 1 {
				evalMonth = 12
				evalYear--
			}
		}
	}

	return models.HijriDate{
		Day:       hDay,
		Month:     evalMonth,
		MonthName: calendar.MonthNames[evalMonth-1],
		Year:      evalYear,
		IsTabular: false,
	}
}
