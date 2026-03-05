package astronomy

import (
	"fmt"
	"os"
)

// NAIF ID (Tetep simpen di sini buat referensi)
const (
	NAIF_SUN   = "SUN"
	NAIF_EARTH = "EARTH"
	NAIF_MOON  = "MOON"
)

type EphemerisManager struct {
	Kernels []string
}

// NewEphemerisManager bakal mastiin semua file kernel ada sebelum aplikasi jalan
func NewEphemerisManager(kernelPaths ...string) (*EphemerisManager, error) {
	for _, path := range kernelPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("kernel NASA ga ketemu, Bang: %s", path)
		}
		// Langsung load ke CSPICE via bridge
		if err := LoadKernel(path); err != nil {
			return nil, err
		}
	}
	return &EphemerisManager{Kernels: kernelPaths}, nil
}

// GetPosition sekarang manggil fungsi NASA asli di spice_bridge.go
func (em *EphemerisManager) GetPosition(target string, et float64) (Vector3, error) {
	pos, err := GetPosition(target, et)
	if err != nil {
		return Vector3{}, err
	}
	// NASA balikin KM, kita bungkus jadi Vector3 biar bisa pake Dot/Norm
	return Vector3{X: pos[0], Y: pos[1], Z: pos[2]}, nil
}
