package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

// GetFullCalendarInfo returns multi-method Hijri date info for the given time and location.
func (s *DateService) GetFullCalendarInfo(t time.Time, lat, lon float64) models.HijriResponse {
	tUTC := t.UTC()

	targetDay := s.GetHijriTargetDate(tUTC, lat, lon)
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	currentH := calendar.GetTabularHijri(stableNoon)

	resp := models.HijriResponse{
		GregorianDate: tUTC,
		Location:      models.LocationInfo{Lat: lat, Lon: lon},
		Methods:       make(map[string]models.MethodResult),
	}

	realtimeTel, _ := s.Astro.GetMoonTelemetry(tUTC, lat, lon)

	methodList := []string{"TABULAR", "MABIMS", "KHGT", "UMM_AL_QURA", "SAUDI"}

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
		}

		if m != "TABULAR" {
			alt := realtimeTel.Altitude
			elong := realtimeTel.Elongation

			if m == "UMM_AL_QURA" {
				meccaTel, err := s.Astro.GetMoonTelemetry(tUTC, 21.4225, 39.8262)
				if err == nil {
					alt = meccaTel.Altitude
					elong = meccaTel.Elongation
				}
			}

			if m == "SAUDI" {
				// Default to Sudair telemetry for global response reference
				sudairTel, err := s.Astro.GetMoonTelemetry(tUTC, 25.5950, 45.6339)
				if err == nil {
					alt = sudairTel.Altitude
					elong = sudairTel.Elongation
				}
			}

			if m == "MABIMS" {
				sabangTel, err := s.Astro.GetMoonTelemetry(tUTC, 5.89, 95.32)
				if err == nil {
					refAlt := sabangTel.Altitude
					refElong := sabangTel.Elongation
					result.ReferenceAltitude = &refAlt
					result.ReferenceElongation = &refElong
				}
			}

			result.CurrentAltitude = &alt
			result.CurrentElongation = &elong
		}

		if m != "TABULAR" && result.HijriDate.Day <= 29 {
			daysTo29 := 29 - result.HijriDate.Day
			searchDateMethod := targetDay.AddDate(0, 0, daysTo29)
			searchDateMethod = time.Date(searchDateMethod.Year(), searchDateMethod.Month(), searchDateMethod.Day(), 12, 0, 0, 0, time.UTC)

			pred, err := s.CalculateMethodPrediction(m, searchDateMethod, lat, lon)
			if err == nil {
				result.Prediction = pred
			}

			if m == "MABIMS" {
				localPred, err := s.CalculateMethodPrediction("MABIMS_LOCAL", searchDateMethod, lat, lon)
				if err == nil {
					if pred != nil {
						localPred.IsNewMonth = pred.IsNewMonth
					}
					result.LocalPrediction = localPred
				}
			}
		}

		resp.Methods[m] = result
	}

	return resp
}
