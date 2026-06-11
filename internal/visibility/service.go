package visibility

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy/ephemeris"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/models"
	"github.com/ardie069/kalender-hijriyah/internal/visibility/criteria"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Service struct {
	Astro *cspice.Adapter
	Ephem *ephemeris.Service
}

func NewService(astro *cspice.Adapter, ephem *ephemeris.Service) *Service {
	return &Service{
		Astro: astro,
		Ephem: ephem,
	}
}

// GenerateVisibilityGrid produces a grid of visibility points for a given date and method.
func (s *Service) GenerateVisibilityGrid(searchDate time.Time, method string) (*models.VisibilityGrid, error) {
	grid := &models.VisibilityGrid{
		Date:   searchDate,
		Method: method,
		Points: []models.VisibilityPoint{},
	}

	latStart, latEnd, latStep := 65.0, -65.0, 4.0
	lonStart, lonEnd, lonStep := -180.0, 180.0, 4.0
	
	searchDate = time.Date(searchDate.Year(), searchDate.Month(), searchDate.Day(), 0, 0, 0, 0, time.UTC)

	ijtima, err := s.Ephem.FindPreviousIjtima(searchDate)
	if err != nil {
		ijtima = searchDate // Fallback
	}
	grid.IjtimaTime = ijtima

	fajarNZ, _ := s.Astro.GetFajr(ijtima.AddDate(0, 0, 1), -41.2865, 174.7762)
	grid.FajarNZTime = fajarNZ

	hDate := calendar.GetTabularHijri(ijtima.AddDate(0, 0, 1))
	grid.MonthName = hDate.MonthName
	grid.Year = hDate.Year

	observationDate := ijtima
	grid.Date = observationDate

	var bestLoc *models.LocationInfo
	var minLonFound = 999.0

	for lat := latStart; lat >= latEnd; lat -= latStep {
		for lon := lonStart; lon <= lonEnd; lon += lonStep {
			sunset, err := s.Astro.GetSunsetFast(observationDate, lat, lon)
			if err != nil {
				continue
			}

			if sunset.Before(ijtima) {
				grid.Points = append(grid.Points, models.VisibilityPoint{
					Lat:      lat,
					Lon:      lon,
					Category: "F",
				})
				continue
			}

			category := "F"
			var alt, elong, arcv, width float64

			if method == "KHGT" {
				alt, elong, arcv, width = s.Astro.CalculateGeocentricParamsGlobal(sunset, lat, lon)
				if alt >= 5.0 && elong >= 8.0 {
					category = "KHGT_YES"
					if lon < minLonFound {
						minLonFound = lon
						bestLoc = &models.LocationInfo{Lat: lat, Lon: lon}
					}
				} else {
					category = "KHGT_NO"
				}
			} else {
				alt, elong, arcv, width = s.Astro.CalculateTopocentricParamsGlobal(sunset, lat, lon)
				
				wArcmin := width * 60.0
				odeh := criteria.Odeh{}
				ctx := criteria.Context{
					Altitude: alt,
					Elongation: elong,
					ArcOfLight: arcv,
					Width: wArcmin,
				}
				cat := string(odeh.Evaluate(ctx))
				
				if cat != "E" && cat != "F" && cat != "G" {
					category = "ODEH_" + cat
				} else {
					category = cat
				}
			}

			grid.Points = append(grid.Points, models.VisibilityPoint{
				Lat:       lat,
				Lon:       lon,
				Category:  category,
				Altitude:  alt,
				Elong:     elong,
				ArcV:      arcv,
				Width:     width,
				SunsetUTC: float64(sunset.Hour()) + float64(sunset.Minute())/60.0 + float64(sunset.Second())/3600.0,
			})
		}
	}
	
	grid.BestLocation = bestLoc
	return grid, nil
}
