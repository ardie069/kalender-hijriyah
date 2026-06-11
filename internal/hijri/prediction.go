package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/decision"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

// PredictHijriDate: Menghitung transisi tanggal
func PredictHijriDate(t time.Time, isNewMonth bool, currentMonth, currentYear int) models.HijriDate {
	if isNewMonth {
		// Transisi ke bulan baru (Tanggal 1)
		newMonth := (currentMonth % 12) + 1
		newYear := currentYear
		if newMonth == 1 {
			newYear++
		}
		return models.HijriDate{
			Day:       1,
			Month:     newMonth,
			MonthName: calendar.MonthNames[newMonth-1],
			Year:      newYear,
		}
	}

	// Gagal ganti bulan -> Bulan sebelumnya digenapkan (Istikmal 30 hari)
	return models.HijriDate{
		Day:       30,
		Month:     currentMonth,
		MonthName: calendar.MonthNames[currentMonth-1],
		Year:      currentYear,
	}
}

func (s *DateService) CalculateMethodPrediction(m string, searchDate time.Time, lat, lon float64) (*models.HilalPrediction, error) {
	checkLat, checkLon := lat, lon
	switch m {
	case "UMM_AL_QURA":
		checkLat, checkLon = 21.4225, 39.8262
	case "SAUDI":
		checkLat, checkLon = 25.5950, 45.6339 // Sudair (lat: 25.5950, lon: 45.6339)
	case "MABIMS":
		checkLat, checkLon = 5.89, 95.32
	case "MABIMS_LOCAL":
		// checkLat, checkLon stays
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

	ijtima, err := s.Ephem.FindIjtima(searchDate)
	if err != nil {
		ijtima = searchDate
	}

	localTime, tzName := s.Tz.GetLocalTimeInfo(sunset, checkLat, checkLon)
	localMoonset, _ := s.Tz.GetLocalTimeInfo(moonset, checkLat, checkLon)

	pred := &models.HilalPrediction{
		CheckDateUTC:       sunset,
		CheckDateLocal:     localTime,
		TimezoneName:       tzName,
		IjtimaTime:         ijtima,
		Altitude:           tel.Altitude,
		AltitudeApparent:   tel.AltitudeApparent,
		Elongation:         tel.Elongation,
		ElongationGeo:      tel.ElongationGeo,
		AgeHours:           sunset.Sub(ijtima).Hours(),
		Location:           &models.LocationInfo{Lat: checkLat, Lon: checkLon},
		MoonsetTimeLocal:   localMoonset,
		MoonsetDiffMinutes: moonset.Sub(sunset).Minutes(),
	}

	switch m {
	case "TABULAR":
		currentH := calendar.GetTabularHijri(searchDate)
		monthLen := 29
		if currentH.Month%2 != 0 {
			monthLen = 30
		} else if currentH.Month == 12 {
			if (11*currentH.Year+14)%30 < 11 {
				monthLen = 30
			}
		}
		pred.IsNewMonth = (monthLen == 29)
		pred.CheckDateUTC = searchDate

	case "KHGT":
		khgtResult := s.Scan.ScanGlobalKHGT(sunset, ijtima)
		if khgtResult.BestVisibility != nil {
			loc := khgtResult.BestVisibility.Location
			localTimeUGHC, tzNameUGHC := s.Tz.GetLocalTimeInfo(khgtResult.BestVisibility.CheckDateUTC, loc.Lat, loc.Lon)
			*pred = *khgtResult.BestVisibility
			pred.CheckDateLocal = localTimeUGHC
			pred.TimezoneName = tzNameUGHC
			
			pred.KHGTGlobalValid = &khgtResult.IsGlobalValid
			pred.KHGTAmericaException = &khgtResult.IsAmericaException
			pred.CheckDateUTC = khgtResult.Date
		}
		pred.AltitudeApparent = nil
		pred.IsNewMonth = khgtResult.IsGlobalValid
		
	case "UMM_AL_QURA":
		hisabMecca := s.UmmAlQura.EvaluateHisabMecca(sunset, ijtima)
		
		rukyatSaudi := s.evaluateSaudiRukyat(sunset, ijtima)
		
		hDate := calendar.GetTabularHijri(ijtima.AddDate(0, 0, 15))
		isCriticalMonth := (hDate.Month == 9 || hDate.Month == 10 || hDate.Month == 12)
		
		if isCriticalMonth {
			pred.IsNewMonth = rukyatSaudi
		} else {
			pred.IsNewMonth = hisabMecca
			if hisabMecca != rukyatSaudi {
				pred.IsNewMonth = rukyatSaudi
			}
		}
		
	case "SAUDI":
		saudi := decision.Saudi{}
		ctxSudair := decision.Context{
			IjtimaTime:  ijtima,
			SunsetTime:  sunset,
			MoonsetTime: moonset,
		}
		isNewMonth := saudi.IsVisible(ctxSudair)
		if !isNewMonth {
			sunsetTumair, _ := s.Astro.GetSunset(searchDate, 25.7039, 45.8614)
			moonsetTumair, _ := s.Astro.GetMoonset(sunsetTumair, 25.7039, 45.8614)
			ctxTumair := decision.Context{
				IjtimaTime:  ijtima,
				SunsetTime:  sunsetTumair,
				MoonsetTime: moonsetTumair,
			}
			isNewMonth = saudi.IsVisible(ctxTumair)
		}
		pred.IsNewMonth = isNewMonth
		
	case "MABIMS", "MABIMS_LOCAL":
		mabims := decision.MABIMS{}
		ctx := decision.Context{
			Altitude:      tel.Altitude,
			Elongation:    tel.Elongation,
			ElongationGeo: tel.ElongationGeo,
		}
		pred.IsNewMonth = mabims.IsVisible(ctx)
		
	default:
		pred.IsNewMonth = false
	}

	return pred, nil
}

