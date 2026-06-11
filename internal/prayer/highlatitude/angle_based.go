package highlatitude

import (
	"math"
	"time"
)

type AngleBased struct{}

func (a AngleBased) Adjust(p AdjustParams) (time.Time, time.Time) {
	isha := p.Isha

	fajrPortion := math.Abs(p.FajrAngle) / 60.0
	fajr := p.Sunrise.Add(-time.Duration(float64(p.NightDuration) * fajrPortion))

	if p.IshaOffsetMin <= 0 {
		ishaPortion := math.Abs(p.IshaAngle) / 60.0
		isha = p.Maghrib.Add(time.Duration(float64(p.NightDuration) * ishaPortion))
	}

	return fajr, isha
}

func (a AngleBased) Name() string {
	return "ANGLE_BASED"
}
