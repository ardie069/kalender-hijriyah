package madhab

type Hanafi struct{}

func (h Hanafi) ShadowFactor() int {
	return 2
}

func (h Hanafi) Name() string {
	return "Hanafi"
}
