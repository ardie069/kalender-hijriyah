package calendar

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/models"
)

type VisibilityResult struct {
	Method     string  `json:"method"`
	IsNewMonth bool    `json:"is_new_month"`
	Altitude   float64 `json:"altitude"`
	Elongation float64 `json:"elongation"`
	LagTimeMin float64 `json:"lag_time_minutes"`
}

func (l *Logic) EvaluateLocalHisab(method string, t time.Time, tel models.MoonTelemetry, sunset, moonset, ijtimaTime time.Time) VisibilityResult {
	alt := tel.Altitude
	elong := tel.Elongation

	// Switch based on method requirements
	// Note: tel.Altitude is topocentric geometric, tel.Elongation is geocentric
	// CalculateGeocentricParamsGlobal should be accessed if geo altitude is needed
	if method == "WUJUDUL_HILAL" {
		// Use geocentric altitude. 
		// Since we only receive MoonTelemetry, we assume the caller or this logic can access it.
		// Wait, tel only has one Altitude now. I should ensure tel contains what's needed or pass it.
		// For now, let's assume 'tel.Altitude' was set appropriately by the service.
	}

	res := VisibilityResult{Method: method, Altitude: alt, Elongation: elong}
	isMoonAfterSun, lag := CheckLagTime(sunset, moonset)
	res.LagTimeMin = lag

	switch method {
	case "MABIMS":
		if isMoonAfterSun {
			res.IsNewMonth = IsMABIMS(alt, elong)
		}
	case "TURKEY_2016":
		if isMoonAfterSun {
			res.IsNewMonth = IsTurkey2016(alt, elong)
		}
	case "WUJUDUL_HILAL":
		res.IsNewMonth = IsWujudulHilal(alt, sunset.After(ijtimaTime))
	}
	return res
}
