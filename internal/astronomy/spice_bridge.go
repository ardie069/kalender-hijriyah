package astronomy

/*
#cgo CFLAGS: -I/home/ardie069/develop/cspice/include
#cgo LDFLAGS: /home/ardie069/develop/cspice/lib/cspice.a -lm
#include "SpiceUsr.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// LoadKernel: Inisialisasi file de440.bsp dengan laporan error detail
func LoadKernel(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	C.reset_c()
	C.furnsh_c(cPath)

	if C.failed_c() != 0 {
		var buffer [1024]C.char
		C.getmsg_c(C.CString("LONG"), 1024, &buffer[0])
		errMsg := C.GoString(&buffer[0])
		C.reset_c()
		return fmt.Errorf("\n[NASA SPICE ERROR]\nDetail: %s", errMsg)
	}
	return nil
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

// GetPosition: Ambil koordinat XYZ (KM) dari target relatif ke Bumi
func GetPosition(target string, et float64) ([3]float64, error) {
	var pos [3]C.SpiceDouble
	var lt C.SpiceDouble

	t := C.CString(target)
	f := C.CString("J2000")
	a := C.CString("LT+S")
	o := C.CString("EARTH")
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(f))
	defer C.free(unsafe.Pointer(a))
	defer C.free(unsafe.Pointer(o))

	C.reset_c()
	C.spkezr_c(t, C.SpiceDouble(et), f, a, o, &pos[0], &lt)

	if C.failed_c() != 0 {
		return [3]float64{}, fmt.Errorf("NASA SPICE: Gagal ambil posisi %s", target)
	}

	return [3]float64{float64(pos[0]), float64(pos[1]), float64(pos[2])}, nil
}
