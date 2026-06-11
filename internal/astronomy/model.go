package astronomy

// CelestialBody represents a standard celestial object structure
type CelestialBody struct {
	Name   string
	Radius float64
}

var (
	Earth = CelestialBody{Name: "EARTH", Radius: 6378.137}
	Moon  = CelestialBody{Name: "MOON", Radius: 1737.4}
	Sun   = CelestialBody{Name: "SUN", Radius: 696340.0}
)
