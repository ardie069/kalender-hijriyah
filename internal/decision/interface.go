package decision

import "time"

type Context struct {
	Altitude       float64
	Elongation     float64
	IjtimaTime     time.Time
	SunsetTime     time.Time
	MoonsetTime    time.Time
	IjtimaHappened bool
}

type Method interface {
	IsVisible(ctx Context) bool
	Name() string
}
