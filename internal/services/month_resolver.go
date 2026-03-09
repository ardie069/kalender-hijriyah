package services

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

// ResolveMonthStart mencari kapan tepatnya tanggal 1 bulan Hijriyah ini menurut suatu metode (Global / Lokal)
// dan mengembalikan struktur `models.HijriDate` untuk target hari ini.
func (s *HijriService) ResolveDynamicHijriDate(m string, targetDay time.Time, lat, lon float64) models.HijriDate {
	// Pijakan dasar: Gunakan tabular untuk tahu kita di bulan dan tahun berapa secara kasar.
	// Tabular cukup akurat dalam margin +/- 2 hari.
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	baseH := calendar.GetTabularHijri(stableNoon)

	prevIjtima, _ := s.Cal.FindPreviousIjtima(targetDay)

	// 1. Kita ambil tanggal dasar Ijtima (H-0)
	sunsetCheckDate := time.Date(prevIjtima.Year(), prevIjtima.Month(), prevIjtima.Day(), 12, 0, 0, 0, time.UTC)
	var sunsetCheck time.Time
	var isNewMonth bool

	// 2. Tentukan Tanggal Evaluasi (sunsetCheckDate) berdasarkan Kriteria
	// Jika kriteria lokal, kita cek pakai sunset lokal.
	// Jika UMM_AL_QURA, kita cek pakai sunset Mekkah.
	switch m {
	case "UMM_AL_QURA":
		meccaLat, meccaLon := 21.4225, 39.8262
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
		}
		moonsetMecca, _ := s.Astro.GetMoonset(sunsetCheck, meccaLat, meccaLon)
		isNewMonth = calendar.IsUmmAlQura(prevIjtima, sunsetCheck, moonsetMecca)

	case "UGHC_KHGT":
		// UGHC tidak bergantung pada satu sunset spesifik, melainkan global.
		// Evaluasi dimulai dari hari Ijtima.
		isNewMonth, _ = s.Cal.ScanGlobalUGHC(sunsetCheckDate, prevIjtima)
		// Jika hari Ijtima gagal, konvensi kalender global adalah menggenapkan bulan sebelumnya (30 hari).
		// Evaluasi UGHC secara spesifik sudah mencakup seluruh bumi untuk hari tersebut.
		if !isNewMonth {
			// Coba evaluasi hari berikutnya
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			isNewMonth, _ = s.Cal.ScanGlobalUGHC(sunsetCheckDate, prevIjtima)
		}

	default:
		// Metode Lokal (MABIMS, Wujudul Hilal, dll)
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, lat, lon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, lat, lon)
		}
		moonsetCheck, _ := s.Astro.GetMoonset(sunsetCheck, lat, lon)
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, lat, lon)
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck.Altitude, telCheck.Elongation, sunsetCheck, moonsetCheck, prevIjtima)
		isNewMonth = resLocal.IsNewMonth
	}

	// 3. Konklusi: Kapan tanggal 1-nya?
	// Jika kriteria terpenuhi pada sunsetCheckDate, maka besok siangnya adalah Tanggal 1.
	// Jika ternyata gagal juga (sangat jarang jika sudah dievaluasi +1 hari), maka lusa adalah Tanggal 1.
	var monthStartDate time.Time
	if isNewMonth {
		monthStartDate = sunsetCheckDate.AddDate(0, 0, 1)
	} else {
		monthStartDate = sunsetCheckDate.AddDate(0, 0, 2)
	}

	// 4. Hitung Days Elapsed
	// Berapa hari dari `monthStartDate` ke `targetDay`?
	targetDayNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	monthStartNoon := time.Date(monthStartDate.Year(), monthStartDate.Month(), monthStartDate.Day(), 12, 0, 0, 0, time.UTC)

	daysElapsed := int(targetDayNoon.Sub(monthStartNoon).Hours()/24.0) + 1

	// Kita bisa gunakan baseH (Tabular) untuk menyesuaikan Month/Year yang benar,
	// Jika elapsed < 1, berarti kita masih di bulan SEBELUMNYA.
	evalMonth := baseH.Month
	evalYear := baseH.Year

	// Perbaikan overflow/underflow untuk Day (Kritis!)
	if daysElapsed < 1 {
		evalMonth--
		if evalMonth < 1 {
			evalMonth = 12
			evalYear--
		}
		// Idealnya umur bulan Hijri sebelumnya adalah 29/30 bergantung ijtima yg lebih lama lagi,
		// tapi pendekatan kompromi sementara (fallback) agar logika nggak rekursif tanpa batas berhentinya:
		// Diasumsikan sisa hari mundur dari 30.
		daysElapsed += 30
	} else if daysElapsed > 30 {
		daysElapsed -= 30
		evalMonth++
		if evalMonth > 12 {
			evalMonth = 1
			evalYear++
		}
	}

	return models.HijriDate{
		Day:       daysElapsed,
		Month:     evalMonth,
		MonthName: calendar.MonthNames[evalMonth-1],
		Year:      evalYear,
		IsTabular: false,
	}
}
