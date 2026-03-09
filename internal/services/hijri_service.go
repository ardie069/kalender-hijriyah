package services

import (
	"fmt"
	"time"
	"github.com/ringsaturn/tzf"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

type HijriService struct {
	finder tzf.F
	Astro  *astronomy.Adapter
	Cal    *calendar.Logic
}

func NewHijriService(astro *astronomy.Adapter, cal *calendar.Logic) (*HijriService, error) {
	f, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, err
	}
	return &HijriService{
		finder: f,
		Astro:  astro,
		Cal:    cal,
	}, nil
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

		// 3. Telemetri Real-time (Detik ini)
		if !result.HijriDate.IsTabular {
			alt := realtimeTel.Altitude
			elong := realtimeTel.Elongation
			result.CurrentAltitude = &alt
			result.CurrentElongation = &elong

			// Jika sudah direview dan ternyata kita lagi di hari ke 30, abaikan prediksi karena Hilal fix nggak dicek lagi di hari ke-30
			// (Bulan otomatis jadi 1 besoknya). 
			// Kita hanya memunculkan hasil prediksi pada hari ke-29.
			if result.HijriDate.Day <= 29 {
				daysTo29 := 29 - result.HijriDate.Day
				searchDateMethod := targetDay.AddDate(0, 0, daysTo29)
				searchDateMethod = time.Date(searchDateMethod.Year(), searchDateMethod.Month(), searchDateMethod.Day(), 12, 0, 0, 0, time.UTC)

				sunset29, _ := s.Astro.GetSunset(searchDateMethod, lat, lon)
				moonset29, _ := s.Astro.GetMoonset(sunset29, lat, lon)
				tel29, _ := s.Astro.GetMoonTelemetry(sunset29, lat, lon)

				// Cari Ijtima yang berhubungan dengan akhir bulan ini
				targetIjtima, _ := s.Cal.FindIjtima(searchDateMethod.AddDate(0, 0, -5))

				ageHours := sunset29.Sub(targetIjtima).Hours()

				localTime, tzName := s.GetLocalTimeInfo(sunset29, lat, lon)

				// 4. Data Prediksi Akhir Bulan (Hilal Insight)
				pred := models.HilalPrediction{
					CheckDateUTC:   sunset29,
					CheckDateLocal: localTime,
					TimezoneName:   tzName,
					IjtimaTime:     targetIjtima,
					Altitude:       tel29.Altitude,
					Elongation:     tel29.Elongation,
					AgeHours:       ageHours,
				}

				// 5. Evaluasi Kriteria (Hanya jika bukan Tabular)
				switch m {
				case "UGHC_KHGT":
					isNew, globalPred := s.Cal.ScanGlobalUGHC(sunset29, targetIjtima)
					if globalPred != nil {
						loc := globalPred.Location
						localTimeUGHC, tzNameUGHC := s.GetLocalTimeInfo(globalPred.CheckDateUTC, loc.Lat, loc.Lon)
						
						pred = *globalPred
						pred.CheckDateLocal = localTimeUGHC
						pred.TimezoneName = tzNameUGHC
					}
					pred.IsNewMonth = isNew
				case "UMM_AL_QURA":
					meccaLat, meccaLon := 21.4225, 39.8262
					sunsetMecca, _ := s.Astro.GetSunset(searchDateMethod, meccaLat, meccaLon)
					moonsetMecca, _ := s.Astro.GetMoonset(sunsetMecca, meccaLat, meccaLon)
					telMecca, _ := s.Astro.GetMoonTelemetry(sunsetMecca, meccaLat, meccaLon)

					localMecca, tzMecca := s.GetLocalTimeInfo(sunsetMecca, meccaLat, meccaLon)
					pred.CheckDateUTC = sunsetMecca
					pred.CheckDateLocal = localMecca
					pred.TimezoneName = tzMecca
					pred.Altitude = telMecca.Altitude
					pred.Elongation = telMecca.Elongation
					pred.AgeHours = sunsetMecca.Sub(targetIjtima).Hours()
					pred.Location = &models.LocationInfo{Lat: meccaLat, Lon: meccaLon}
					pred.IsNewMonth = calendar.IsUmmAlQura(targetIjtima, sunsetMecca, moonsetMecca)
				default:
					resLocal := s.Cal.EvaluateLocalHisab(m, sunset29, tel29.Altitude, tel29.Elongation, sunset29, moonset29, targetIjtima)
					pred.Location = &models.LocationInfo{Lat: lat, Lon: lon}
					pred.IsNewMonth = resLocal.IsNewMonth
				}

				result.Prediction = &pred
			}
		}

		resp.Methods[m] = result
	}

	return resp
}

func (s *HijriService) GetLocalTimeInfo(utc time.Time, lat, lon float64) (string, string) {
	// 1. Cari nama timezone (misal: "America/Anchorage")
	tzName := s.finder.GetTimezoneName(lon, lat)
	
	// 2. Load lokasi berdasarkan database IANA (Go udah punya ini di time package)
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return utc.Format("2006-01-02 15:04:05"), "UTC"
	}

	// 3. Konversi UTC ke Waktu Lokal (Otomatis handle DST!)
	local := utc.In(loc)
	
	// Contoh return: "2026-03-18 21:06:57", "AKDT (UTC-8)"
	_, offset := local.Zone()
	return local.Format("2006-01-02 15:04:05"), fmt.Sprintf("%s (UTC %d)", local.Format("MST"), offset/3600)
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
