package calendar

import (
	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"math"
	"time"
)

func (l *Logic) FindIjtima(approxDate time.Time) (time.Time, error) {
	// Cari dalam rentang 30 hari (biar pasti ketemu 1 fase konjungsi)
	startET, _ := astronomy.Str2et(approxDate.AddDate(0, 0, -15).Format(astronomy.TimeFormat))
	endET, _ := astronomy.Str2et(approxDate.AddDate(0, 0, 15).Format(astronomy.TimeFormat))

	low, high := startET, endET
	var mid float64

	// Gunakan 60 iterasi biar makin sakti
	for i := 0; i < 60; i++ {
		mid = (low + high) / 2
		if l.getLongitudeDiff(mid) > 0 {
			high = mid
		} else {
			low = mid
		}
	}
	return astronomy.Et2Utc(mid), nil
}

// FindPreviousIjtima: Nyari ijtima terakhir sebelum tanggal target
func (l *Logic) FindPreviousIjtima(t time.Time) (time.Time, error) {
	// Mundur 30 hari buat nyari konjungsi sebelumnya
	return l.FindIjtima(t.AddDate(0, 0, -29))
}

func (l *Logic) getLongitudeDiff(et float64) float64 {
	sunPos, _ := l.Manager.GetGeocentric(astronomy.Sun, et, astronomy.FrameEclipJ2000)
	moonPos, _ := l.Manager.GetGeocentric(astronomy.Moon, et, astronomy.FrameEclipJ2000)

	sunLong := math.Atan2(sunPos.Y, sunPos.X)
	moonLong := math.Atan2(moonPos.Y, moonPos.X)

	diff := moonLong - sunLong
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}

	return diff
}
