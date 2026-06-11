package decision

type KHGT struct{}

func (k KHGT) IsVisible(ctx Context) bool {
	// Syarat Ijtima qoblal sunset & Alt > 0 derajat
	return ctx.IjtimaHappened && ctx.Altitude > 0.0
}

func (k KHGT) Name() string {
	return "KHGT"
}
