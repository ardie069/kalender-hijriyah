package astronomy

import (
	"math"
	"time"
)

// GetTimeByAltitude: Cari waktu ketika matahari mencapai altitude tertentu (bisection search)
func (em *EphemerisManager) GetTimeByAltitude(date time.Time, lat, lon, targetAlt float64, rising bool) (time.Time, error) {
	approxNoonUTC := 12.0 - (lon / 15.0)
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	noon := base.Add(time.Duration(approxNoonUTC*3600) * time.Second)

	var start, end time.Time
	if rising {
		start = noon.Add(-10 * time.Hour)
		end = noon.Add(2 * time.Hour)
	} else {
		start = noon.Add(-2 * time.Hour)
		end = noon.Add(10 * time.Hour)
	}

	low := TimeToEt(start)
	high := TimeToEt(end)

	for range 30 {
		mid := (low + high) / 2
		pos, err := em.GetTopocentricPosition("SUN", mid, lat, lon)
		if err != nil {
			return time.Time{}, err
		}
		currentAlt, _ := em.GetLocalAltAz(pos, lat, lon)

		if rising {
			if currentAlt < targetAlt {
				low = mid
			} else {
				high = mid
			}
		} else {
			if currentAlt > targetAlt {
				low = mid
			} else {
				high = mid
			}
		}
	}

	return Et2Utc((low + high) / 2), nil
}

// GetSolarTransit: Cari waktu kulminasi matahari / Dhuhr (golden section search)
func (em *EphemerisManager) GetSolarTransit(date time.Time, lat, lon float64) (time.Time, error) {
	approxNoonUTC := 12.0 - (lon / 15.0)
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	noon := base.Add(time.Duration(approxNoonUTC*3600) * time.Second)

	low := TimeToEt(noon.Add(-2 * time.Hour))
	high := TimeToEt(noon.Add(2 * time.Hour))

	for range 20 {
		m1 := low + (high-low)*0.382
		m2 := low + (high-low)*0.618

		pos1, err := em.GetTopocentricPosition("SUN", m1, lat, lon)
		if err != nil {
			return time.Time{}, err
		}
		alt1, _ := em.GetLocalAltAz(pos1, lat, lon)

		pos2, err := em.GetTopocentricPosition("SUN", m2, lat, lon)
		if err != nil {
			return time.Time{}, err
		}
		alt2, _ := em.GetLocalAltAz(pos2, lat, lon)

		if alt1 < alt2 {
			low = m1
		} else {
			high = m2
		}
	}
	return Et2Utc((low + high) / 2), nil
}

// GetAsrTime: Hitung waktu Ashar dengan shadow factor default (Syafi'i = 1)
func (em *EphemerisManager) GetAsrTime(dhuhr time.Time, lat, lon float64) (time.Time, error) {
	return em.GetAsrTimeWithFactor(dhuhr, lat, lon, 1)
}

// GetAsrTimeWithFactor menghitung Asr dengan shadow factor kustom.
// shadowFactor = 1 → Syafi'i/Maliki/Hanbali, shadowFactor = 2 → Hanafi
func (em *EphemerisManager) GetAsrTimeWithFactor(dhuhr time.Time, lat, lon float64, shadowFactor int) (time.Time, error) {
	// 1. Cari altitude matahari pas Dhuhr buat dapet panjang bayangan minimum
	etDhuhr := TimeToEt(dhuhr)
	posDhuhr, err := em.GetTopocentricPosition("SUN", etDhuhr, lat, lon)
	if err != nil {
		return time.Time{}, err
	}
	altDhuhr, _ := em.GetLocalAltAz(posDhuhr, lat, lon)

	// Zenith Angle (Z) = 90 - Altitude
	zenithDhuhrRad := (90.0 - altDhuhr) * math.Pi / 180.0

	// Kriteria: bayangan = shadowFactor (panjang benda) + tan(zenithDhuhr)
	asrShadow := float64(shadowFactor) + math.Tan(zenithDhuhrRad)
	targetAsrAlt := math.Atan(1.0/asrShadow) * 180.0 / math.Pi

	// 2. Cari kapan matahari nyentuh altitude tersebut di sore hari
	return em.GetTimeByAltitude(dhuhr, lat, lon, targetAsrAlt, false)
}

// GetSunrise: Waktu terbit matahari (altitude = -0.833°)
func (em *EphemerisManager) GetSunrise(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, true)
}

// GetSunset: Waktu terbenam matahari (altitude = -0.833°)
func (em *EphemerisManager) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, false)
}

// GetIsha: Waktu Isya (altitude = -18.0°)
func (em *EphemerisManager) GetIsha(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -18.0, false)
}
// GetMoonset: Waktu terbenam bulan (altitude = -0.583° untuk standar toposentrik)
func (em *EphemerisManager) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.583, false)
}

// GetFajr: Waktu Subuh (altitude = -18.0°)
func (em *EphemerisManager) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -18.0, true)
}
