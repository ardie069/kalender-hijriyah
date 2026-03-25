package calendar

import (
	"fmt"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/models"
)

func FormatHijri(h models.HijriDate) string {
	return fmt.Sprintf("%d %s %d", h.Day, h.MonthName, h.Year)
}

// PredictHijriDate: Menghitung transisi tanggal
func PredictHijriDate(t time.Time, isNewMonth bool, currentMonth, currentYear int) models.HijriDate {
	if isNewMonth {
		// Transisi ke bulan baru (Tanggal 1)
		newMonth := (currentMonth % 12) + 1
		newYear := currentYear
		if newMonth == 1 {
			newYear++
		}
		return models.HijriDate{
			Day:       1,
			Month:     newMonth,
			MonthName: MonthNames[newMonth-1],
			Year:      newYear,
		}
	}

	// Gagal ganti bulan -> Bulan sebelumnya digenapkan (Istikmal 30 hari)
	// Kita pake logika: prevMonth = currentMonth
	// Karena kita sedang berada di hari ke-29/30 bulan tersebut.
	return models.HijriDate{
		Day:       30,
		Month:     currentMonth,
		MonthName: MonthNames[currentMonth-1],
		Year:      currentYear,
	}
}
