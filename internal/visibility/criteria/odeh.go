package criteria

type Odeh struct{}

func (o Odeh) Evaluate(ctx Context) Category {
	w := ctx.Width
	v := ctx.ArcOfLight - (7.1651 - 6.3226*w + 0.7319*w*w - 0.1018*w*w*w)

	switch {
	case v >= 5.65:
		return CategoryA
	case v >= 2.0:
		return CategoryB
	case v >= -0.96:
		return CategoryC
	default:
		return CategoryD // Impossible
	}
}

func (o Odeh) Name() string {
	return "Odeh"
}