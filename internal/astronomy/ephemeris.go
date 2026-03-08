package astronomy

import (
	"fmt"
	"os"
	"sync"
)

type EphemerisManager struct {
	mu sync.Mutex
	Kernels []string
}

func NewEphemerisManager(kernelPaths ...string) (*EphemerisManager, error) {
	em := &EphemerisManager{Kernels: kernelPaths}

	for _, path := range kernelPaths {
		// Cek filenya ada apa nggak
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("kernel ga ketemu: %s", path)
		}

		// Panggil LoadKernel dari spice_bridge.go
		if err := LoadKernel(path); err != nil {
			return nil, err
		}
	}
	return em, nil
}

func (em *EphemerisManager) GetGeocentric(target string, et float64, frame string) (Vector3, error) {
	em.mu.Lock()
	defer em.mu.Unlock()
	return getPosition(target, frame, CorrLTS, Earth, et)
}
