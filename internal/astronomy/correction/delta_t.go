package correction

import "time"

// CalculateDeltaT returns the Delta T value for a given time.
// Delta T is the difference between Terrestrial Time (TT) and Coordinated Universal Time (UTC).
// Currently uses a fixed approximation for ~2024-2026.
func CalculateDeltaT(t time.Time) float64 {
	return 69.184
}
