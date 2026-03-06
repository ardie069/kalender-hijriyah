package services

import (
	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"time"
)

type LocationInfo struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type MethodResult struct {
	HijriDate         calendar.HijriDate `json:"hijri_date"`
	CurrentAltitude   float64            `json:"current_altitude"`
	CurrentElongation float64            `json:"current_elongation"`
	Prediction        HilalPrediction    `json:"prediction"`
}

type HilalPrediction struct {
	CheckDateUTC time.Time `json:"check_date_utc"`
	IsNewMonth   bool      `json:"is_new_month"`
	IjtimaTime   time.Time `json:"ijtima_time"`
	Altitude     float64   `json:"altitude"`
	Elongation   float64   `json:"elongation"`
	AgeHours     float64   `json:"age_hours"`
}

type HijriResponse struct {
	GregorianDate time.Time               `json:"gregorian_date"`
	Location      LocationInfo            `json:"location"`
	Methods       map[string]MethodResult `json:"methods"`
}

type HijriService struct {
	Astro *astronomy.Adapter
	Cal   *calendar.Logic
}

func (s *HijriService) GetFullCalendarInfo(t time.Time, lat, lon float64) HijriResponse {
	tUTC := t.UTC()

	// 1. Inisialisasi Data Dasar
	sunsetToday, _ := s.Astro.GetSunset(tUTC, lat, lon)
	currentH := calendar.GetTabularHijri(tUTC)
	realtimeTel, _ := s.Astro.GetMoonTelemetry(tUTC, lat, lon)

	// 2. Tentukan Tanggal Display (Logic Maghrib)
	isAfterSunset := tUTC.After(sunsetToday)
	displayDay := currentH.Day
	if isAfterSunset {
		displayDay++
	}

	// 3. LOGIKA PREDIKSI: Cari kapan Sunset di hari ke-29 bulan ini
	// Caranya: Kita cari Ijtima berikutnya, lalu ambil sunset tepat sebelumnya.
	nextIjtima, _ := s.Cal.FindIjtima(tUTC)

	searchDate := time.Date(nextIjtima.Year(), nextIjtima.Month(), nextIjtima.Day(), 10, 0, 0, 0, time.UTC)

	sunset29, _ := s.Astro.GetSunset(searchDate, lat, lon)

	moonset29, _ := s.Astro.GetMoonset(sunset29, lat, lon)
	tel29, _ := s.Astro.GetMoonTelemetry(sunset29, lat, lon)

	ageHours := sunset29.Sub(nextIjtima).Hours()

	resp := HijriResponse{
		GregorianDate: tUTC,
		Location:      LocationInfo{Lat: lat, Lon: lon},
		Methods:       make(map[string]MethodResult),
	}

	methodList := []string{"MABIMS", "WUJUDUL_HILAL", "UMM_AL_QURA", "TURKEY_2016", "UGHC_KHGT"}

	for _, m := range methodList {
		var result MethodResult

		// A. Data Real-time (Buat tampilan sekarang)
		result.HijriDate = calendar.HijriDate{
			Day:       displayDay,
			Month:     currentH.Month,
			MonthName: currentH.MonthName,
			Year:      currentH.Year,
		}
		result.CurrentAltitude = realtimeTel.Altitude
		result.CurrentElongation = realtimeTel.Elongation

		// B. Data Prediksi (The "Future" Insight)
		pred := HilalPrediction{
			CheckDateUTC: sunset29,
			IjtimaTime:   nextIjtima,
			Altitude:     tel29.Altitude,
			Elongation:   tel29.Elongation,
			AgeHours:     ageHours,
		}

		// C. Jalankan Evaluasi berdasarkan Metode
		switch m {
		case "UGHC_KHGT":
			result.HijriDate.Day = displayDay + 1
			pred.IsNewMonth = s.Cal.ScanGlobalUGHC(sunset29, nextIjtima)

		case "UMM_AL_QURA":
			result.HijriDate.Day = displayDay + 1
			// Umm Al-Qura punya logic sendiri di Makkah
			pred.IsNewMonth = calendar.IsUmmAlQura(nextIjtima, sunset29, moonset29)

		default:
			// Metode Lokal (MABIMS, Wujudul Hilal, dkk)
			resLocal := s.Cal.EvaluateLocalHisab(m, sunset29, tel29.Altitude, tel29.Elongation, sunset29, moonset29, nextIjtima)
			pred.IsNewMonth = resLocal.IsNewMonth
		}

		result.Prediction = pred
		resp.Methods[m] = result
	}

	return resp
}
