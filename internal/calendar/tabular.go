package calendar

import (
	"math"
	"time"
)

func GetTabularHijri(t time.Time) HijriDate {
	// 1. Julian Day
	jd := float64(t.Unix())/86400.0 + 2440587.5

	// 2. Days since Epoch Hijri (1948439.5)
	days := int(math.Floor(jd - 1948439.5))

	// 3. Cycle 30 years (10631 days)
	cycle := days / 10631
	rem := days % 10631

	// 4. Year within cycle
	yearInCycle := int(float64(rem) / 354.366)
	if yearInCycle > 29 {
		yearInCycle = 29
	}

	year := cycle*30 + yearInCycle + 1

	// 5. Day within year
	dayInYear := rem - int(math.Floor(float64(yearInCycle)*354.366+0.5))

	// 6. Find Month and Day
	// Daftar kumulatif hari dalam bulan Hijriyah (Tabular)
	monthDays := []int{0, 30, 59, 89, 118, 148, 177, 207, 236, 266, 295, 325, 355}

	month := 12
	for i := 1; i <= 12; i++ {
		if dayInYear <= monthDays[i] {
			month = i
			break
		}
	}

	day := dayInYear - monthDays[month-1]
	if day <= 0 {
		day = 1
	}

	return HijriDate{
		Day:       day,
		Month:     month,
		MonthName: MonthNames[month-1],
		Year:      year,
	}
}
