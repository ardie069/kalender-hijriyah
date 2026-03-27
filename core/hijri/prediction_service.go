package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// calculateMethodPrediction computes hilal prediction for a specific method at the given search date.
func (s *DateService) calculateMethodPrediction(m string, searchDate time.Time, lat, lon float64) (*models.HilalPrediction, error) {
	// Standard search coordinates
	checkLat, checkLon := lat, lon
	switch m {
	case "UMM_AL_QURA":
		checkLat, checkLon = 21.4225, 39.8262 // Mecca
	case "MABIMS":
		checkLat, checkLon = 5.89, 95.32 // Sabang (Westernmost Indonesia)
	case "MABIMS_LOCAL":
		// checkLat, checkLon stays as lat, lon (user's location)
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

	localTime, tzName := s.Tz.GetLocalTimeInfo(sunset, checkLat, checkLon)
	localMoonset, _ := s.Tz.GetLocalTimeInfo(moonset, checkLat, checkLon)

	pred := &models.HilalPrediction{
		CheckDateUTC:       sunset,
		CheckDateLocal:     localTime,
		TimezoneName:       tzName,
		IjtimaTime:         ijtima,
		Altitude:           tel.Altitude, // Default is topocentric
		AltitudeApparent:   tel.AltitudeApparent,
		Elongation:         tel.Elongation,
		AgeHours:           sunset.Sub(ijtima).Hours(),
		Location:           &models.LocationInfo{Lat: checkLat, Lon: checkLon},
		MoonsetTimeLocal:   localMoonset,
		MoonsetDiffMinutes: moonset.Sub(sunset).Minutes(),
	}

	// Method-specific evaluation
	switch m {
	case "TABULAR":
		// Hijri month length in Tabular: 1,3,5,7,9,11 are 30 days.
		// searchDate is day 29 of the CURRENT month.
		currentH := calendar.GetTabularHijri(searchDate)
		monthLen := 29
		if currentH.Month%2 != 0 {
			monthLen = 30
		} else if currentH.Month == 12 {
			// Leap year check: (11*Year + 14) % 30 < 11 are leap years (30 days)
			if (11*currentH.Year+14)%30 < 11 {
				monthLen = 30
			}
		}
		pred.IsNewMonth = (monthLen == 29)
		pred.CheckDateUTC = searchDate // Stable reference for tabular
	case "WUJUDUL_HILAL":
		// Wujudul Hilal requires geocentric altitude
		altGeo, _ := s.Astro.CalculateGeocentricParamsGlobal(sunset, checkLat, checkLon)
		pred.Altitude = altGeo
		pred.AltitudeApparent = nil // Exclude for geocentric-based methods
		pred.IsNewMonth = calendar.IsWujudulHilal(altGeo, sunset.After(ijtima))
	case "KHGT":
		khgtResult := s.Cal.ScanGlobalKHGT(sunset, ijtima)
		if khgtResult.BestVisibility != nil {
			loc := khgtResult.BestVisibility.Location
			localTimeUGHC, tzNameUGHC := s.Tz.GetLocalTimeInfo(khgtResult.BestVisibility.CheckDateUTC, loc.Lat, loc.Lon)
			*pred = *khgtResult.BestVisibility
			pred.CheckDateLocal = localTimeUGHC
			pred.TimezoneName = tzNameUGHC
			
			// Inject explicit flags 
			pred.KHGTGlobalValid = &khgtResult.IsGlobalValid
			pred.KHGTAmericaException = &khgtResult.IsAmericaException

			// CRITICAL: Stable CheckDateUTC for KHGT projection.
			// khgtResult.Date is the beginning of the scanned UTC day.
			// We use this to avoid Hawaii's sunset (crossing the date line) shifting the projection by +1 day.
			pred.CheckDateUTC = khgtResult.Date
		}
		pred.AltitudeApparent = nil // CRITICAL: Exclude for KHGT
		pred.IsNewMonth = khgtResult.IsGlobalValid
	case "UMM_AL_QURA":
		pred.IsNewMonth = calendar.IsUmmAlQura(ijtima, sunset, moonset)
	case "MABIMS_LOCAL":
		resLocal := s.Cal.EvaluateLocalHisab("MABIMS", sunset, tel, sunset, moonset, ijtima)
		pred.IsNewMonth = resLocal.IsNewMonth
	default:
		resLocal := s.Cal.EvaluateLocalHisab(m, sunset, tel, sunset, moonset, ijtima)
		pred.IsNewMonth = resLocal.IsNewMonth
	}

	return pred, nil
}
