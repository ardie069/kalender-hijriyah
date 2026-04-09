package models

import "time"

type MoonTelemetry struct {
	Altitude         float64   `json:"altitude"` // Geometric (Geo or Topo depending on request)
	AltitudeApparent *float64   `json:"altitude_apparent,omitempty"`
	Azimuth          float64   `json:"azimuth"`
	Elongation       float64   `json:"elongation"`
	Illumination     float64   `json:"illumination"`
	DistanceKM       float64   `json:"distance_km"`
	AgeHours         float64    `json:"age_hours"`
	PhaseName        string     `json:"phase_name"`
	Moonrise         *time.Time `json:"moonrise,omitempty"`
	Moonset          *time.Time `json:"moonset,omitempty"`
	Timestamp        time.Time  `json:"timestamp"`
}

type VisibilityPoint struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Category string  `json:"category"` // A, B, C, D, E, F
	Altitude float64 `json:"alt"`
	Elong    float64 `json:"elong"`
	ArcV     float64 `json:"arcv"`   // Moon Alt - Sun Alt
	Width    float64 `json:"width"`  // Crescent width in degrees
	SunsetUTC float64 `json:"sunset_utc"` // hours since start of day
}

type VisibilityGrid struct {
	Date         time.Time         `json:"date"`
	MonthName    string            `json:"month_name"`
	Year         int               `json:"year"`
	Method       string            `json:"method"`
	Points       []VisibilityPoint `json:"points"`
	IjtimaTime   time.Time         `json:"ijtima_time"`
	FajarNZTime  time.Time         `json:"fajar_nz_time"`
	BestLocation *LocationInfo      `json:"best_location,omitempty"`
}


