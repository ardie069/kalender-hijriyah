package criteria

type Yallop struct{}

func (y Yallop) Evaluate(ctx Context) Category {
	L := ctx.Elongation
	h := ctx.Altitude
	q := (h - (-0.014*L*L + 0.1*L + 2.53)) / 10.0

	switch {
	case q > 0.422:
		return CategoryA
	case q > 0.222:
		return CategoryB
	case q > 0.022:
		return CategoryC
	case q > -0.178:
		return CategoryD
	case q > -0.278:
		return CategoryE
	default:
		return CategoryF
	}
}

func (y Yallop) Name() string {
	return "Yallop"
}
