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

func (a *Adapter) CalculateGeocentricParams(sun, moon Vector3) (altitude, elongation float64) {
	elongation = math.Acos(sun.Unit().Dot(moon.Unit())) * (180.0 / math.Pi)

	altitude = math.Asin(moon.Unit().Z) * (180.0 / math.Pi)

	return altitude, elongation
}

// GetFajr: Mencari waktu subuh (Solar Altitude = -18 derajat)
func (a *Adapter) GetFajr(date time.Time, latitude, longitude float64) (time.Time, error) {
	startET, _ := Str2et(date.Format("2006-01-02") + " 00:00:00 UTC")
	// Subuh: Matahari naik mendekati ufuk, target -18 derajat
	targetAlt := -18.0 * (math.Pi / 180.0)

	// Cari di 12 jam pertama hari itu
	low, high := startET, startET+43200
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, _ := a.Manager.GetTopocentricPosition(Sun, mid, latitude, longitude)
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, latitude, longitude)
		currentAlt := currentAltDeg * (math.Pi / 180.0)
		if currentAlt < targetAlt {
			low = mid
		} else {
			high = mid
		}
	}
	return Et2Utc(low), nil
}

func (a *Adapter) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	startET, _ := Str2et(date.Format("2006-01-02") + " 12:00:00 UTC")
	et := a.findRoot(startET, 86400, -0.583*(math.Pi/180.0), Moon, lat, lon)
	return Et2Utc(et), nil
}

func (a *Adapter) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	startET, _ := Str2et(date.Format("2006-01-02") + " 06:00:00 UTC")
	et := a.findRoot(startET, 43200, -0.833*(math.Pi/180.0), Sun, lat, lon)
	return Et2Utc(et), nil
}
