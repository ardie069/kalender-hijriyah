package calendar

import (
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"math"
	"time"
)

func (l *Logic) FindIjtima(approxDate time.Time) (time.Time, error) {
	// Cari dalam rentang 30 hari (biar pasti ketemu 1 fase konjungsi)
	startET, err := astronomy.Str2et(approxDate.AddDate(0, 0, -15).Format(astronomy.TimeFormat))
	if err != nil {
		return time.Time{}, err
	}
	endET, err := astronomy.Str2et(approxDate.AddDate(0, 0, 15).Format(astronomy.TimeFormat))
	if err != nil {
		return time.Time{}, err
	}

	low, high := startET, endET
	var mid float64

	// Gunakan 60 iterasi biar makin sakti
	for range 60 {
		mid = (low + high) / 2
		diff, err := l.getLongitudeDiff(mid)
		if err != nil {
			return time.Time{}, err
		}
		if diff > 0 {
			high = mid
		} else {
			low = mid
		}
	}
	return astronomy.Et2Utc(mid), nil
}

// FindPreviousIjtima: Nyari ijtima terakhir sebelum atau sekitar tanggal target
func (l *Logic) FindPreviousIjtima(t time.Time) (time.Time, error) {
	// Mundur 16 hari buat nyari konjungsi yang mendasari bulan ini.
	// FindIjtima nyari +/- 15 hari dari titik ini, jadi total nyari di [t-31, t-1].
	// Ini memastikan ijtima yang baru terjadi "hari ini" tidak terambil sebagai ijtima bulan lalu.
	return l.FindIjtima(t.AddDate(0, 0, -16))
}

func (l *Logic) getLongitudeDiff(et float64) (float64, error) {
	sunPos, err := l.Manager.GetGeocentric(astronomy.Sun, et, astronomy.FrameEclipJ2000)
	if err != nil {
		return 0, err
	}
	moonPos, err := l.Manager.GetGeocentric(astronomy.Moon, et, astronomy.FrameEclipJ2000)
	if err != nil {
		return 0, err
	}

	sunLong := math.Atan2(sunPos.Y, sunPos.X)
	moonLong := math.Atan2(moonPos.Y, moonPos.X)

	diff := moonLong - sunLong
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}

	return diff, nil
}
