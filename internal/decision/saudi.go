package decision

type Saudi struct{}

func (s Saudi) IsVisible(ctx Context) bool {
	// Kriteria: ijtimak sebelum maghrib dan saat maghrib bulan berada di atas ufuk
	// Yang berarti MoonsetTime > SunsetTime
	return ctx.IjtimaTime.Before(ctx.SunsetTime) && ctx.MoonsetTime.After(ctx.SunsetTime)
}

func (s Saudi) Name() string {
	return "SAUDI"
}
