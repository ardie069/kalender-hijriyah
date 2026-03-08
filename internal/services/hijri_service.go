package services

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

type HijriService struct {
	Astro *astronomy.Adapter
	Cal   *calendar.Logic
}

func (s *HijriService) GetFullCalendarInfo(t time.Time, lat, lon float64) models.HijriResponse {
	tUTC := t.UTC()

	// 1. Inisialisasi Data Dasar
	sunsetToday, _ := s.Astro.GetSunset(tUTC, lat, lon)
	realtimeTel, _ := s.Astro.GetMoonTelemetry(tUTC, lat, lon)

	// 2. Tentukan Tanggal Base Tabular (Logic Maghrib & Midnight Fix)
	isAfterSunset := tUTC.After(sunsetToday)

	// Set Target Date ke Jam 12:00 UTC pada hari siang yang dituju
	targetDay := tUTC
	if isAfterSunset {
		targetDay = tUTC.Add(24 * time.Hour)
	}
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	currentH := calendar.GetTabularHijri(stableNoon)

	resp := models.HijriResponse{
		GregorianDate: tUTC,
		Location:      models.LocationInfo{Lat: lat, Lon: lon},
		Methods:       make(map[string]models.MethodResult),
	}

	// Tambahin TABULAR ke list biar user bisa bandingin
	methodList := []string{"MABIMS", "WUJUDUL_HILAL", "UGHC_KHGT", "UMM_AL_QURA", "TABULAR"}

	for _, m := range methodList {
		var result models.MethodResult

		// 1. Tentukan Offset & IsTabular
		offset := 0
		isTabular := false

		switch m {
		case "TABULAR":
			offset = 0
			isTabular = true
		case "UGHC_KHGT", "UMM_AL_QURA", "TURKEY_2016":
			offset = 0
			isTabular = false
		case "MABIMS", "WUJUDUL_HILAL":
			offset = -1
			isTabular = false
		default:
			offset = -1
			isTabular = false
		}

		// 2. Setup Data Hijri Awal (Berdasarkan Offset)
		day := currentH.Day + offset
		month := currentH.Month
		year := currentH.Year

		if day > 30 {
			day -= 30
			month++
			if month > 12 {
				month = 1
				year++
			}
		} else if day < 1 {
			month--
			if month < 1 {
				month = 12
				year--
			}
			day += 30
		}

		result.HijriDate = models.HijriDate{
			Day:       day,
			Month:     month,
			MonthName: calendar.MonthNames[month-1],
			Year:      year,
			IsTabular: isTabular,
		}

		// 3. Telemetri Real-time (Detik ini)
		if !isTabular {
			alt := realtimeTel.Altitude
			elong := realtimeTel.Elongation
			result.CurrentAltitude = &alt
			result.CurrentElongation = &elong

			daysTo29 := 29 - result.HijriDate.Day
			searchDateMethod := targetDay.AddDate(0, 0, daysTo29)
			searchDateMethod = time.Date(searchDateMethod.Year(), searchDateMethod.Month(), searchDateMethod.Day(), 12, 0, 0, 0, time.UTC)

			sunset29, _ := s.Astro.GetSunset(searchDateMethod, lat, lon)
			moonset29, _ := s.Astro.GetMoonset(sunset29, lat, lon)
			tel29, _ := s.Astro.GetMoonTelemetry(sunset29, lat, lon)

			// Cari Ijtima yang berhubungan dengan akhir bulan ini
			// Kita mulai pencarian dari 5 hari sebelum sunset29 (sekitar hari ke-24)
			// Ini menjamin kita menemukan Ijtima penentu bulan tersebut.
			targetIjtima, _ := s.Cal.FindIjtima(searchDateMethod.AddDate(0, 0, -5))

			ageHours := sunset29.Sub(targetIjtima).Hours()

			// 4. Data Prediksi Akhir Bulan (Hilal Insight)
			pred := models.HilalPrediction{
				CheckDateUTC: sunset29,
				IjtimaTime:   targetIjtima,
				Altitude:     tel29.Altitude,
				Elongation:   tel29.Elongation,
				AgeHours:     ageHours,
			}

			// 5. Evaluasi Kriteria (Hanya jika bukan Tabular)
			switch m {
			case "UGHC_KHGT":
				pred.IsNewMonth = s.Cal.ScanGlobalUGHC(sunset29, targetIjtima)
			case "UMM_AL_QURA":
				pred.IsNewMonth = calendar.IsUmmAlQura(targetIjtima, sunset29, moonset29)
			default:
				resLocal := s.Cal.EvaluateLocalHisab(m, sunset29, tel29.Altitude, tel29.Elongation, sunset29, moonset29, targetIjtima)
				pred.IsNewMonth = resLocal.IsNewMonth
			}

			evaluatingMonth := result.HijriDate.Month
			evaluatingYear := result.HijriDate.Year

			// Override HANYA pada jendela selang 24 jam SETELAH sunset29
			hoursSinceSunset29 := tUTC.Sub(sunset29).Hours()
			if hoursSinceSunset29 >= 0 && hoursSinceSunset29 < 24 {
				if pred.IsNewMonth {
					nextM := (evaluatingMonth % 12) + 1
					nextY := evaluatingYear
					if nextM == 1 {
						nextY++
					}
					result.HijriDate.Day = 1
					result.HijriDate.Month = nextM
					result.HijriDate.MonthName = calendar.MonthNames[nextM-1]
					result.HijriDate.Year = nextY
				} else {
					result.HijriDate.Day = 30
					result.HijriDate.Month = evaluatingMonth
					result.HijriDate.MonthName = calendar.MonthNames[evaluatingMonth-1]
					result.HijriDate.Year = evaluatingYear
				}
				result.HijriDate.IsTabular = false
			}

			// Safe fallback jika masih lolos 31 pada hari-hari selain crossover
			if result.HijriDate.Day > 30 {
				result.HijriDate.Day -= 30
				result.HijriDate.Month++
				if result.HijriDate.Month > 12 {
					result.HijriDate.Month = 1
					result.HijriDate.Year++
				}
				result.HijriDate.MonthName = calendar.MonthNames[result.HijriDate.Month-1]
			}

			result.Prediction = &pred
		}

		resp.Methods[m] = result
	}

	return resp
}

// GetTabularOnly khusus buat pencarian tanggal tanpa telemetri astronomi
func (s *HijriService) GetTabularOnly(t time.Time) models.MethodResult {
	hDate := calendar.GetTabularHijri(t)
	hDate.IsTabular = true

	return models.MethodResult{
		HijriDate: hDate,
		Prediction: nil,
	}
}
