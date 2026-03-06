package astronomy

import (
	"math"
	"time"
)

type Adapter struct {
	Manager *EphemerisManager
}

// findRoot: Fungsi Bisection generik buat nyari waktu kejadian
func (a *Adapter) findRoot(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) float64 {
	low, high := startET, startET+duration
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, _ := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)
		currentAlt := currentAltDeg * (math.Pi / 180.0)

		// Logic spesifik buat Moonset/Sunset (Alt harus turun)
		if currentAlt > targetAlt {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (a *Adapter) CalculateGeocentricParams(sunPos, moonPos Vector3, lat, lon float64) (altitude, elongation float64) {
	// 1. Elongation: Sudut pisah matahari dan bulan di pusat bumi
	elongation = math.Acos(sunPos.Unit().Dot(moonPos.Unit())) * (180.0 / math.Pi)

	// 2. Geocentric Altitude at specific location.
	// Kita ambil Vektor posisional GEO (tapi dengan referensi horizon lokal si observer)
	// Memakai GetLocalAltAz yang proyeksikan relatif ke Zenith lintang bujur observasi.
	altitude, _ = a.Manager.GetLocalAltAz(moonPos, lat, lon)

	return altitude, elongation
}

// GetFajr: Mencari waktu subuh (Solar Altitude = -18 derajat)
func (a *Adapter) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	startET, _ := Str2et(date.Format("2006-01-02") + " 00:00:00 UTC")
	targetAlt := -18.0 * (math.Pi / 180.0)

	et := a.findRootUp(startET, 86400, targetAlt, Sun, lat, lon)
	return Et2Utc(et), nil
}

// findRootUp: Bisection naik khusus untuk Sunrise/Fajr
func (a *Adapter) findRootUp(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) float64 {
	// 1. Cari jam kasar kapan Altitude naik melewati target
	var low, high float64
	foundInterval := false

	for offset := 0.0; offset < duration; offset += 3600.0 {
		pos1, _ := a.Manager.GetTopocentricPosition(target, startET+offset, lat, lon)
		alt1, _ := a.Manager.GetLocalAltAz(pos1, lat, lon)
		pos2, _ := a.Manager.GetTopocentricPosition(target, startET+offset+3600.0, lat, lon)
		alt2, _ := a.Manager.GetLocalAltAz(pos2, lat, lon)

		radTarget := targetAlt * (180.0 / math.Pi)
		// Fajr is when alt1 < target && alt2 > target
		if alt1 < radTarget && alt2 >= radTarget {
			low = startET + offset
			high = startET + offset + 3600.0
			foundInterval = true
			break
		}
	}

	if !foundInterval {
		return startET + duration/2 // Fallback
	}

	// 2. Bisection Halus
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, _ := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)
		currentAlt := currentAltDeg * (math.Pi / 180.0)

		if currentAlt < targetAlt {
			low = mid // Masih di bawah target, root di kanan
		} else {
			high = mid // Udah lewat target, root di kiri
		}
	}
	return low
}

func (a *Adapter) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	startET, _ := Str2et(date.Format("2006-01-02") + " 12:00:00 UTC")
	et := a.findRoot(startET, 86400, -0.583*(math.Pi/180.0), Moon, lat, lon)
	return Et2Utc(et), nil
}

func (a *Adapter) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	// Start at 00:00 UTC of that day
	startET, _ := Str2et(date.Format("2006-01-02") + " 00:00:00 UTC")
	// Search over the whole 24 hours
	et := a.findRootDown(startET, 86400, -0.833*(math.Pi/180.0), Sun, lat, lon)
	return Et2Utc(et), nil
}

// findRootDown: Bisection turun khusus untuk Sunset/Moonset
func (a *Adapter) findRootDown(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) float64 {
	// 1. Cari jam kasar kapan Altitude turun melewati target (Grid search 1 jam)
	var low, high float64
	foundInterval := false

	for offset := 0.0; offset < duration; offset += 3600.0 {
		pos1, _ := a.Manager.GetTopocentricPosition(target, startET+offset, lat, lon)
		alt1, _ := a.Manager.GetLocalAltAz(pos1, lat, lon)
		pos2, _ := a.Manager.GetTopocentricPosition(target, startET+offset+3600.0, lat, lon)
		alt2, _ := a.Manager.GetLocalAltAz(pos2, lat, lon)

		// Sunset is when alt1 > target && alt2 < target
		if alt1 > targetAlt*(180.0/math.Pi) && alt2 <= targetAlt*(180.0/math.Pi) {
			low = startET + offset
			high = startET + offset + 3600.0
			foundInterval = true
			break
		}
	}

	if !foundInterval {
		return startET + duration/2 // Fallback
	}

	// 2. Bisection Halus (Precision)
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, _ := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)
		currentAlt := currentAltDeg * (math.Pi / 180.0)

		if currentAlt > targetAlt {
			low = mid // Masih di atas ufuk, root ada di kanan
		} else {
			high = mid // Udah tenggelam, root ada di kiri
		}
	}
	return low
}
