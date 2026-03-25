package calendar

import (
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
)

type Logic struct {
	Astro   *astronomy.Adapter
	Manager *astronomy.EphemerisManager
}

func NewLogic(astro *astronomy.Adapter, manager *astronomy.EphemerisManager) *Logic {
	return &Logic{
		Astro:   astro,
		Manager: manager,
	}
}
