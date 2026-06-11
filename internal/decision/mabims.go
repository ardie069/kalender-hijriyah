package decision

type MABIMS struct{}

func (m MABIMS) IsVisible(ctx Context) bool {
	// Syarat Alt >= 3 derajat & Elong >= 6.4 derajat
	return ctx.Altitude >= 3.0 && ctx.Elongation >= 6.4
}

func (m MABIMS) Name() string {
	return "MABIMS"
}
