package astronomy

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/models"
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

	// 1. Geocentric Positions (Inertial frame is fine for elongation)
	sunPosGeo, _ := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	moonPosGeo, _ := a.Manager.GetGeocentric(Moon, et, FrameJ2000)
	elongGeo := math.Acos(sunPosGeo.Unit().Dot(moonPosGeo.Unit())) * (180.0 / math.Pi)

	// 3. Topocentric Positions (In ECEF)
	_, _ = a.Manager.GetTopocentricPosition(Sun, et, latitude, longitude)
	topoMoon, _ := a.Manager.GetTopocentricPosition(Moon, et, latitude, longitude)

	// 4. Illumination & Phase (Using inertial frame)
	sunToMoon := sunPosGeo.Sub(moonPosGeo).Unit()
	earthToMoon := moonPosGeo.Scale(-1).Unit()
	phaseAngle := math.Acos(sunToMoon.Dot(earthToMoon))
	illumination := (1.0 + math.Cos(phaseAngle)) / 2.0 * 100.0

	// 5. Alt/Az
	altTopo, azTopo := a.Manager.GetLocalAltAz(topoMoon, latitude, longitude)

	altApparent := altTopo + ApplyRefraction(altTopo)
	return models.MoonTelemetry{
		Altitude:         altTopo,
		AltitudeApparent: &altApparent,
		Azimuth:          azTopo,
		Elongation:       elongGeo,
		Illumination:     illumination,
		DistanceKM:       moonPosGeo.Norm(),
		PhaseName:        getPhaseName(illumination, sunToMoon.Dot(earthToMoon) < 0),
		Timestamp:        dt.UTC(),
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

// CalculateGeocentricParamsGlobal returns geocentric elongation and geocentric altitude relative to local horizon
func (a *Adapter) CalculateGeocentricParamsGlobal(dt time.Time, lat, lon float64) (altitude, elongation float64) {
	et, _ := Str2et(dt.UTC().Format(TimeFormat))

	sunPosGeo, _ := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	moonPosGeo, _ := a.Manager.GetGeocentric(Moon, et, FrameJ2000)
	elongation = math.Acos(sunPosGeo.Unit().Dot(moonPosGeo.Unit())) * (180.0 / math.Pi)

	// Use Body-Fixed for Altitude
	moonPosFixed, _ := a.Manager.GetGeocentric(Moon, et, "IAU_EARTH")
	altitude, _ = a.Manager.GetLocalAltAz(moonPosFixed, lat, lon)

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

func (a *Adapter) GetSunsetFast(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetSunsetFast(date, lat, lon)
}

func (a *Adapter) GetIsha(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetIsha(date, lat, lon)
}
