package decision

type UmmAlQura struct{}

func (s UmmAlQura) IsVisible(ctx Context) bool {
	// Syarat Ijtima qoblal ghurub & Moonset ba'da ghurub
	return ctx.IjtimaTime.Before(ctx.SunsetTime) && ctx.MoonsetTime.After(ctx.SunsetTime)
}

func (s UmmAlQura) Name() string {
	return "UMM_AL_QURA"
}
