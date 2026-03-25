package models

import "time"

type MoonTelemetry struct {
	Altitude         float64   `json:"altitude"` // Geometric (Geo or Topo depending on request)
	AltitudeApparent *float64   `json:"altitude_apparent,omitempty"`
	Azimuth          float64   `json:"azimuth"`
	Elongation       float64   `json:"elongation"`
	Illumination     float64   `json:"illumination"`
	DistanceKM       float64   `json:"distance_km"`
	AgeHours         float64   `json:"age_hours"`
	PhaseName        string    `json:"phase_name"`
	Timestamp        time.Time `json:"timestamp"`
}
