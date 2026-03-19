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
	methodList := []string{"TABULAR", "MABIMS", "WUJUDUL_HILAL", "UGHC_KHGT", "UMM_AL_QURA"}

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
			}
		}
		resp.Methods[m] = result
	}

	return resp
}

func (s *HijriService) calculateMethodPrediction(m string, searchDate time.Time, lat, lon float64) (*models.HilalPrediction, error) {
	// Standard search coordinates
	checkLat, checkLon := lat, lon
	switch m {
	case "UMM_AL_QURA":
		checkLat, checkLon = 21.4225, 39.8262 // Mecca
	case "MABIMS":
		checkLat, checkLon = 5.89, 95.32 // Sabang (Westernmost Indonesia)
	}

	sunset, err := s.Astro.GetSunset(searchDate, checkLat, checkLon)
	if err != nil {
		return nil, err
	}

	moonset, err := s.Astro.GetMoonset(sunset, checkLat, checkLon)
	if err != nil {
		return nil, err
	}

	tel, err := s.Astro.GetMoonTelemetry(sunset, checkLat, checkLon)
	if err != nil {
		return nil, err
	}

	ijtima, err := s.Cal.FindIjtima(searchDate.AddDate(0, 0, -5))
	if err != nil {
		ijtima = searchDate // Fallback (should not happen with FindIjtima)
	}

	localTime, tzName := s.GetLocalTimeInfo(sunset, checkLat, checkLon)

	pred := &models.HilalPrediction{
		CheckDateUTC:      sunset,
		CheckDateLocal:    localTime,
		TimezoneName:      tzName,
		IjtimaTime:        ijtima,
		AltitudeGeometric: tel.Altitude,
		AltitudeApparent:  tel.Altitude + astronomy.ApplyRefraction(tel.Altitude),
		Elongation:        tel.Elongation,
		AgeHours:          sunset.Sub(ijtima).Hours(),
		Location:          &models.LocationInfo{Lat: checkLat, Lon: checkLon},
	}

	// Method-specific evaluation
	switch m {
	case "UGHC_KHGT":
		isNew, globalPred := s.Cal.ScanGlobalUGHC(sunset, ijtima)
		if globalPred != nil {
			loc := globalPred.Location
			localTimeUGHC, tzNameUGHC := s.GetLocalTimeInfo(globalPred.CheckDateUTC, loc.Lat, loc.Lon)
			*pred = *globalPred
			pred.CheckDateLocal = localTimeUGHC
			pred.TimezoneName = tzNameUGHC
		}
		pred.IsNewMonth = isNew
	case "UMM_AL_QURA":
		pred.IsNewMonth = calendar.IsUmmAlQura(ijtima, sunset, moonset)
	default:
		resLocal := s.Cal.EvaluateLocalHisab(m, sunset, tel.Altitude, tel.Elongation, sunset, moonset, ijtima)
		pred.IsNewMonth = resLocal.IsNewMonth
	}

	return pred, nil
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

func (s *HijriService) GetLocation(lat, lon float64) *time.Location {
	tzName := s.finder.GetTimezoneName(lon, lat)
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.UTC
	}
	return loc
}

// GetTabularOnly khusus buat pencarian tanggal tanpa telemetri astronomi
func (s *HijriService) GetTabularOnly(t time.Time) models.MethodResult {
	hDate := calendar.GetTabularHijri(t)
	hDate.IsTabular = true

	return models.MethodResult{
		HijriDate:  hDate,
		Prediction: nil,
	}
}

// GetHijriTargetDate: Returns the Gregorian date that corresponds to the current Hijri day.
// It handles rollover at sunset (or Isha for critical months on the 29th).
func (s *HijriService) GetHijriTargetDate(t time.Time, lat, lon float64) time.Time {
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
