package astronomy

/*
#cgo LDFLAGS: -l:cspice.a -l:csupport.a -lm
#include "SpiceUsr.h"
#include <stdlib.h>

void silence_spice_c() {
    erract_c("SET", 0, "REPORT");
    errdev_c("SET", 0, "NULL");
}

void get_position_c(const char *t, double et, const char *f, const char *c, const char *o, double *pos, double *lt) {
    spkpos_c(t, et, f, c, o, pos, lt);
}

void get_msg_c(const char *type, int len, char *buffer) {
    getmsg_c(type, len, buffer);
}
*/
import "C"

import (
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"
)

// SATU KUNCI UNTUK SEMUA: SPICE tidak thread-safe sama sekali.
var spiceMu sync.Mutex

func getSpiceErrorInternal() error {
	var buffer [1024]C.char
	cLong := C.CString("LONG")
	defer C.free(unsafe.Pointer(cLong))
	C.get_msg_c(cLong, 1024, &buffer[0])
	return fmt.Errorf("NASA SPICE ERROR: %s", C.GoString(&buffer[0]))
}

// LoadKernel: Inisialisasi file kernel dan bungkam output terminal SPICE
func LoadKernel(path string) error {
	spiceMu.Lock()
	defer spiceMu.Unlock()

	C.silence_spice_c()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	C.reset_c()
	C.furnsh_c(cPath)

	if C.failed_c() != 0 {
		return getSpiceErrorInternal()
	}
	return nil
}

// Str2et: Ubah UTC String ke Ephemeris Time (TDB) - GUNAKAN HANYA JIKA PERLU
func Str2et(utcStr string) (float64, error) {
	spiceMu.Lock()
	defer spiceMu.Unlock()

	var et C.SpiceDouble
	cStr := C.CString(utcStr)
	defer C.free(unsafe.Pointer(cStr))

	C.reset_c()
	C.str2et_c(cStr, &et)

	if C.failed_c() != 0 {
		return 0, getSpiceErrorInternal()
	}
	return float64(et), nil
}

// TimeToEt: Versi Matematika Murni (TIDAK BUTUH MUTEX)
// Sangat disarankan untuk digunakan di dalam loop scanning agar tidak rebutan Mutex
func TimeToEt(t time.Time) float64 {
	// Konversi UTC ke TDB (pendekatan J2000)
	// J2000 Epoch: 946728000 Unix
	// Delta T 2026: ~69.184 detik
	const unixJ2000 = 946728000
	const deltaT = 69.184
	return float64(t.Unix()) - unixJ2000 + deltaT
}

func getPositionInternal(target, frame, corr, observer string, et float64) (Vector3, error) {
	// FUNGSI INI HARUS DIPANGGIL DI DALAM BLOK LOCK
	var pos [3]C.SpiceDouble
	var lt C.SpiceDouble

	t, f, c, o := C.CString(target), C.CString(frame), C.CString(corr), C.CString(observer)
	defer func() {
		C.free(unsafe.Pointer(t))
		C.free(unsafe.Pointer(f))
		C.free(unsafe.Pointer(c))
		C.free(unsafe.Pointer(o))
	}()

	C.reset_c()
	C.get_position_c(t, C.double(et), f, c, o, (*C.double)(&pos[0]), (*C.double)(&lt))
	if C.failed_c() != 0 {
		return Vector3{}, getSpiceErrorInternal()
	}
	return Vector3{float64(pos[0]), float64(pos[1]), float64(pos[2])}, nil
}

func getPosition(target, frame, corr, observer string, et float64) (Vector3, error) {
	spiceMu.Lock()
	defer spiceMu.Unlock()
	return getPositionInternal(target, frame, corr, observer, et)
}

func (em *EphemerisManager) GetTopocentricPosition(target string, et, lat, lon float64) (Vector3, error) {
	spiceMu.Lock()
	defer spiceMu.Unlock()

	// Math Geodetic ke Rectangular (Bisa di luar Lock sebenarnya, tapi biar aman satu blok)
	re := 6378.137
	f := 1.0 / 298.257223563
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0
	cosLat := math.Cos(latRad)
	sinLat := math.Sin(latRad)
	N := re / math.Sqrt(1-f*(2-f)*sinLat*sinLat)
	altKM := 0.45
	obsX := (N + altKM) * cosLat * math.Cos(lonRad)
	obsY := (N + altKM) * cosLat * math.Sin(lonRad)
	obsZ := (N*(1-f)*(1-f) + altKM) * sinLat

	posGeo, err := getPositionInternal(target, "IAU_EARTH", "LT+S", "EARTH", et)
	if err != nil {
		return Vector3{}, err
	}

	return Vector3{
		X: posGeo.X - obsX,
		Y: posGeo.Y - obsY,
		Z: posGeo.Z - obsZ,
	}, nil
}

func (em *EphemerisManager) GetLocalAltAz(pos Vector3, lat, lon float64) (float64, float64) {
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	sinLat, cosLat := math.Sin(latRad), math.Cos(latRad)
	sinLon, cosLon := math.Sin(lonRad), math.Cos(lonRad)

	// Transformasi ke Toposentrik (North-East-Down)
	x := -sinLat*cosLon*pos.X - sinLat*sinLon*pos.Y + cosLat*pos.Z
	y := -sinLon*pos.X + cosLon*pos.Y
	z := cosLat*cosLon*pos.X + cosLat*sinLon*pos.Y + sinLat*pos.Z

	alt := math.Asin(z/math.Sqrt(x*x+y*y+z*z)) * 180.0 / math.Pi
	az := math.Atan2(y, x) * 180.0 / math.Pi
	if az < 0 {
		az += 360.0
	}
	return alt, az
}

func Et2Utc(et float64) time.Time {
	// Inverse dari TimeToEt
	const unixJ2000 = 946728000
	const deltaT = 69.184
	unix := int64(et + unixJ2000 - deltaT)
	return time.Unix(unix, 0).UTC()
}

func (em *EphemerisManager) GetTimeByAltitude(date time.Time, lat, lon, targetAlt float64, rising bool) (time.Time, error) {
	// 1. Tentukan estimasi jam 12:00 siang UTC di lokasi tersebut
	// Noon UTC = 12.0 - (lon / 15.0)
	approxNoonUTC := 12.0 - (lon / 15.0)
	
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	noon := base.Add(time.Duration(approxNoonUTC*3600) * time.Second)

	var start, end time.Time
	if rising {
		// Cari dari 6 jam sebelum noon sampe noon (Subuh/Terbit)
		start = noon.Add(-7 * time.Hour)
		end = noon.Add(1 * time.Hour)
	} else {
		// Cari dari noon sampe 7 jam setelah noon (Maghrib/Isya)
		start = noon.Add(-1 * time.Hour)
		end = noon.Add(7 * time.Hour)
	}

	// Bisection search
	low := TimeToEt(start)
	high := TimeToEt(end)

	for range 20 {
		mid := (low + high) / 2
		pos, _ := em.GetTopocentricPosition("SUN", mid, lat, lon)
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

func (em *EphemerisManager) GetSolarTransit(date time.Time, lat, lon float64) time.Time {
	// Estimasi noon UTC
	approxNoonUTC := 12.0 - (lon / 15.0)
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	noon := base.Add(time.Duration(approxNoonUTC*3600) * time.Second)

	// Cari altitude tertinggi di sekitar estimasi noon (+/- 2 jam)
	low := TimeToEt(noon.Add(-2 * time.Hour))
	high := TimeToEt(noon.Add(2 * time.Hour))

	for range 20 {
		m1 := low + (high-low)*0.382
		m2 := low + (high-low)*0.618

		pos1, _ := em.GetTopocentricPosition("SUN", m1, lat, lon)
		alt1, _ := em.GetLocalAltAz(pos1, lat, lon)

		pos2, _ := em.GetTopocentricPosition("SUN", m2, lat, lon)
		alt2, _ := em.GetLocalAltAz(pos2, lat, lon)

		if alt1 < alt2 {
			low = m1
		} else {
			high = m2
		}
	}
	return Et2Utc((low + high) / 2)
}

func (em *EphemerisManager) GetAsrTime(dhuhr time.Time, lat, lon float64) time.Time {
	// 1. Cari altitude matahari pas Dhuhr buat dapet panjang bayangan minimum
	etDhuhr := TimeToEt(dhuhr)
	posDhuhr, _ := em.GetTopocentricPosition("SUN", etDhuhr, lat, lon)
	altDhuhr, _ := em.GetLocalAltAz(posDhuhr, lat, lon)

	// Zenith Angle (Z) = 90 - Altitude
	zenithDhuhrRad := (90.0 - altDhuhr) * math.Pi / 180.0

	// Kriteria Syafi'i: bayangan = 1 (panjang benda) + tan(zenithDhuhr)
	asrShadow := 1.0 + math.Tan(zenithDhuhrRad)
	targetAsrAlt := math.Atan(1.0/asrShadow) * 180.0 / math.Pi

	// 2. Cari kapan matahari nyentuh altitude tersebut di sore hari
	asrTime, _ := em.GetTimeByAltitude(dhuhr, lat, lon, targetAsrAlt, false)
	return asrTime
}

func (em *EphemerisManager) GetSunrise(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, true)
}

func (em *EphemerisManager) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	return em.GetTimeByAltitude(date, lat, lon, -0.833, false)
}
