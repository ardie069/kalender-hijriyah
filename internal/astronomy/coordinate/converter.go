package coordinate

import "math"

// DegToRad converts degrees to radians
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// RadToDeg converts radians to degrees
func RadToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}
