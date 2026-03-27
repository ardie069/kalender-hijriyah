package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
	"github.com/ardie069/kalender-hijriyah/core/timezone"
)

// DateService handles Hijri date resolution across multiple methods.
type DateService struct {
	Tz    *timezone.Service
	Astro *astronomy.Adapter
	Cal   *calendar.Logic
}

// NewDateService creates a new DateService.
func NewDateService(astro *astronomy.Adapter, cal *calendar.Logic, tz *timezone.Service) *DateService {
	return &DateService{
		Tz:    tz,
		Astro: astro,
		Cal:   cal,
	}
}

// GetFullCalendarInfo returns multi-method Hijri date info for the given time and location.
func (s *DateService) GetFullCalendarInfo(t time.Time, lat, lon float64) models.HijriResponse {
	tUTC := t.UTC()

	// 1. Get Base Telemetry and Data
	targetDay := s.GetHijriTargetDate(tUTC, lat, lon)
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	currentH := calendar.GetTabularHijri(stableNoon)

	resp := models.HijriResponse{
		GregorianDate: tUTC,
		Location:      models.LocationInfo{Lat: lat, Lon: lon},
		Methods:       make(map[string]models.MethodResult),
	}

	realtimeTel, _ := s.Astro.GetMoonTelemetry(tUTC, lat, lon)

	// Methods to evaluate
	methodList := []string{"TABULAR", "MABIMS", "WUJUDUL_HILAL", "KHGT", "UMM_AL_QURA"}

	// Process methods sequentially.
	// SPICE is serialized via mutex, so goroutines here only add overhead.
	for _, m := range methodList {
		var result models.MethodResult

		if m == "TABULAR" {
			result.HijriDate = models.HijriDate{
				Day:       currentH.Day,
				Month:     currentH.Month,
				MonthName: currentH.MonthName,
				Year:      currentH.Year,
				IsTabular: true,
			}
		} else {
			result.HijriDate = s.ResolveDynamicHijriDate(m, targetDay, lat, lon)

			// Add Real-time Telemetry
			alt := realtimeTel.Altitude
			elong := realtimeTel.Elongation
			result.CurrentAltitude = &alt
			result.CurrentElongation = &elong

			// 3. Predictions for the 29th of the Hijri Month
			if result.HijriDate.Day <= 29 {
				daysTo29 := 29 - result.HijriDate.Day
				searchDateMethod := targetDay.AddDate(0, 0, daysTo29)
				searchDateMethod = time.Date(searchDateMethod.Year(), searchDateMethod.Month(), searchDateMethod.Day(), 12, 0, 0, 0, time.UTC)

				pred, err := s.calculateMethodPrediction(m, searchDateMethod, lat, lon)
				if err == nil {
					result.Prediction = pred
				}

				if m == "MABIMS" {
					localPred, err := s.calculateMethodPrediction("MABIMS_LOCAL", searchDateMethod, lat, lon)
					if err == nil {
						if pred != nil {
							localPred.IsNewMonth = pred.IsNewMonth
						}
						result.LocalPrediction = localPred
					}
				}
			}
		}
		resp.Methods[m] = result
	}

	return resp
}

// GetTabularOnly returns tabular-only Hijri date without astronomical telemetry.
func (s *DateService) GetTabularOnly(t time.Time) models.MethodResult {
	hDate := calendar.GetTabularHijri(t)
	hDate.IsTabular = true

	return models.MethodResult{
		HijriDate:  hDate,
		Prediction: nil,
	}
}

// GetHijriTargetDate returns the Gregorian date that corresponds to the current Hijri day.
// It handles rollover at sunset (or Isha for critical months on the 29th).
func (s *DateService) GetHijriTargetDate(t time.Time, lat, lon float64) time.Time {
	tUTC := t.UTC()

	sunsetToday, err := s.Astro.GetSunset(tUTC, lat, lon)
	if err != nil {
		sunsetToday = tUTC
	}

	hToday := calendar.GetTabularHijri(tUTC)

	rolloverTime := sunsetToday
	// Critical Months: 8 (Sha'ban), 9 (Ramadan), 11 (Dhu al-Qi'dah)
	// Role over at Isha only on the 29th of these months.
	if (hToday.Month == 8 || hToday.Month == 9 || hToday.Month == 11) && hToday.Day == 29 {
		ishaToday, err := s.Astro.GetIsha(tUTC, lat, lon)
		if err == nil {
			rolloverTime = ishaToday
		}
	}

	if tUTC.After(rolloverTime) {
		return tUTC.Add(24 * time.Hour)
	}
	return tUTC
}
