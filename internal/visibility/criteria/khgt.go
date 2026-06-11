package criteria

type KHGT struct{}

func (k KHGT) Evaluate(ctx Context) Category {
	// Kriteria KHGT (Kongres Istanbul 2016) menggunakan tolak ukur Geosentris:
	// Elongasi >= 8.0 derajat dan Altitude >= 5.0 derajat
	if ctx.Altitude >= 5.0 && ctx.Elongation >= 8.0 {
		return Category("KHGT_YES")
	}
	return Category("KHGT_NO")
}

func (k KHGT) Name() string {
	return "KHGT"
}
