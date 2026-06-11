package ephemeris

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Service struct {
	Astro   *cspice.Adapter
	Manager *cspice.EphemerisManager
}

func NewService(astro *cspice.Adapter, manager *cspice.EphemerisManager) *Service {
	return &Service{
		Astro:   astro,
		Manager: manager,
	}
}

func (s *Service) FindIjtima(approxDate time.Time) (time.Time, error) {
	et, _ := cspice.Str2et(approxDate.Format(cspice.TimeFormat))
	diff, err := s.getLongitudeDiff(et)
	if err != nil {
		return time.Time{}, err
	}

	daysDiff := diff / (12.19 * math.Pi / 180.0)
	estIjtimaET := et - (daysDiff * 86400.0)

	low := estIjtimaET - (2 * 86400.0)
	high := estIjtimaET + (2 * 86400.0)
	var mid float64

	for range 60 {
		mid = (low + high) / 2
		d, err := s.getLongitudeDiff(mid)
		if err != nil {
			return time.Time{}, err
		}
		if d > 0 {
			high = mid
		} else {
			low = mid
		}
	}
	return cspice.Et2Utc(mid), nil
}

func (s *Service) FindPreviousIjtima(t time.Time) (time.Time, error) {
	return s.FindIjtima(t.AddDate(0, 0, -16))
}

func (s *Service) getLongitudeDiff(et float64) (float64, error) {
	sunPos, err := s.Manager.GetGeocentric(cspice.Sun, et, cspice.FrameEclipJ2000)
	if err != nil {
		return 0, err
	}
	moonPos, err := s.Manager.GetGeocentric(cspice.Moon, et, cspice.FrameEclipJ2000)
	if err != nil {
		return 0, err
	}

	sunLong := math.Atan2(sunPos.Y, sunPos.X)
	moonLong := math.Atan2(moonPos.Y, moonPos.X)

	diff := moonLong - sunLong
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}

	return diff, nil
}
