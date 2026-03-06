package astronomy

/*
#cgo LDFLAGS: -l:cspice.a -l:csupport.a -lm
#include "SpiceUsr.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"math"
	"time"
	"unsafe"
)

// LoadKernel: Inisialisasi file de440.bsp dengan laporan error detail
func LoadKernel(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.reset_c()
	C.furnsh_c(cPath)
	if C.failed_c() != 0 {
		return getSpiceError()
	}
	return nil
}

func Et2Utc(et float64) time.Time {
	var buffer [64]C.char
	C.et2utc_c(C.SpiceDouble(et), C.CString("ISOC"), 3, 64, &buffer[0])

	utcStr := C.GoString(&buffer[0])
	t, _ := time.Parse("2006-01-02T15:04:05.000", utcStr)
	return t
}

// Str2et: Ubah UTC String ke Ephemeris Time (TDB)
func Str2et(utcStr string) (float64, error) {
	var et C.SpiceDouble
	cStr := C.CString(utcStr)
	defer C.free(unsafe.Pointer(cStr))

	C.reset_c()
	C.str2et_c(cStr, &et)

	if C.failed_c() != 0 {
		return 0, fmt.Errorf("NASA SPICE: Gagal konversi waktu %s", utcStr)
	}
	return float64(et), nil
}

func (em *EphemerisManager) GetEclipticPosition(target string, et float64) (Vector3, error) {
	var pos [3]C.SpiceDouble
	var lt C.SpiceDouble

	cTarget := C.CString(target)
	cFrame := C.CString("ECLIPJ2000")
	cCorr := C.CString("LT+S")
	cObs := C.CString("EARTH")

	defer C.free(unsafe.Pointer(cTarget))
	defer C.free(unsafe.Pointer(cFrame))
	defer C.free(unsafe.Pointer(cCorr))
	defer C.free(unsafe.Pointer(cObs))

	C.reset_c()
	C.spkezr_c(cTarget, C.SpiceDouble(et), cFrame, cCorr, cObs, &pos[0], &lt)

	if C.failed_c() != 0 {
		return Vector3{}, fmt.Errorf("NASA SPICE: Gagal ambil posisi ekliptika %s", target)
	}

	return Vector3{X: float64(pos[0]), Y: float64(pos[1]), Z: float64(pos[2])}, nil
}

// GetPosition: Ambil koordinat XYZ (KM) dari target relatif ke Bumi
func getPosition(target, frame, corr, observer string, et float64) (Vector3, error) {
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
	C.spkezr_c(t, C.SpiceDouble(et), f, c, o, &pos[0], &lt)
	if C.failed_c() != 0 {
		return Vector3{}, getSpiceError()
	}
	return Vector3{float64(pos[0]), float64(pos[1]), float64(pos[2])}, nil
}

// GetTopocentricPosition: Ambil posisi benda langit dari koordinat pengamat (Lat/Lon)
func (em *EphemerisManager) GetTopocentricPosition(target string, et, lat, lon float64) (Vector3, error) {
	// 1. Konversi Lat/Lon ke Rectangular (IAU_EARTH)
	// Pake konstanta Radii Bumi dari WGS84 biar gak "bola" doang
	re := 6378.137           // Radius Khatulistiwa (KM)
	f := 1.0 / 298.257223563 // Flattening Bumi

	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	// Geodetic to Rectangular (Model Ellipsoid)
	cosLat := math.Cos(latRad)
	sinLat := math.Sin(latRad)
	N := re / math.Sqrt(1-f*(2-f)*sinLat*sinLat)

	// Asumsikan ketinggian Malang rata-rata 450m (0.45 KM)
	altKM := 0.45
	obsX := (N + altKM) * cosLat * math.Cos(lonRad)
	obsY := (N + altKM) * cosLat * math.Sin(lonRad)
	obsZ := (N*(1-f)*(1-f) + altKM) * sinLat

	// 2. Ambil posisi target RELATIF terhadap bumi dalam frame IAU_EARTH
	// Frame IAU_EARTH ini ikut muter bareng bumi!
	posGeo, err := getPosition(target, "IAU_EARTH", "LT+S", "EARTH", et)
	if err != nil {
		return Vector3{}, err
	}

	// 3. Vektor Resultan (Target - Observer) di frame yang sama
	// Hasilnya udah murni posisi Toposentris (dari mata pengamat)
	return Vector3{
		X: posGeo.X - obsX,
		Y: posGeo.Y - obsY,
		Z: posGeo.Z - obsZ,
	}, nil
}

func getSpiceError() error {
	var buffer [1024]C.char
	C.getmsg_c(C.CString("LONG"), 1024, &buffer[0])
	errMsg := C.GoString(&buffer[0])
	C.reset_c()
	return fmt.Errorf("NASA SPICE ERROR: %s", errMsg)
}

// GetLocalAltAz: Versi Manual Projection (Gak butuh recazl_c)
func (em *EphemerisManager) GetLocalAltAz(v Vector3, lat, lon float64) (float64, float64) {
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	// 1. Zenith (Up) vector
	zenX := math.Cos(latRad) * math.Cos(lonRad)
	zenY := math.Cos(latRad) * math.Sin(lonRad)
	zenZ := math.Sin(latRad)

	// 2. North vector
	northX := -math.Sin(latRad) * math.Cos(lonRad)
	northY := -math.Sin(latRad) * math.Sin(lonRad)
	northZ := math.Cos(latRad)

	// 3. East vector
	eastX := -math.Sin(lonRad)
	eastY := math.Cos(lonRad)
	eastZ := 0.0

	// 4. Dot Product
	upProj := v.X*zenX + v.Y*zenY + v.Z*zenZ
	northProj := v.X*northX + v.Y*northY + v.Z*northZ
	eastProj := v.X*eastX + v.Y*eastY + v.Z*eastZ

	// 5. Hitung Altitude
	r := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	altRad := math.Asin(upProj / r)
	altitude := altRad * 180.0 / math.Pi

	// 6. Hitung Azimuth (from North towards East)
	azimuth := math.Atan2(eastProj, northProj) * 180.0 / math.Pi
	if azimuth < 0 {
		azimuth += 360.0
	}

	return altitude, azimuth
}
