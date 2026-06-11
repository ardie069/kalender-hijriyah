package highlatitude

import "time"

type Seventh struct{}

func (s Seventh) Adjust(p AdjustParams) (time.Time, time.Time) {
	seventh := p.NightDuration / 7
	fajr := p.Sunrise.Add(-seventh)
	isha := p.Maghrib.Add(seventh)
	return fajr, isha
}

func (s Seventh) Name() string {
	return "SEVENTH_OF_NIGHT"
}
