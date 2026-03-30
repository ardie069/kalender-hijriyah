package hijri

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// ResolveDynamicHijriDate resolves the Hijri date for a specific method on the target day.
func (s *DateService) ResolveDynamicHijriDate(m string, targetDay time.Time, lat, lon float64) models.HijriDate {
	// 1. Normalize targetDay to Noon for stable arithmetic comparison.
	targetNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)

	// 2. Pass 1: Find the most recent or current Ijtima.
	// We search at targetDay to find the Ijtima governing the current or upcoming month.
	ijtimaRecent, _ := s.Cal.FindIjtima(targetNoon)
	monthStartRecent := s.evalMonthStart(m, ijtimaRecent, lat, lon)

	var finalMonthStart time.Time

	// 3. Logic: If targetNoon is BEFORE monthStartRecent, it belongs to the previous lunation.
	if targetNoon.Before(monthStartRecent) {
		// Must use the previous lunation's Ijtima.
		ijtimaPrev, _ := s.Cal.FindIjtima(ijtimaRecent.AddDate(0, 0, -29))
		finalMonthStart = s.evalMonthStart(m, ijtimaPrev, lat, lon)
	} else {
		finalMonthStart = monthStartRecent
	}

	// 4. Calculate hDay (1-indexed).
	daysElapsed := int(math.Round(targetNoon.Sub(finalMonthStart).Hours() / 24.0))
	hDay := daysElapsed + 1

	// 5. Determine Month and Year based on cumulative lunations from epoch.
	// 1 Muharram 1 H has JD 1948439.5.
	monthStartJD := float64(finalMonthStart.Unix())/86400.0 + 2440587.5
	monthsSinceEpoch := int(math.Round((monthStartJD - 1948439.5) / 29.53059))

	evalYear := (monthsSinceEpoch / 12) + 1
	evalMonth := (monthsSinceEpoch % 12) + 1

	// Safety check for hDay (should be 1-30 in this logic)
	// If hDay > 30, it might be due to a very long month (astronomically rare but possible).
	if hDay > 30 {
		hDay -= 30
		evalMonth++
		if evalMonth > 12 {
			evalMonth = 1
			evalYear++
		}
	}

	return models.HijriDate{
		Day:       hDay,
		Month:     evalMonth,
		MonthName: calendar.MonthNames[evalMonth-1],
		Year:      evalYear,
		IsTabular: false,
		// Reference Ijtima for analytical transparency if needed
	}
}

// evalMonthStart evaluates when a specific Hijri month starts given its governing Ijtima.
func (s *DateService) evalMonthStart(m string, ijtima time.Time, lat, lon float64) time.Time {
	sunsetCheckDate := time.Date(ijtima.Year(), ijtima.Month(), ijtima.Day(), 12, 0, 0, 0, time.UTC)
	var isNewMonth bool

	switch m {
	case "UMM_AL_QURA":
		meccaLat, meccaLon := 21.4225, 39.8262
		sunsetCheck, _ := s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		if sunsetCheck.Before(ijtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		}
		moonsetMecca, _ := s.Astro.GetMoonset(sunsetCheck, meccaLat, meccaLon)
		isNewMonth = calendar.IsUmmAlQura(ijtima, sunsetCheck, moonsetMecca)

	case "MABIMS":
		sabangLat, sabangLon := 5.89, 95.32
		sunsetCheck, _ := s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		if sunsetCheck.Before(ijtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		}
		moonsetCheck, _ := s.Astro.GetMoonset(sunsetCheck, sabangLat, sabangLon)
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, sabangLat, sabangLon)
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck, sunsetCheck, moonsetCheck, ijtima)
		isNewMonth = resLocal.IsNewMonth

	case "KHGT":
		khgtRes := s.Cal.ScanGlobalKHGT(sunsetCheckDate, ijtima)
		isNewMonth = khgtRes.IsGlobalValid
		if !isNewMonth {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			khgtRes2 := s.Cal.ScanGlobalKHGT(sunsetCheckDate, ijtima)
			isNewMonth = khgtRes2.IsGlobalValid
		}

	default:
		// Fallback for unknown methods
		khgtRes := s.Cal.ScanGlobalKHGT(sunsetCheckDate, ijtima)
		isNewMonth = khgtRes.IsGlobalValid
	}

	if isNewMonth {
		return sunsetCheckDate.AddDate(0, 0, 1)
	}
	return sunsetCheckDate.AddDate(0, 0, 2)
}
