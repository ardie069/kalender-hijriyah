package hijri

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// CalendarService handles yearly Hijri calendar generation.
type CalendarService struct {
	DateSvc *DateService
}

// NewCalendarService creates a new CalendarService.
func NewCalendarService(dateSvc *DateService) *CalendarService {
	return &CalendarService{DateSvc: dateSvc}
}

// GetYearlyCalendar generates a full Hijri year calendar for the given method.
func (c *CalendarService) GetYearlyCalendar(year int, lat, lon float64, method string) models.YearlyCalendarResponse {
	resp := models.YearlyCalendarResponse{
		Year:   year,
		Method: method,
		Months: make([]models.CalendarMonth, 0, 12),
	}

	// 1. Resolve start of month 1
	estJD := calendar.HijriEpochJD + float64(year-1)*354.367 + 15
	estGreg := calendar.JDToTime(estJD)
	resolvedStart := c.findHijriMonthStart(method, year, 1, estGreg, lat, lon)

	for m := 1; m <= 12; m++ {
		// 2. Resolve start of NEXT month
		nextM, nextY := m+1, year
		if nextM > 12 {
			nextM = 1
			nextY++
		}

		estNextJD := calendar.HijriEpochJD + float64(year-1)*354.367 + float64(m)*29.53 + 15
		estNextGreg := calendar.JDToTime(estNextJD)
		resolvedNextStart := c.findHijriMonthStart(method, nextY, nextM, estNextGreg, lat, lon)

		// 3. Total days between this month's start and next month's start
		totalDays := int(math.Round(resolvedNextStart.Sub(resolvedStart).Hours() / 24.0))

		resp.Months = append(resp.Months, models.CalendarMonth{
			MonthID:        m,
			MonthName:      calendar.MonthNames[m-1],
			TotalDays:      totalDays,
			Day1Weekday:    int(resolvedStart.Weekday()),
			StartGregorian: resolvedStart.Format("2006-01-02"),
		})

		// 4. Current next is the start for the following month in the loop
		resolvedStart = resolvedNextStart
	}

	return resp
}

// findHijriMonthStart searches around a Gregorian estimate for the Gregorian date
// that marks Day 1 of Hijri year/month under a specific method.
func (c *CalendarService) findHijriMonthStart(method string, hYear, hMonth int, estGreg time.Time, lat, lon float64) time.Time {
	current := estGreg.AddDate(0, 0, -5) // Start 5 days before estimate
	for i := 0; i < 10; i++ {
		res := c.DateSvc.ResolveDynamicHijriDate(method, current, lat, lon)
		if res.Year == hYear && res.Month == hMonth && res.Day == 1 {
			return time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)
		}
		if res.Year > hYear || (res.Year == hYear && res.Month > hMonth) {
			// Went too far
			current = current.AddDate(0, 0, -1)
			continue
		}
		current = current.AddDate(0, 0, 1)
	}
	// Fallback to tabular day 1 if search fails
	return time.Date(estGreg.Year(), estGreg.Month(), estGreg.Day(), 0, 0, 0, 0, time.UTC)
}
