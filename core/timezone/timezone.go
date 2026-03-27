package timezone

import (
	"fmt"
	"time"

	"github.com/ringsaturn/tzf"
)

// Service provides timezone resolution from coordinates.
type Service struct {
	finder tzf.F
}

// NewService creates a new timezone Service.
func NewService() (*Service, error) {
	f, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, err
	}
	return &Service{finder: f}, nil
}

// GetLocalTimeInfo converts a UTC time to local time string and timezone display name.
func (s *Service) GetLocalTimeInfo(utc time.Time, lat, lon float64) (string, string) {
	tzName := s.finder.GetTimezoneName(lon, lat)

	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return utc.Format("2006-01-02 15:04:05"), "UTC"
	}

	local := utc.In(loc)
	_, offset := local.Zone()
	return local.Format("2006-01-02 15:04:05"), fmt.Sprintf("%s (UTC %d)", local.Format("MST"), offset/3600)
}

// GetLocation returns a *time.Location for the given coordinates.
func (s *Service) GetLocation(lat, lon float64) *time.Location {
	tzName := s.finder.GetTimezoneName(lon, lat)
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.UTC
	}
	return loc
}
