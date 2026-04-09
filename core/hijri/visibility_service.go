package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// GenerateVisibilityGrid produces a grid of visibility points for a given date and method.
func (s *DateService) GenerateVisibilityGrid(searchDate time.Time, method string) (*models.VisibilityGrid, error) {
	grid := &models.VisibilityGrid{
		Date:   searchDate,
		Method: method,
		Points: []models.VisibilityPoint{},
	}

	// Grid configuration
	// Step 4.0 provides a perfect balance of speed (~1.5s vs ~10s) and visual detail
	latStart, latEnd, latStep := 65.0, -65.0, 4.0
	lonStart, lonEnd, lonStep := -180.0, 180.0, 4.0
	
	// Normalize searchDate to 00:00:00 UTC to avoid timezone/time-of-day artifacts
	searchDate = time.Date(searchDate.Year(), searchDate.Month(), searchDate.Day(), 0, 0, 0, 0, time.UTC)

	// --- 1. PRIORITY CALCULATIONS (METADATA) ---
	// Use FindPreviousIjtima to always anchor to the conjunction that started 
	// the current Hijri cycle, preventing mid-month jumps.
	ijtima, err := s.Cal.FindPreviousIjtima(searchDate)
	if err != nil {
		ijtima = searchDate // Fallback
	}
	grid.IjtimaTime = ijtima

	// Fajar NZ: Dawn in New Zealand (Wellington -41.2865, 174.7762)
	// For KHGT, we compare Ijtima with the Fajr that marks the start of the 
	// potential new global month (the Fajr following the conjunction).
	fajarNZ, _ := s.Astro.GetFajr(ijtima.AddDate(0, 0, 1), -41.2865, 174.7762)
	grid.FajarNZTime = fajarNZ

	// --- 2. HIJRI LABELING ---
	// Determine which Hijri Month we are investigating based on Ijtima
	hDate := calendar.GetTabularHijri(ijtima.AddDate(0, 0, 1))
	grid.MonthName = hDate.MonthName
	grid.Year = hDate.Year

	// --- 3. OBSERVATION ANCHORING ---
	// To prevent the "2 May" disconnect, we MUST perform the grid scan 
	// on the actual evening where the Hilal sighting is relevant.
	// This is typically the evening of the found Ijtima.
	observationDate := ijtima
	grid.Date = observationDate // Update visual date to match anchor

	// --- 4. GRID GENERATION (VISUALIZATION) ---
	var bestLoc *models.LocationInfo
	var minLonFound = 999.0

	for lat := latStart; lat >= latEnd; lat -= latStep {
		for lon := lonStart; lon <= lonEnd; lon += lonStep {
			// CRITICAL: We use observationDate (H-29) instead of the input searchDate
			sunset, err := s.Astro.GetSunsetFast(observationDate, lat, lon)
			if err != nil {
				continue
			}

			// Ijtima check: Hilal only visible after ijtima
			if sunset.Before(ijtima) {
				grid.Points = append(grid.Points, models.VisibilityPoint{
					Lat:      lat,
					Lon:      lon,
					Category: "F",
				})
				continue
			}

			// Calculate Params based on Method (Topo vs Geo)
			category := "F"
			var alt, elong, arcv, width float64

			if method == "KHGT" {
				// KHGT: Geocentric standard (Turkish Congress 2016)
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
				// ODEH: Topocentric standard (Refraction + Parallax)
				alt, elong, arcv, width = s.Astro.CalculateTopocentricParamsGlobal(sunset, lat, lon)
				
				// Standard Odeh width (arcmin)
				wArcmin := width * 60.0
				cat := calendar.GetOdehCategory(arcv, wArcmin)
				
				// Standardize Categories to prevent frontend overlap
				if cat != "E" && cat != "F" && cat != "G" {
					category = "ODEH_" + cat
				} else {
					category = cat // E, F, or G
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
