package highlatitude

import "time"

type MiddleNight struct{}

func (m MiddleNight) Adjust(p AdjustParams) (time.Time, time.Time) {
	half := p.NightDuration / 2
	fajr := p.Sunrise.Add(-half)
	isha := p.Maghrib.Add(half)
	return fajr, isha
}

func (m MiddleNight) Name() string {
	return "MIDDLE_OF_NIGHT"
}
