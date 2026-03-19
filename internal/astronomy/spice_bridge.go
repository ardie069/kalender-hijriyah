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

func ApplyRefraction(geoAlt float64) float64 {
	// Jika terlalu rendah, jangan hitung (karena rumus Tan bisa meledak)
	if geoAlt < -1.0 {
		return 0
	}

	denominator := math.Tan((geoAlt + 10.3/(geoAlt+5.11)) * math.Pi / 180.0)
	if math.Abs(denominator) < 1e-5 {
		return 0 // Safety: hindari div-by-zero
	}

	R_arcmin := 1.02 / denominator

	// Limitasi: Refraksi tidak boleh terlalu besar secara fisik di aplikasi ini
	if R_arcmin > 60.0 {
		R_arcmin = 34.0 // Nilai standar horizon
	} else if R_arcmin < 0 {
		R_arcmin = 0
	}

	return R_arcmin / 60.0
}
