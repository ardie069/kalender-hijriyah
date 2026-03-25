package calendar

import (
	"github.com/ardie069/kalender-hijriyah/internal/models"
	"time"
)

// KHGT Parameters berdasarkan Muhammadiyah
const (
	// Global scan coverage
	GlobalLonStart = -180.0
	GlobalLonEnd   = 180.0
	GlobalLatStart = 65.0
	GlobalLatEnd   = -65.0

	// Amerika continent bounds (daratan perkiraan)
	AmericaLonStart = -168.0
	AmericaLonEnd   = -30.0
	AmericaLatStart = 60.0
	AmericaLatEnd   = -55.0

	// NZ untuk pengecualian
	NZLat = -41.2866 // Wellington
	NZLon = 174.7756

	// Kriteria KHGT (Turki 2016)
	MinElongation = 8.0
	MinAltitude   = 5.0
)

func (l *Logic) ScanGlobalKHGT(targetDateUTC time.Time, ijtimaTimeUTC time.Time) (bool, *models.HilalPrediction) {
	deadlineUTC := time.Date(
		targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(),
		23, 59, 59, 999999999, time.UTC,
	)

	result := &models.KHGTResult{
		Date:               targetDateUTC,
		IsGlobalValid:      false,
		IsAmericaException: false,
		ValidLocations:     []models.HilalPrediction{},
	}

	var bestPred *models.HilalPrediction
	maxScore := -999.0

	// Scan global grid
Loop:
	for lon := GlobalLonStart; lon <= GlobalLonEnd; lon += 5.0 {
		for lat := GlobalLatStart; lat >= GlobalLatEnd; lat -= 5.0 {
			// Evaluasi titik ini dengan KRITERIA GEOSENTRIS
			isValid, pred := l.evaluateGeocentricPoint(
				targetDateUTC, ijtimaTimeUTC, deadlineUTC,
				lat, lon, &result.IsAmericaException,
			)

			if pred != nil {
				score := pred.Altitude + pred.Elongation
				if score > maxScore {
					maxScore = score
					bestPred = pred
				}
			}

			if isValid {
				result.ValidLocations = append(result.ValidLocations, *pred)
				result.IsGlobalValid = true

				// Early exit jika sudah ketemu dan bukan Amerika exception
				if !result.IsAmericaException {
					break Loop
				}
			}
		}
	}
	return result.IsGlobalValid, bestPred
}

func (l *Logic) evaluateGeocentricPoint(
	targetDateUTC, ijtimaTimeUTC, deadlineUTC time.Time,
	lat, lon float64,
	americaExceptionFlag *bool,
) (bool, *models.HilalPrediction) {

	// Hitung sunset di lokasi ini
	sunsetUTC, err := l.Astro.GetSunset(targetDateUTC, lat, lon)
	if err != nil {
		return false, nil
	}

	// Syarat dasar: sunset harus setelah ijtima
	if sunsetUTC.Before(ijtimaTimeUTC) {
		return false, nil
	}

	// KRITIS: Hitung altitude dan elongasi GEOSENTRIS
	// (dari pusat bumi, tanpa koreksi lat/lon)
	altGeo, elongGeo := l.Astro.CalculateGeocentricParamsGlobal(sunsetUTC, lat, lon)

	ageHours := sunsetUTC.Sub(ijtimaTimeUTC).Hours()
	pred := &models.HilalPrediction{
		CheckDateUTC:      sunsetUTC,
		IjtimaTime:        ijtimaTimeUTC,
		Altitude:          altGeo,
		Elongation:        elongGeo,
		AgeHours:          ageHours,
		Location:          &models.LocationInfo{Lat: lat, Lon: lon},
		IsNewMonth:        false,
	}

	// Evaluasi kriteria utama (Turki 2016)
	if altGeo >= MinAltitude && elongGeo >= MinElongation {
		// Kasus 1: Terpenuhi SEBELUM pukul 24.00 UTC di bagian mana pun di dunia
		if sunsetUTC.Before(deadlineUTC) || sunsetUTC.Equal(deadlineUTC) {
			pred.IsNewMonth = true
			return true, pred
		}

		// Kasus 2: Terpenuhi SETELAH pukul 24.00 UTC
		// Bulan baru dimulai JIKA:
		// 1. Terjadi di wilayah daratan Benua Amerika
		// 2. Ijtimak terjadi sebelum fajar di Selandia Baru
		isAmericaLand := (lon >= AmericaLonStart && lon <= AmericaLonEnd &&
			lat >= AmericaLatEnd && lat <= AmericaLatStart)

		if isAmericaLand {
			// Cek ijtima sebelum fajar NZ (Wellington)
			// Syarat: h+1 fajar NZ
			nextDay := targetDateUTC.AddDate(0, 0, 1)
			fajrNZ, err := l.Astro.GetFajr(nextDay, NZLat, NZLon)
			if err == nil && ijtimaTimeUTC.Before(fajrNZ) {
				pred.IsNewMonth = true
				*americaExceptionFlag = true
				return true, pred
			}
		}

		// Terpenuhi setelah 24:00 UTC tapi tidak memenuhi syarat exception → TIDAK SAH
		return false, pred
	}

	return false, pred
}
