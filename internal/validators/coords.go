package validators

import (
	"fmt"
	"math"
)

func ValidateCoords(lat, lon float64) error {
	if math.IsNaN(lat) || math.IsNaN(lon) || math.IsInf(lat, 0) || math.IsInf(lon, 0) {
		return fmt.Errorf("latitude/longitude harus berupa angka valid")
	}
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude harus berada di rentang -90 s/d 90")
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude harus berada di rentang -180 s/d 180")
	}
	return nil
}

// For backward compatibility with handlers
func ValidateCoordinates(lat, lon float64) error {
	return ValidateCoords(lat, lon)
}
