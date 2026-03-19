package astronomy

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/models"
)

type Adapter struct {
	Manager *EphemerisManager
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

func (a *Adapter) GetMoonTelemetry(dt time.Time, latitude, longitude float64) (models.MoonTelemetry, error) {
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

	return models.MoonTelemetry{
		Altitude:     alt,
		Azimuth:      az,
		Elongation:   elongation,
		Illumination: illumination,
		DistanceKM:   moonPos.Norm(),
		PhaseName:    getPhaseName(illumination, sunToMoon.Dot(earthToMoon) < 0),
		Timestamp:    dt.UTC(),
	}, nil
}

// CalculateGeocentricParams returns geocentric elongation and topocentric altitude
func (a *Adapter) CalculateGeocentricParams(et float64, lat, lon float64) (altitude, elongation float64) {
	sunPos, err := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	if err != nil {
		return 0, 0
	}
	moonPos, err := a.Manager.GetGeocentric(Moon, et, FrameJ2000)
	if err != nil {
		return 0, 0
	}

	elongation = math.Acos(sunPos.Unit().Dot(moonPos.Unit())) * (180.0 / math.Pi)

	topoPos, err := a.Manager.GetTopocentricPosition(Moon, et, lat, lon)
	if err != nil {
		return 0, elongation
	}
	altitude, _ = a.Manager.GetLocalAltAz(topoPos, lat, lon)

	return altitude, elongation
}

func (a *Adapter) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetFajr(date, lat, lon)
}

func (a *Adapter) GetMoonset(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetMoonset(date, lat, lon)
}

func (a *Adapter) GetSunset(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetSunset(date, lat, lon)
}

func (a *Adapter) GetIsha(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetIsha(date, lat, lon)
}
