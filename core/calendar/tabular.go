package calendar

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/models"
)

var HijriEpoch = time.Date(622, time.July, 16, 0, 0, 0, 0, time.UTC)
var LeapYears = map[int]bool{
	2: true, 5: true, 7: true, 10: true, 13: true,
	16: true, 18: true, 21: true, 24: true, 26: true, 29: true,
}
var MonthNames = []string{
	"Muharram", "Safar", "Rabi'ul Awwal", "Rabi'ul Akhir",
	"Jumadil Awwal", "Jumadil Akhir", "Rajab", "Sya'ban",
	"Ramadan", "Syawal", "Zulqa'dah", "Zulhijjah",
}

func GetTabularHijri(t time.Time) models.HijriDate {
	// 1. Konversi ke Julian Day (JD)
	// 2440587.5 adalah JD saat 1 Jan 1970 00:00 UTC
	jd := float64(t.Unix())/86400.0 + 2440587.5

	// 2. Days since Epoch Hijri (1 Muharram 1 H = JD 1948439.5)
	days := int(math.Floor(jd - 1948439.5))

	// 3. Siklus 30 Tahun (10.631 hari)
	cycle := days / 10631
	rem := days % 10631

	// 4. Hitung Tahun dalam Siklus
	year := cycle * 30
	for i := 1; i <= 30; i++ {
		daysInYear := 354
		if LeapYears[i] {
			daysInYear = 355
		}
		if rem < daysInYear {
			year += i
			break
		}
		rem -= daysInYear
	}

	// 5. Hitung Bulan dan Hari (rem sekarang adalah dayInYear)
	month := 0
	day := 0
	for m := 1; m <= 12; m++ {
		daysInMonth := 29
		if m%2 != 0 { // Bulan Ganjil = 30 hari
			daysInMonth = 30
		}
		// Zulhijjah (Bulan 12) di tahun kabisat = 30 hari
		if m == 12 && LeapYears[year%30] {
			daysInMonth = 30
		}

		if rem < daysInMonth {
			month = m
			day = rem + 1
			break
		}
		rem -= daysInMonth
	}

	return models.HijriDate{
		Day:       day,
		Month:     month,
		MonthName: MonthNames[month-1],
		Year:      year,
	}
}
