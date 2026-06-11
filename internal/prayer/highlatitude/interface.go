package highlatitude

import "time"

type AdjustParams struct {
	Fajr          time.Time
	Isha          time.Time
	Sunrise       time.Time
	Maghrib       time.Time
	NightDuration time.Duration
	FajrAngle     float64
	IshaAngle     float64
	IshaOffsetMin float64
}

type Method interface {
	Adjust(params AdjustParams) (time.Time, time.Time)
	Name() string
}
