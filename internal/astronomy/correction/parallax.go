package correction

import "math"

// CalculateHorizontalParallax computes the horizontal parallax in degrees
// given the distance to the body (in km) and the Earth's radius (Re).
func CalculateHorizontalParallax(distanceKm, earthRadiusKm float64) float64 {
	if distanceKm <= 0 {
		return 0
	}
	return math.Asin(earthRadiusKm/distanceKm) * (180.0 / math.Pi)
}
