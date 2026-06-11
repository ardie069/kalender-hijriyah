package prayer

// Calculation holds the fundamental math for computing prayer times
type Calculation struct {
	BaseLat float64
	BaseLon float64
}

// Compute returns default computed values
func (c *Calculation) Compute() float64 {
	return 0.0
}
