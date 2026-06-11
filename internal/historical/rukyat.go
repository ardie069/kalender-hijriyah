package historical

import "time"

// RukyatRecord represents an individual historical moon sighting report
type RukyatRecord struct {
	ID          string
	Date        time.Time
	Location    string
	Latitude    float64
	Longitude   float64
	IsSeen      bool
	Observer    string
	Description string
}
