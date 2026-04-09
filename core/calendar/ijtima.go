package calendar

import (
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"math"
	"time"
)

func (l *Logic) FindIjtima(approxDate time.Time) (time.Time, error) {
	// 1. Get current longitude difference as an estimate
	et, _ := astronomy.Str2et(approxDate.Format(astronomy.TimeFormat))
	diff, err := l.getLongitudeDiff(et)
	if err != nil {
		return time.Time{}, err
	}

	// Moon moves ~12.2 degrees per day relative to Sun.
	// Estimated days since/until ijtima: diff (radians) / (12.2 * PI / 180)
	daysDiff := diff / (12.19 * math.Pi / 180.0)
	estIjtimaET := et - (daysDiff * 86400.0)

	// 2. Search in a narrow 4-day window around the estimate
	low := estIjtimaET - (2 * 86400.0)
	high := estIjtimaET + (2 * 86400.0)
	var mid float64

	for range 60 {
		mid = (low + high) / 2
		d, err := l.getLongitudeDiff(mid)
		if err != nil {
			return time.Time{}, err
		}
		if d > 0 {
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
