package ummalqura

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy/ephemeris"
	"github.com/ardie069/kalender-hijriyah/internal/decision"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Service struct {
	Astro *cspice.Adapter
	Ephem *ephemeris.Service
}

func NewService(astro *cspice.Adapter, ephem *ephemeris.Service) *Service {
	return &Service{
		Astro: astro,
		Ephem: ephem,
	}
}

// EvaluateHisabMecca evaluates whether the crescent is visible based on the Hisab Mecca (Umm Al Qura) method
func (s *Service) EvaluateHisabMecca(sunsetCheckDate time.Time, ijtima time.Time) bool {
	meccaLat, meccaLon := 21.4225, 39.8262
	sunsetCheck, _ := s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
	
	if sunsetCheck.Before(ijtima) {
		sunsetCheckDate = sunsetCheckDate.AddDate(0, 0, 1)
		sunsetCheck, _ = s.Astro.GetSunset(sunsetCheckDate, meccaLat, meccaLon)
	}
	moonsetMecca, _ := s.Astro.GetMoonset(sunsetCheck, meccaLat, meccaLon)

	ummAlQura := decision.UmmAlQura{}
	ctx := decision.Context{
		IjtimaTime:  ijtima,
		SunsetTime:  sunsetCheck,
		MoonsetTime: moonsetMecca,
	}
	return ummAlQura.IsVisible(ctx)
}
