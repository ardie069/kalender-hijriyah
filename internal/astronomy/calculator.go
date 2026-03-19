package astronomy

import (
	"math"
	"time"
)

type Adapter struct {
	Manager *EphemerisManager
}

// findRoot: Fungsi Bisection generik buat nyari waktu kejadian
func (a *Adapter) findRoot(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) (float64, error) {
	low, high := startET, startET+duration
	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, err := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		if err != nil {
			return 0, err
		}
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)
		currentAlt := currentAltDeg * (math.Pi / 180.0)

		if currentAlt > targetAlt {
			low = mid
		} else {
			high = mid
		}
	}
	return low, nil
}

func (a *Adapter) CalculateGeocentricParams(et float64, lat, lon float64) (altitude, elongation float64) {
	sunPos, err := a.Manager.GetGeocentric(Sun, et, "J2000")
	if err != nil {
		return 0, 0
	}
	moonPos, err := a.Manager.GetGeocentric(Moon, et, "J2000")
	if err != nil {
		return 0, 0
	}

	elongation = math.Acos(sunPos.Unit().Dot(moonPos.Unit())) * (180.0 / math.Pi)

	topoPos, err := a.Manager.GetTopocentricPosition(Moon, et, lat, lon)
	if err != nil {
		return 0, elongation
	}
	altitude, _ = a.Manager.GetLocalAltAz(topoPos, lat, lon)

	return altitude, elongation
}

// ... (GetFajr remains similar but handle findRootUp error)
func (a *Adapter) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon / 15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	targetAlt := -18.0 * (math.Pi / 180.0)

	et, err := a.findRootUp(startET, 86400, targetAlt, Sun, lat, lon)
	if err != nil {
		return time.Time{}, err
	}
	return Et2Utc(et), nil
}

// findRootUp: Bisection naik khusus untuk Sunrise/Fajr
func (a *Adapter) findRootUp(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) (float64, error) {
	targetAltDeg := targetAlt * (180.0 / math.Pi)
	var low, high float64
	foundInterval := false

	for offset := 0.0; offset < duration; offset += 3600.0 {
		pos1, err := a.Manager.GetTopocentricPosition(target, startET+offset, lat, lon)
		if err != nil {
			return 0, err
		}
		alt1, _ := a.Manager.GetLocalAltAz(pos1, lat, lon)

		pos2, err := a.Manager.GetTopocentricPosition(target, startET+offset+3600.0, lat, lon)
		if err != nil {
			return 0, err
		}
		alt2, _ := a.Manager.GetLocalAltAz(pos2, lat, lon)

		if alt1 < targetAltDeg && alt2 >= targetAltDeg {
			low = startET + offset
			high = startET + offset + 3600.0
			foundInterval = true
			break
		}
	}

	if !foundInterval {
		return startET + duration/2, nil // Fallback
	}

	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, err := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		if err != nil {
			return 0, err
		}
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)

		if currentAltDeg < targetAltDeg {
			low = mid
		} else {
			high = mid
		}
	}
	return low, nil
}

func (a *Adapter) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon / 15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	et, err := a.findRoot(startET, 86400, -0.583*(math.Pi/180.0), Moon, lat, lon)
	if err != nil {
		return time.Time{}, err
	}
	return Et2Utc(et), nil
}

func (a *Adapter) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon / 15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	et, err := a.findRootDown(startET, 86400, -0.833*(math.Pi/180.0), Sun, lat, lon)
	if err != nil {
		return time.Time{}, err
	}
	return Et2Utc(et), nil
}

func (a *Adapter) GetIsha(date time.Time, lat, lon float64) (time.Time, error) {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	localStart := targetDate.Add(time.Duration(-lon / 15.0 * float64(time.Hour)))
	startET := TimeToEt(localStart)
	et, err := a.findRootDown(startET, 86400, -18.0*(math.Pi/180.0), Sun, lat, lon)
	if err != nil {
		return time.Time{}, err
	}
	return Et2Utc(et), nil
}

// findRootDown: Bisection turun khusus untuk Sunset/Moonset
func (a *Adapter) findRootDown(startET float64, duration float64, targetAlt float64, target string, lat, lon float64) (float64, error) {
	var low, high float64
	foundInterval := false

	for offset := 0.0; offset < duration; offset += 3600.0 {
		pos1, err := a.Manager.GetTopocentricPosition(target, startET+offset, lat, lon)
		if err != nil {
			return 0, err
		}
		alt1, _ := a.Manager.GetLocalAltAz(pos1, lat, lon)

		pos2, err := a.Manager.GetTopocentricPosition(target, startET+offset+3600.0, lat, lon)
		if err != nil {
			return 0, err
		}
		alt2, _ := a.Manager.GetLocalAltAz(pos2, lat, lon)

		if alt1 > targetAlt*(180.0/math.Pi) && alt2 <= targetAlt*(180.0/math.Pi) {
			low = startET + offset
			high = startET + offset + 3600.0
			foundInterval = true
			break
		}
	}

	if !foundInterval {
		return startET + duration/2, nil // Fallback
	}

	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		pos, err := a.Manager.GetTopocentricPosition(target, mid, lat, lon)
		if err != nil {
			return 0, err
		}
		currentAltDeg, _ := a.Manager.GetLocalAltAz(pos, lat, lon)
		currentAlt := currentAltDeg * (math.Pi / 180.0)

		if currentAlt > targetAlt {
			low = mid
		} else {
			high = mid
		}
	}
	return low, nil
}
