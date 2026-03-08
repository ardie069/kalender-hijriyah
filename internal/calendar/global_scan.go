package calendar

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
)

// KHGT Parameters
const (
	AmericasLonStart = -30.0  // Mulai dari pesisir Timur (Greenland/Brasil)
	AmericasLonEnd   = -165.0 // Berakhir di Alaska/Hawaii
	NZLat            = -41.0  // Wellington, Selandia Baru
	NZLon            = 174.0  // Garis penanggalan internasional
)

func (l *Logic) ScanGlobalUGHC(targetDateUTC time.Time, ijtimaTime time.Time) bool {
	deadline := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 23, 59, 59, 0, time.UTC)

	// CSPICE is strictly NOT thread-safe for error tracing and state!
	// Running this in goroutines causes SPICE(BADSUBSCRIPT) trace stack corruption.
	// Since SPICE calls are ~1-2 microseconds each, sequential scanning is completely fine.
	for lon := AmericasLonStart; lon >= AmericasLonEnd; lon -= 2.0 {
		for lat := 60.0; lat >= -60.0; lat -= 5.0 {
			if l.checkCoordinate(targetDateUTC, ijtimaTime, deadline, lat, lon) {
				return true
			}
		}
	}

	return false
}

// Helper function supaya kode lebih bersih dan fokus pada satu titik koordinat
func (l *Logic) checkCoordinate(targetDateUTC, ijtimaTime, deadline time.Time, lat, lon float64) bool {
	sunsetTime, err := l.Astro.GetSunset(targetDateUTC, lat, lon)

	// Handle transisi tanggal di Amerika
	if err != nil || sunsetTime.Before(ijtimaTime) {
		// Coba sunset hari berikutnya jika ijtima belum terjadi
		sunsetTime, _ = l.Astro.GetSunset(targetDateUTC.AddDate(0, 0, 1), lat, lon)
	}

	// Syarat Dasar: Sunset harus setelah Ijtima
	if sunsetTime.Before(ijtimaTime) {
		return false
	}

	et, _ := astronomy.Str2et(sunsetTime.Format(astronomy.TimeFormat))

	sunPos, _ := l.Manager.GetGeocentric(astronomy.Sun, et, "J2000")
	moonPos, _ := l.Manager.GetGeocentric(astronomy.Moon, et, "J2000")

	alt, elong := l.Astro.CalculateGeocentricParams(sunPos, moonPos, lat, lon)

	// Kriteria Turki 2016 (KHGT)
	if alt >= 5.0 && elong >= 8.0 {
		// Cek Pengecualian Amerika (Syarat B lu)
		// Jika sunset > deadline tapi di Amerika, tetap SAH selama Ijtima NZ sebelum Fajar
		if sunsetTime.After(deadline) {
			fajrNZ, _ := l.Astro.GetFajr(targetDateUTC, NZLat, NZLon)
			return ijtimaTime.Before(fajrNZ)
		}
		return true
	}

	return false
}
