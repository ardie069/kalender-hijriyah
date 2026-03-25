package services

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// ResolveMonthStart mencari kapan tepatnya tanggal 1 bulan Hijriyah ini menurut suatu metode (Global / Lokal)
// dan mengembalikan struktur `models.HijriDate` untuk target hari ini.
func (s *HijriService) ResolveDynamicHijriDate(m string, targetDay time.Time, lat, lon float64) models.HijriDate {
	// Pijakan dasar: Gunakan tabular untuk tahu kita di bulan dan tahun berapa secara kasar.
	// Tabular cukup akurat dalam margin +/- 2 hari.
	stableNoon := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 12, 0, 0, 0, time.UTC)
	baseH := calendar.GetTabularHijri(stableNoon)

	// 1. Tentukan Ijtima yang mendasari bulan ini.
	// FindPreviousIjtima akan mencari dalam rentang [targetDay-30, targetDay].
	prevIjtima, _ := s.Cal.FindPreviousIjtima(targetDay)

	// 1. Kita ambil tanggal dasar Ijtima (H-0)
	sunsetCheckDate := time.Date(prevIjtima.Year(), prevIjtima.Month(), prevIjtima.Day(), 12, 0, 0, 0, time.UTC)
	var sunsetCheck time.Time
	var isNewMonth bool

	// 2. Tentukan Tanggal Evaluasi (sunsetCheckDate) berdasarkan Kriteria
	// Jika kriteria lokal, kita cek pakai sunset lokal.
	// Jika UMM_AL_QURA, kita cek pakai sunset Mekkah.
	// Jika MABIMS, kita cek pakai sunset Sabang (untuk konsistensi regional).
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

	case "MABIMS":
		sabangLat, sabangLon := 5.89, 95.32
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		if sunsetCheck.Before(prevIjtima) {
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, sabangLat, sabangLon)
		}
		moonsetCheck, _ := s.Astro.GetMoonset(sunsetCheck, sabangLat, sabangLon)
		telCheck, _ := s.Astro.GetMoonTelemetry(sunsetCheck, sabangLat, sabangLon)
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck, sunsetCheck, moonsetCheck, prevIjtima)
		isNewMonth = resLocal.IsNewMonth

	case "KHGT":
		// KHGT tidak bergantung pada satu sunset spesifik, melainkan global.
		// Evaluasi dimulai dari hari Ijtima.
		isNewMonth, _ = s.Cal.ScanGlobalKHGT(sunsetCheckDate, prevIjtima)
		// Jika hari Ijtima gagal, konvensi kalender global adalah menggenapkan bulan sebelumnya (30 hari).
		// Evaluasi KHGT secara spesifik sudah mencakup seluruh bumi untuk hari tersebut.
		if !isNewMonth {
			// Coba evaluasi hari berikutnya
			sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
			isNewMonth, _ = s.Cal.ScanGlobalKHGT(sunsetCheckDate, prevIjtima)
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
		resLocal := s.Cal.EvaluateLocalHisab(m, sunsetCheck, telCheck, sunsetCheck, moonsetCheck, prevIjtima)
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

	daysElapsed := int(math.Round(targetDayNoon.Sub(monthStartNoon).Hours() / 24.0))

	// Jika targetDay sama dengan monthStartDate, itu adalah tanggal 1.
	// Jika targetDay adalah besoknya monthStartDate, itu adalah tanggal 2.
	// Jadi H = (Target - Start) + 1
	hDay := daysElapsed + 1

	evalMonth := baseH.Month
	evalYear := baseH.Year

	// Jika hDay < 1, berarti kita masih di bulan sebelumnya
	if hDay < 1 {
		evalMonth--
		if evalMonth < 1 {
			evalMonth = 12
			evalYear--
		}
		// SANGAT PENTING: Jika bulan sebelumnya adalah 29 atau 30 hari,
		// kita asumsikan 30 saja sebagai fallback jika tidak menghitung ijtima sebelumnya lagi.
		hDay += 30
	} else if hDay > 30 {
		// Jika lebih dari 30 hari, berarti sudah masuk bulan berikutnya
		hDay -= 30
		evalMonth++
		if evalMonth > 12 {
			evalMonth = 1
			evalYear++
		}
	} else if hDay > 29 && !isNewMonth {
		// Kasus khusus: Jika ijtima belum tercapai tapi tabular sudah ganti bulan,
		// hDay bisa jadi 29 atau 30 di bulan sebelumnya.
		// Tabular Syawal (10) 1, tapi Rukyat masih Ramadan (9) 30.
		if baseH.Day == 1 {
			evalMonth--
			if evalMonth < 1 {
				evalMonth = 12
				evalYear--
			}
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
