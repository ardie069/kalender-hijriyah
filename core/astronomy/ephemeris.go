package astronomy

import (
	"fmt"
	"math"
	"os"
	"time"
)

type EphemerisManager struct {
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
	return getPosition(target, frame, CorrLTS, Earth, et)
}

// DayVectors: Penampung vektor geosentris Sun/Moon selama 24 jam
type DayVectors struct {
	Sun  []Vector3
	Moon []Vector3
	Date time.Time
}

// PrecalculateDayVectors: Hitung posisi geosentris Sun/Moon setiap jam untuk 48 jam (IAU_EARTH frame)
func (em *EphemerisManager) PrecalculateDayVectors(dt time.Time) *DayVectors {
	dv := &DayVectors{
		Date: dt,
		Sun:  make([]Vector3, 49),
		Moon: make([]Vector3, 49),
	}

	base := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.UTC)
	for h := 0; h <= 48; h++ {
		checkTime := base.Add(time.Duration(h) * time.Hour)
		et := TimeToEt(checkTime)

		sunPos, _ := em.GetGeocentric(Sun, et, "IAU_EARTH")
		moonPos, _ := em.GetGeocentric(Moon, et, "IAU_EARTH")

		dv.Sun[h] = sunPos
		dv.Moon[h] = moonPos
	}
	return dv
}

// InterpolateVector: Interpolasi linier antara dua jam (mendukung hingga 48 jam)
func (dv *DayVectors) InterpolateVector(body string, t time.Time) Vector3 {
	diff := t.Sub(time.Date(dv.Date.Year(), dv.Date.Month(), dv.Date.Day(), 0, 0, 0, 0, time.UTC))
	hours := diff.Hours()
	if hours < 0 {
		hours = 0
	} else if hours > 47.99 {
		hours = 47.99
	}

	hLow := int(math.Floor(hours))
	hHigh := hLow + 1
	f := hours - float64(hLow)

	var vLow, vHigh Vector3
	if body == "SUN" {
		vLow, vHigh = dv.Sun[hLow], dv.Sun[hHigh]
	} else {
		vLow, vHigh = dv.Moon[hLow], dv.Moon[hHigh]
	}

	return Vector3{
		X: vLow.X + (vHigh.X-vLow.X)*f,
		Y: vLow.Y + (vHigh.Y-vLow.Y)*f,
		Z: vLow.Z + (vHigh.Z-vLow.Z)*f,
	}
}
