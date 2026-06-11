package madhab

type Shafii struct{}

func (s Shafii) ShadowFactor() int {
	return 1
}

func (s Shafii) Name() string {
	return "Syafi'i"
}
