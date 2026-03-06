package calendar

import (
	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"time"
)

// KHGT Parameters
const (
	AmericasLonStart = -30.0  // Mulai dari pesisir Timur (Greenland/Brasil)
	AmericasLonEnd   = -165.0 // Berakhir di Alaska/Hawaii
	NZLat            = -41.0  // Wellington, Selandia Baru
	NZLon            = 174.0  // Garis penanggalan internasional
)

func (l *Logic) ScanGlobalUGHC(targetDateUTC time.Time, ijtimaTime time.Time) bool {
	// 1. DEADLINE: 24:00 UTC
	deadline := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 23, 59, 59, 0, time.UTC)

	// 2. SCAN BENUA AMERIKA (Daratan)
	// Kita scan wilayah yang matahari terbenamnya paling mendekati deadline UTC
	foundInAmericas := false
	for lon := AmericasLonStart; lon >= AmericasLonEnd; lon -= 1.0 {
		// Scan di berbagai latitude Amerika (Utara ke Selatan)
		for lat := 60.0; lat >= -60.0; lat -= 10.0 {
			sunsetTime, err := l.Astro.GetSunset(targetDateUTC, lat, lon)
			if err != nil || sunsetTime.After(deadline) {
				continue
			}

			et, _ := astronomy.Str2et(sunsetTime.Format(astronomy.TimeFormat))
			sunPos, _ := l.Manager.GetGeocentric(astronomy.Sun, et, astronomy.FrameJ2000)
			moonPos, _ := l.Manager.GetGeocentric(astronomy.Moon, et, astronomy.FrameJ2000)

			alt, elong := l.Astro.CalculateGeocentricParams(sunPos, moonPos)

			if IsTurkey2016(alt, elong) {
				foundInAmericas = true
				break
			}
		}
		if foundInAmericas {
			break
		}
	}

	// 3. SELANDIA BARU OVERRIDE (Butir 4)
	// Jika kriteria terpenuhi di Amerika, cek apakah ijtima di NZ terjadi sebelum Subuh
	if foundInAmericas {
		fajrNZ, _ := l.Astro.GetFajr(targetDateUTC, NZLat, NZLon)
		// Jika Ijtima terjadi sebelum subuh di wilayah UTC+13
		if ijtimaTime.Before(fajrNZ) {
			return true
		}
	}

	return false
}
