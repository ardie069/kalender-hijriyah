package calendar

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

// KHGT Parameters
const (
	AmericasLonStart = -30.0  // Mulai dari pesisir Timur (Greenland/Brasil)
	AmericasLonEnd   = -165.0 // Berakhir di Alaska/Hawaii
	NZLat            = -41.0  // Wellington, Selandia Baru
	NZLon            = 174.0  // Garis penanggalan internasional
)

func (l *Logic) ScanGlobalUGHC(targetDateUTC time.Time, ijtimaTime time.Time) (bool, *models.HilalPrediction) {
	deadline := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 23, 59, 59, 0, time.UTC)

	var bestPred *models.HilalPrediction
	var maxAlt float64 = -999.0

	// CSPICE is strictly NOT thread-safe for error tracing and state!
	// Running this in goroutines causes SPICE(BADSUBSCRIPT) trace stack corruption.
	// Since SPICE calls are ~1-2 microseconds each, sequential scanning is completely fine.
	for lon := AmericasLonStart; lon >= AmericasLonEnd; lon -= 2.0 {
		for lat := 60.0; lat >= -60.0; lat -= 5.0 {
			found, pred := l.checkCoordinate(targetDateUTC, ijtimaTime, deadline, lat, lon)
			if pred != nil && pred.AltitudeGeometric > maxAlt {
				maxAlt = pred.AltitudeGeometric
				bestPred = pred
			}
			if found {
				return true, pred
			}
		}
	}

	return false, bestPred
}

// Helper function supaya kode lebih bersih dan fokus pada satu titik koordinat
func (l *Logic) checkCoordinate(targetDateUTC, ijtimaTime, deadline time.Time, lat, lon float64) (bool, *models.HilalPrediction) {
	sunsetTime, err := l.Astro.GetSunset(targetDateUTC, lat, lon)

	// Syarat Dasar: Sunset harus setelah Ijtima
	// Jika gagal, berarti di lokasi itu hilal mustahil terlihat hari ini. Tidak boleh lompat ke besok.
	if err != nil || sunsetTime.Before(ijtimaTime) {
		return false, nil
	}

	et := astronomy.TimeToEt(sunsetTime)

	alt, elong := l.Astro.CalculateGeocentricParams(et, lat, lon)

	ageHours := sunsetTime.Sub(ijtimaTime).Hours()
	pred := &models.HilalPrediction{
		CheckDateUTC:      sunsetTime,
		IsNewMonth:        false,
		IjtimaTime:        ijtimaTime,
		AltitudeGeometric: alt,
		AltitudeApparent:  alt + astronomy.ApplyRefraction(alt),
		Elongation:        elong,
		AgeHours:          ageHours,
		Location:          &models.LocationInfo{Lat: lat, Lon: lon},
	}

	// Kriteria Turki 2016 (KHGT)
	if alt >= 5.0 && elong >= 8.0 {
		// Cek Pengecualian Amerika (Syarat B lu)
		// Jika sunset > deadline tapi di Amerika, tetap SAH selama Ijtima NZ sebelum Fajar
		if sunsetTime.After(deadline) {
			fajrNZ, err := l.Astro.GetFajr(targetDateUTC.AddDate(0, 0, 1), NZLat, NZLon)
			if err == nil && ijtimaTime.Before(fajrNZ) {
				pred.IsNewMonth = true
				return true, pred
			}
			return false, pred
		}
		pred.IsNewMonth = true
		return true, pred
	}

	return false, pred
}
