package shared

// Location represents a geographical coordinate
type Location struct {
	Latitude  float64
	Longitude float64
}

// Validation response
type ValidationResult struct {
	IsValid bool
	Message string
}
