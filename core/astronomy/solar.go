package astronomy

import (
	"math"
	"time"
)

// GetTimeByAltitude: Mencari waktu saat benda mencapai altitude tertentu dalam jendela 36 jam dari tengah malam.
func (em *EphemerisManager) GetTimeByAltitude(date time.Time, lat, lon, targetAlt float64, rising bool, body string) (time.Time, error) {
	// Calculate approximate local noon in UTC
	approxNoonUTC := 12.0 - (lon / 15.0)
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	noon := base.Add(time.Duration(approxNoonUTC*3600) * time.Second)

	var start, end time.Time
	if rising {
		// Morning events (Fajr, Sunrise): Search from local midnight (-12h from noon)
		start = noon.Add(-12 * time.Hour)
		end = noon.Add(6 * time.Hour)
	} else {
		// Evening events (Sunset, Isha): Search around and after local noon
		start = noon.Add(-6 * time.Hour)
		end = noon.Add(14 * time.Hour)
	}

	// Moon cycles are longer, extend window to ensure we don't miss crossing
	if body == "MOON" {
		end = start.Add(36 * time.Hour)
	}

	return em.FindAltitudeCrossing(body, start, end, lat, lon, targetAlt, rising)
}

// FindAltitudeCrossing: Scan interval demi interval untuk menemukan crossing altitude.
func (em *EphemerisManager) FindAltitudeCrossing(body string, start, end time.Time, lat, lon, targetAlt float64, rising bool) (time.Time, error) {
	step := time.Hour
	for t := start; t.Before(end); t = t.Add(step) {
		t1 := t
		t2 := t.Add(step)
		if t2.After(end) {
			t2 = end
		}

		et1 := TimeToEt(t1)
		et2 := TimeToEt(t2)

		pos1, _ := em.GetTopocentricPosition(body, et1, lat, lon)
		alt1, _ := em.GetLocalAltAz(pos1, lat, lon)
		if body == "SUN" || body == "MOON" {
			alt1 += ApplyRefraction(alt1)
		}

		pos2, _ := em.GetTopocentricPosition(body, et2, lat, lon)
		alt2, _ := em.GetLocalAltAz(pos2, lat, lon)
		if body == "SUN" || body == "MOON" {
			alt2 += ApplyRefraction(alt2)
		}

		if (rising && alt1 < targetAlt && alt2 > targetAlt) || (!rising && alt1 > targetAlt && alt2 < targetAlt) {
			low := et1
			high := et2
			for range 30 {
				mid := (low + high) / 2
				posM, _ := em.GetTopocentricPosition(body, mid, lat, lon)
				altM, _ := em.GetLocalAltAz(posM, lat, lon)
				if body == "SUN" || body == "MOON" {
					altM += ApplyRefraction(altM)
				}
				if rising {
					if altM < targetAlt {
						low = mid
					} else {
						high = mid
					}
				} else {
					if altM > targetAlt {
						low = mid
					} else {
						high = mid
					}
				}
			}
			return Et2Utc((low + high) / 2), nil
		}
	}
	return time.Time{}, nil
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
		alt1 += ApplyRefraction(alt1)

		pos2, err := em.GetTopocentricPosition("SUN", m2, lat, lon)
		if err != nil {
			return time.Time{}, err
		}
		alt2, _ := em.GetLocalAltAz(pos2, lat, lon)
		alt2 += ApplyRefraction(alt2)

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
	return em.GetTimeByAltitude(dhuhr, lat, lon, targetAsrAlt, false, "SUN")
}

// GetSunrise: Waktu terbit matahari (altitude = -0.833°)
func (em *EphemerisManager) GetSunrise(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, true, "SUN")
}

// GetSunset: Waktu terbenam matahari (altitude = -0.833°)
func (em *EphemerisManager) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, false, "SUN")
}

// GetIsha: Waktu Isya (altitude = -18.0°)
func (em *EphemerisManager) GetIsha(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -18.0, false, "SUN")
}

// GetMoonset: Waktu terbenam bulan (altitude = 0° sesuai permintaan pengguna untuk crossing horizon)
func (em *EphemerisManager) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, 0.0, false, "MOON")
}

// GetMoonrise: Waktu terbit bulan (altitude = 0° sesuai permintaan pengguna untuk crossing horizon)
func (em *EphemerisManager) GetMoonrise(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, 0.0, true, "MOON")
}

// GetFajr: Waktu Subuh (altitude = -18.0°)
func (em *EphemerisManager) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -18.0, true, "SUN")
}

// GetSunsetFast: Waktu terbenam matahari dengan iterasi lebih sedikit (untuk scan global)
func (em *EphemerisManager) GetSunsetFast(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitudeFast(date, lat, lon, -0.833, false, "SUN")
}

// GetTimeByAltitudeFast: Versi cepat dari GetTimeByAltitude dengan iterasi terbatas (12 iterasi ~ 3 menit akurasi)
func (em *EphemerisManager) GetTimeByAltitudeFast(date time.Time, lat, lon, targetAlt float64, rising bool, body string) (time.Time, error) {
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

	for i := 0; i < 12; i++ {
		mid := (low + high) / 2
		pos, err := em.GetTopocentricPosition(body, mid, lat, lon)
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
