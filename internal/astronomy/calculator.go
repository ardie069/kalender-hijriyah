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

func (a *Adapter) CalculateGeocentricParams(et float64, lat, lon float64) (altitude, elongation float64) {
	// Ambil posisi J2000 Inertial untuk mengukur sudut Elongasi absolut antar benda langit
	sunPos, _ := a.Manager.GetGeocentric(Sun, et, "J2000")
	moonPos, _ := a.Manager.GetGeocentric(Moon, et, "J2000")

	// 1. Elongation: Sudut pisah matahari dan bulan di pusat bumi
	elongation = math.Acos(sunPos.Unit().Dot(moonPos.Unit())) * (180.0 / math.Pi)

	// 2. Topocentric Altitude at specific location.
	// Wajib mengambil vektor berdasar perputaran bumi rotasional (IAU_EARTH) ke titik topocenter,
	// TIDAK BOLEH memasukkan J2000 yang bersifat inersial ke fungsi proyeksi AltAz lokal!
	topoPos, _ := a.Manager.GetTopocentricPosition(Moon, et, lat, lon)
	altitude, _ = a.Manager.GetLocalAltAz(topoPos, lat, lon)

	return altitude, elongation
}

// GetFajr: Mencari waktu subuh (Solar Altitude = -18 derajat)
func (a *Adapter) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	// Ambil 00:00 LOKAL hari itu (dengan menggeser UTC start berdasar bujur)
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon/15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	targetAlt := -18.0 * (math.Pi / 180.0)

	et := a.findRootUp(startET, 86400, targetAlt, Sun, lat, lon)
	return Et2Utc(et), nil
}

// findRootUp: Bisection naik khusus untuk Sunrise/Fajr
func (a *Adapter) findRootUp(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) float64 {
	// targetAlt is in radians, convert to degrees for consistent comparison with GetLocalAltAz output
	targetAltDeg := targetAlt * (180.0 / math.Pi)

	// 1. Cari jam kasar kapan Altitude naik melewati target
	var low, high float64
	foundInterval := false

	for offset := 0.0; offset < duration; offset += 3600.0 {
		pos1, _ := a.Manager.GetTopocentricPosition(target, startET+offset, lat, lon)
		alt1, _ := a.Manager.GetLocalAltAz(pos1, lat, lon)
		pos2, _ := a.Manager.GetTopocentricPosition(target, startET+offset+3600.0, lat, lon)
		alt2, _ := a.Manager.GetLocalAltAz(pos2, lat, lon)

		if alt1 < targetAltDeg && alt2 >= targetAltDeg {
			low = startET + offset
			high = startET + offset + 3600.0
			foundInterval = true
			break
		}
	}

	if !foundInterval {
		return startET + duration/2 // Fallback
	}

	// 2. Bisection Halus (degrees throughout)
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, _ := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)

		if currentAltDeg < targetAltDeg {
			low = mid // Masih di bawah target, root di kanan
		} else {
			high = mid // Udah lewat target, root di kiri
		}
	}
	return low
}

func (a *Adapter) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, time.UTC)
	// Kita mulai pencarian dari noon lokal untuk menghindari moonset hari sebelumnya
	localStart := targetDate.Add(time.Duration(-lon/15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	et := a.findRoot(startET, 86400, -0.583*(math.Pi/180.0), Moon, lat, lon)
	return Et2Utc(et), nil
}

func (a *Adapter) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	// Start at 00:00 LOCAL TIME of that day
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon/15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
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
