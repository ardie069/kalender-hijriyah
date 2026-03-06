package calendar

import (
	"fmt"
	"time"
)

var MonthNames = []string{
	"Muharram", "Shafar", "Rabi'ul Awwal", "Rabi'ul Thani",
	"Jumadil Ula", "Jumadil Thani", "Rajab", "Sha'ban",
	"Ramadan", "Shawwal", "Dhul Qa'dah", "Dhul Hijjah",
}

type HijriDate struct {
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	MonthName string `json:"month_name"`
	Year      int    `json:"year"`
}

func (h HijriDate) String() string {
	return fmt.Sprintf("%d %s %d", h.Day, h.MonthName, h.Year)
}

// PredictHijriDate: Menghitung transisi tanggal
func PredictHijriDate(t time.Time, isNewMonth bool, currentMonth, currentYear int) HijriDate {
	if isNewMonth {
		// Transisi ke bulan baru (Tanggal 1)
		newMonth := (currentMonth % 12) + 1
		newYear := currentYear
		if newMonth == 1 {
			newYear++
		}
		return HijriDate{
			Day:       1,
			Month:     newMonth,
			MonthName: MonthNames[newMonth-1],
			Year:      newYear,
		}
	}

	// Gagal ganti bulan -> Bulan sebelumnya digenapkan (Istikmal 30 hari)
	// Kita pake logika: prevMonth = currentMonth
	// Karena kita sedang berada di hari ke-29/30 bulan tersebut.
	return HijriDate{
		Day:       30,
		Month:     currentMonth,
		MonthName: MonthNames[currentMonth-1],
		Year:      currentYear,
	}
}
