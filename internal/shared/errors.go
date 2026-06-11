package shared

import "errors"

// Common domain errors
var (
	ErrInvalidMethod = errors.New("invalid calculation method provided")
	ErrLocationNotFound = errors.New("location parameters are missing or invalid")
	ErrEphemerisLoad = errors.New("failed to load ephemeris data")
)
