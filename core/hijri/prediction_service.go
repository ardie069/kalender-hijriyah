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

	pred := &models.HilalPrediction{
		CheckDateUTC:     sunset,
		CheckDateLocal:   localTime,
		TimezoneName:     tzName,
		IjtimaTime:       ijtima,
		Altitude:         tel.Altitude, // Default is topocentric
		AltitudeApparent: tel.AltitudeApparent,
		Elongation:       tel.Elongation,
		AgeHours:         sunset.Sub(ijtima).Hours(),
		Location:         &models.LocationInfo{Lat: checkLat, Lon: checkLon},
	}

	// Method-specific evaluation
	switch m {
	case "WUJUDUL_HILAL":
		// Wujudul Hilal requires geocentric altitude
		altGeo, _ := s.Astro.CalculateGeocentricParamsGlobal(sunset, checkLat, checkLon)
		pred.Altitude = altGeo
		pred.AltitudeApparent = nil // Exclude for geocentric-based methods
		pred.IsNewMonth = calendar.IsWujudulHilal(altGeo, sunset.After(ijtima))
	case "KHGT":
		isNew, globalPred := s.Cal.ScanGlobalKHGT(sunset, ijtima)
		if globalPred != nil {
			loc := globalPred.Location
			localTimeUGHC, tzNameUGHC := s.Tz.GetLocalTimeInfo(globalPred.CheckDateUTC, loc.Lat, loc.Lon)
			*pred = *globalPred
			pred.CheckDateLocal = localTimeUGHC
			pred.TimezoneName = tzNameUGHC
		}
		pred.AltitudeApparent = nil // CRITICAL: Exclude for KHGT
		pred.IsNewMonth = isNew
	case "UMM_AL_QURA":
		pred.IsNewMonth = calendar.IsUmmAlQura(ijtima, sunset, moonset)
	default:
		resLocal := s.Cal.EvaluateLocalHisab(m, sunset, tel, sunset, moonset, ijtima)
		pred.IsNewMonth = resLocal.IsNewMonth
	}

	return pred, nil
}
