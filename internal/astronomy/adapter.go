package astronomy

import (
	"math"
	"time"
)

type MoonTelemetry struct {
	Altitude     float64   `json:"altitude"`
	Azimuth      float64   `json:"azimuth"`
	Elongation   float64   `json:"elongation"`
	Illumination float64   `json:"illumination"`
	DistanceKM   float64   `json:"distance_km"`
	PhaseName    string    `json:"phase_name"`
	Timestamp    time.Time `json:"timestamp"`
}

func GetAdapter(manager *EphemerisManager) *Adapter {
	return &Adapter{
		Manager: manager,
	}
}

// getPhaseName: Nentuin nama fase berdasarkan persen cahaya dan status waxing/waning
func getPhaseName(illum float64, isWaxing bool) string {
	switch {
	case illum < 0.1:
		return "New Moon"
	case illum >= 0.1 && illum < 50:
		if isWaxing {
			return "Waxing Crescent"
		}
		return "Waning Crescent"
	case illum >= 50 && illum < 51:
		if isWaxing {
			return "First Quarter"
		}
		return "Third Quarter"
	case illum >= 51 && illum < 99.9:
		if isWaxing {
			return "Waxing Gibbous"
		}
		return "Waning Gibbous"
	case illum >= 99.9:
		return "Full Moon"
	default:
		return "Unknown"
	}
}

func (a *Adapter) GetMoonTelemetry(dt time.Time, latitude, longitude float64) (MoonTelemetry, error) {
	et, _ := Str2et(dt.UTC().Format(TimeFormat))

	// 1. Get Positions
	sunPos, _ := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	moonPos, _ := a.Manager.GetGeocentric(Moon, et, FrameJ2000)
	// Earth is at (0,0,0) in Geocentric mode, so moonPos is already Earth->Moon vector

	// 2. Elongation (Angle Earth-Sun vs Earth-Moon)
	elongation := math.Acos(sunPos.Unit().Dot(moonPos.Unit())) * (180.0 / math.Pi)

	// 3. Illumination & Phase (Angle Sun-Moon-Earth)
	sunToMoon := sunPos.Sub(moonPos).Unit()
	earthToMoon := moonPos.Scale(-1).Unit() // Earth (0,0,0) minus MoonPos
	phaseAngle := math.Acos(sunToMoon.Dot(earthToMoon))
	illumination := (1.0 + math.Cos(phaseAngle)) / 2.0 * 100.0

	// 4. Topocentric conversion (Alt/Az)
	topoPos, _ := a.Manager.GetTopocentricPosition(Moon, et, latitude, longitude)
	alt, az := a.Manager.GetLocalAltAz(topoPos, latitude, longitude)

	return MoonTelemetry{
		Altitude:     alt,
		Azimuth:      az,
		Elongation:   elongation,
		Illumination: illumination,
		DistanceKM:   moonPos.Norm(),
		PhaseName:    getPhaseName(illumination, sunToMoon.Dot(earthToMoon) < 0),
		Timestamp:    dt.UTC(),
	}, nil
}
