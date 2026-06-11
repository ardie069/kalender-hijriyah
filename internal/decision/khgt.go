package decision

type KHGT struct{}

func (k KHGT) IsVisible(ctx Context) bool {
	// Syarat KHGT (Kongres Istanbul 2016):
	// Altitude >= 5 derajat dan Elongasi >= 8 derajat (Geosentris)
	return ctx.Altitude >= 5.0 && ctx.Elongation >= 8.0
}

func (k KHGT) Name() string {
	return "KHGT"
}
