package models

import "time"

type MoonTelemetry struct {
	Altitude     float64   `json:"altitude"`
	Azimuth      float64   `json:"azimuth"`
	Elongation   float64   `json:"elongation"`
	Illumination float64   `json:"illumination"`
	DistanceKM   float64   `json:"distance_km"`
	PhaseName    string    `json:"phase_name"`
	Timestamp    time.Time `json:"timestamp"`
}
