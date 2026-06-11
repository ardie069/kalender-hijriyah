package correction

import "math"

// ApplyRefraction calculates the atmospheric refraction correction in degrees
// for a given true altitude (in degrees).
func ApplyRefraction(geoAlt float64) float64 {
	if geoAlt < -1.0 {
		return 0
	}

	denominator := math.Tan((geoAlt + 10.3/(geoAlt+5.11)) * math.Pi / 180.0)
	if math.Abs(denominator) < 1e-5 {
		return 0
	}

	rArcmin := 1.02 / denominator

	if rArcmin > 60.0 {
		rArcmin = 34.0
	} else if rArcmin < 0 {
		rArcmin = 0
	}

	return rArcmin / 60.0
}
