package calendar

import (
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Logic struct {
	Astro   *cspice.Adapter
	Manager *cspice.EphemerisManager
}

func NewLogic(astro *cspice.Adapter, manager *cspice.EphemerisManager) *Logic {
	return &Logic{
		Astro:   astro,
		Manager: manager,
	}
}
