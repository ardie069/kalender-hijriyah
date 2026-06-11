package hijri

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy/ephemeris"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/calendar/ummalqura"
	"github.com/ardie069/kalender-hijriyah/internal/decision"
	"github.com/ardie069/kalender-hijriyah/internal/models"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/timezone"
	"github.com/ardie069/kalender-hijriyah/internal/visibility/scan"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type DateService struct {
	Tz    *timezone.Service
	Astro *cspice.Adapter
	Ephem     *ephemeris.Service
	Scan      *scan.Scanner
	UmmAlQura *ummalqura.Service
}

func NewDateService(astro *cspice.Adapter, ephem *ephemeris.Service, s *scan.Scanner, umm *ummalqura.Service, tz *timezone.Service) *DateService {
	return &DateService{
		Tz:        tz,
		Astro:     astro,
		Ephem:     ephem,
		Scan:      s,
		UmmAlQura: umm,
	}
}

type CalendarService struct {
	DateSvc *DateService
}

func NewCalendarService(dateSvc *DateService) *CalendarService {
	return &CalendarService{DateSvc: dateSvc}
}

func (s *DateService) GetHijriTargetDate(t time.Time, lat, lon float64) time.Time {
	tUTC := t.UTC()

	sunsetToday, err := s.Astro.GetSunset(tUTC, lat, lon)
	if err != nil {
		sunsetToday = tUTC
	}

	hToday := calendar.GetTabularHijri(tUTC)

	rolloverTime := sunsetToday
	if (hToday.Month == 8 || hToday.Month == 9 || hToday.Month == 11) && hToday.Day == 29 {
		ishaToday, err := s.Astro.GetIsha(tUTC, lat, lon)
		if err == nil {
			rolloverTime = ishaToday
		}
	}

	if tUTC.After(rolloverTime) {
		return tUTC.Add(24 * time.Hour)
	}
	return tUTC
}

func (s *DateService) ResolveDynamicHijriDate(m string, targetDay time.Time, lat, lon float64) models.HijriDate {
	targetNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)

	ijtimaRecent, _ := s.Ephem.FindIjtima(targetNoon)
	monthStartRecent := s.evalMonthStart(m, ijtimaRecent, lat, lon)

	var finalMonthStart time.Time

	if targetNoon.Before(monthStartRecent) {
		ijtimaPrev, _ := s.Ephem.FindIjtima(ijtimaRecent.AddDate(0, 0, -29))
		finalMonthStart = s.evalMonthStart(m, ijtimaPrev, lat, lon)
	} else {
		finalMonthStart = monthStartRecent
	}

	daysElapsed := int(math.Round(targetNoon.Sub(finalMonthStart).Hours() / 24.0))
	hDay := daysElapsed + 1

	monthStartJD := float64(finalMonthStart.Unix())/86400.0 + 2440587.5
	monthsSinceEpoch := int(math.Round((monthStartJD - 1948439.5) / 29.53059))

	evalYear := (monthsSinceEpoch / 12) + 1
	evalMonth := (monthsSinceEpoch % 12) + 1

	if hDay > 30 {
		hDay -= 30
		evalMonth++
		if evalMonth > 12 {
			evalMonth = 1
			evalYear++
		}
	}

	return models.HijriDate{
		Day:       hDay,
		Month:     evalMonth,
		MonthName: calendar.MonthNames[evalMonth-1],
		Year:      evalYear,
		IsTabular: false,
	}
}

func (s *DateService) evalMonthStart(m string, ijtima time.Time, lat, lon float64) time.Time {
	sunsetCheckDate := time.Date(ijtima.Year(), ijtima.Month(), ijtima.Day(), 12, 0, 0, 0, time.UTC)
	var isNewMonth bool

	switch m {
	case "UMM_AL_QURA":
		hisabMecca := s.UmmAlQura.EvaluateHisabMecca(sunsetCheckDate, ijtima)

		// Evaluasi Rukyat Saudi (Sudair & Tumair)
		rukyatSaudi := s.evaluateSaudiRukyat(sunsetCheckDate, ijtima)

		// Tentukan bulan target
		hDate := calendar.GetTabularHijri(ijtima.AddDate(0, 0, 15))
		isCriticalMonth := (hDate.Month == 9 || hDate.Month == 10 || hDate.Month == 12)

		if isCriticalMonth {
			// Ramadan, Syawal, Zulhijjah WAJIB pakai rukyat
			isNewMonth = rukyatSaudi
		} else {
			// Bulan lain pakai Hisab Umm Al Qura
			isNewMonth = hisabMecca
			// Koreksi otomatis jika Hisab tidak sama dengan potensi Rukyat (kasus langka)
			if hisabMecca != rukyatSaudi {
				isNewMonth = rukyatSaudi
			}
		}

	case "SAUDI":
		isNewMonth = s.evaluateSaudiRukyat(sunsetCheckDate, ijtima)

	case "MABIMS":
		sabangLat, sabangLon := 5.89, 95.32
		sunsetCheck, _ := s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		if sunsetCheck.Before(ijtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		}
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, sabangLat, sabangLon)
		
		mabims := decision.MABIMS{}
		ctx := decision.Context{
			Altitude:   telCheck.Altitude,
			Elongation: telCheck.Elongation,
		}
		isNewMonth = mabims.IsVisible(ctx)

	case "KHGT":
		khgtRes := s.Scan.ScanGlobalKHGT(sunsetCheckDate, ijtima)
		isNewMonth = khgtRes.IsGlobalValid
		if !isNewMonth {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			khgtRes2 := s.Scan.ScanGlobalKHGT(sunsetCheckDate, ijtima)
			isNewMonth = khgtRes2.IsGlobalValid
		}

	default:
		khgtRes := s.Scan.ScanGlobalKHGT(sunsetCheckDate, ijtima)
		isNewMonth = khgtRes.IsGlobalValid
	}

	if isNewMonth {
		return sunsetCheckDate.AddDate(0, 0, 1)
	}
	return sunsetCheckDate.AddDate(0, 0, 2)
}

func (s *DateService) evaluateSaudiRukyat(sunsetCheckDate time.Time, ijtima time.Time) bool {
	// Sudair (lat: 25.5950, lon: 45.6339)
	sudairLat, sudairLon := 25.5950, 45.6339
	sunsetSudair, _ := s.Astro.GetSunset(sunsetCheckDate, sudairLat, sudairLon)
	if sunsetSudair.Before(ijtima) {
		sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
		sunsetSudair, _ = s.Astro.GetSunset(sunsetCheckDate, sudairLat, sudairLon)
	}
	moonsetSudair, _ := s.Astro.GetMoonset(sunsetSudair, sudairLat, sudairLon)

	// Tumair (lat: 25.7039, lon: 45.8614)
	tumairLat, tumairLon := 25.7039, 45.8614
	sunsetTumair, _ := s.Astro.GetSunset(sunsetCheckDate, tumairLat, tumairLon)
	moonsetTumair, _ := s.Astro.GetMoonset(sunsetTumair, tumairLat, tumairLon)

	saudiRule := decision.Saudi{}
	ctxSudair := decision.Context{
		IjtimaTime:  ijtima,
		SunsetTime:  sunsetSudair,
		MoonsetTime: moonsetSudair,
	}
	ctxTumair := decision.Context{
		IjtimaTime:  ijtima,
		SunsetTime:  sunsetTumair,
		MoonsetTime: moonsetTumair,
	}
	return saudiRule.IsVisible(ctxSudair) || saudiRule.IsVisible(ctxTumair)
}

func (s *DateService) GetTabularOnly(t time.Time) models.MethodResult {
	hDate := calendar.GetTabularHijri(t)
	hDate.IsTabular = true

	return models.MethodResult{
		HijriDate:  hDate,
		Prediction: nil,
	}
}

func (s *DateService) GetGregorianMonthInfo(year, month int, lat, lon float64) models.GregorianMonthResponse {
	resp := models.GregorianMonthResponse{
		Year:  year,
		Month: month,
		Lat:   lat,
		Lon:   lon,
		Days:  make([]models.MonthDayInfo, 0, 31),
	}

	startDate := time.Date(year, time.Month(month), 1, 12, 0, 0, 0, time.UTC)
	for d := startDate; d.Month() == time.Month(month); d = d.AddDate(0, 0, 1) {
		dayInfo := models.MonthDayInfo{
			GregorianDate: d.Format("2006-01-02"),
			Corrections:   make(map[string]models.HijriDate),
		}

		dayInfo.Tabular = calendar.GetTabularHijri(d)
		dayInfo.Tabular.IsTabular = true

		methods := []string{"KHGT", "UMM_AL_QURA", "MABIMS"}
		for _, m := range methods {
			resolved := s.ResolveDynamicHijriDate(m, d, lat, lon)
			dayInfo.Corrections[m] = resolved
		}

		resp.Days = append(resp.Days, dayInfo)
	}

	return resp
}

func (c *CalendarService) GetYearlyCalendar(year int, lat, lon float64, method string) models.YearlyCalendarResponse {
	resp := models.YearlyCalendarResponse{
		Year:   year,
		Method: method,
		Months: make([]models.CalendarMonth, 0, 12),
	}

	estJD := calendar.HijriEpochJD + float64(year-1)*354.367 + 15
	estGreg := calendar.JDToTime(estJD)
	resolvedStart := c.findHijriMonthStart(method, year, 1, estGreg, lat, lon)

	for m := 1; m <= 12; m++ {
		nextM, nextY := m+1, year
		if nextM > 12 {
			nextM = 1
			nextY++
		}

		estNextJD := calendar.HijriEpochJD + float64(year-1)*354.367 + float64(m)*29.53 + 15
		estNextGreg := calendar.JDToTime(estNextJD)
		resolvedNextStart := c.findHijriMonthStart(method, nextY, nextM, estNextGreg, lat, lon)

		totalDays := int(math.Round(resolvedNextStart.Sub(resolvedStart).Hours() / 24.0))

		resp.Months = append(resp.Months, models.CalendarMonth{
			MonthID:        m,
			MonthName:      calendar.MonthNames[m-1],
			TotalDays:      totalDays,
			Day1Weekday:    int(resolvedStart.Weekday()),
			StartGregorian: resolvedStart.Format("2006-01-02"),
		})

		resolvedStart = resolvedNextStart
	}

	return resp
}

func (c *CalendarService) findHijriMonthStart(method string, hYear, hMonth int, estGreg time.Time, lat, lon float64) time.Time {
	current := estGreg.AddDate(0, 0, -5)
	for i := 0; i < 10; i++ {
		res := c.DateSvc.ResolveDynamicHijriDate(method, current, lat, lon)
		if res.Year == hYear && res.Month == hMonth && res.Day == 1 {
			return time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)
		}
		if res.Year > hYear || (res.Year == hYear && res.Month > hMonth) {
			current = current.AddDate(0, 0, -1)
			continue
		}
		current = current.AddDate(0, 0, 1)
	}
	return time.Date(estGreg.Year(), estGreg.Month(), estGreg.Day(), 0, 0, 0, 0, time.UTC)
}
