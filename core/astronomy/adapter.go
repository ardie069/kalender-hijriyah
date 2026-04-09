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

	// 1. Topocentric Positions (For real observers)
	topoSun, _ := a.Manager.GetTopocentricPosition(Sun, et, latitude, longitude)
	topoMoon, _ := a.Manager.GetTopocentricPosition(Moon, et, latitude, longitude)
	elongTopo := math.Acos(topoSun.Unit().Dot(topoMoon.Unit())) * (180.0 / math.Pi)

	// 2. Geocentric Positions (For illumination/phase calculations)
	sunPosGeo, _ := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	moonPosGeo, _ := a.Manager.GetGeocentric(Moon, et, FrameJ2000)

	// 3. Topocentric Positions (In ECEF)
	// No changes needed here, just updating the return mapping below

	// 4. Illumination & Phase (Using inertial frame)
	sunToMoon := sunPosGeo.Sub(moonPosGeo).Unit()
	earthToMoon := moonPosGeo.Scale(-1).Unit()
	phaseAngle := math.Acos(sunToMoon.Dot(earthToMoon))
	illumination := (1.0 + math.Cos(phaseAngle)) / 2.0 * 100.0

	// Determine Waxing/Waning using Ecliptic Longitude
	sunEclip, _ := a.Manager.GetGeocentric(Sun, et, FrameEclipJ2000)
	moonEclip, _ := a.Manager.GetGeocentric(Moon, et, FrameEclipJ2000)
	sunLong := math.Atan2(sunEclip.Y, sunEclip.X)
	moonLong := math.Atan2(moonEclip.Y, moonEclip.X)
	
	diff := moonLong - sunLong
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	isWaxing := diff > 0

	// 5. Alt/Az
	altTopo, azTopo := a.Manager.GetLocalAltAz(topoMoon, latitude, longitude)

	altApparent := altTopo + ApplyRefraction(altTopo)
	// 6. Moonrise & Moonset: Cari terbit hari ini, dan terbenam setelah terbit tersebut.
	moonrise, _ := a.Manager.GetMoonrise(dt, latitude, longitude)
	
	// Cari terbenam. Jika terbit ditemukan, cari terbenam setelah waktu terbit.
	// Jika tidak, cari terbenam terdekat hari ini.
	var moonset time.Time
	if !moonrise.IsZero() {
		// Cari terbenam dalam 24 jam setelah terbit
		moonset, _ = a.Manager.FindAltitudeCrossing("MOON", moonrise, moonrise.Add(24*time.Hour), latitude, longitude, 0.0, false)
	} else {
		moonset, _ = a.Manager.GetMoonset(dt, latitude, longitude)
	}

	return models.MoonTelemetry{
		Altitude:         altTopo,
		AltitudeApparent: &altApparent,
		Azimuth:          azTopo,
		Elongation:       elongTopo,
		Illumination:     illumination,
		DistanceKM:       moonPosGeo.Norm(),
		PhaseName:        getPhaseName(illumination, isWaxing),
		Moonrise:         &moonrise,
		Moonset:          &moonset,
		Timestamp:        dt.UTC(),
	}, nil
}

// CalculateGeocentricParams returns topocentric parameters for local evaluation (e.g., MABIMS)
func (a *Adapter) CalculateGeocentricParams(et float64, lat, lon float64) (altitude, elongation float64) {
	topoSun, _ := a.Manager.GetTopocentricPosition(Sun, et, lat, lon)
	topoMoon, err := a.Manager.GetTopocentricPosition(Moon, et, lat, lon)
	if err != nil {
		return 0, 0
	}

	elongation = math.Acos(topoSun.Unit().Dot(topoMoon.Unit())) * (180.0 / math.Pi)
	altitude, _ = a.Manager.GetLocalAltAz(topoMoon, lat, lon)

	return altitude, elongation
}

// CalculateGeocentricParamsGlobal returns geocentric elongation and geocentric altitude relative to local horizon
func (a *Adapter) CalculateGeocentricParamsGlobal(dt time.Time, lat, lon float64) (altitude, elongation, arcv, width float64) {
	et := TimeToEt(dt)

	sunPosGeo, _ := a.Manager.GetGeocentric(Sun, et, FrameJ2000)
	moonPosGeo, _ := a.Manager.GetGeocentric(Moon, et, FrameJ2000)
	elongation = math.Acos(sunPosGeo.Unit().Dot(moonPosGeo.Unit())) * (180.0 / math.Pi)

	// --- Manual Earth Rotation for Altitude ---
	d := et / 86400.0
	gmstDeg := 280.46061837 + 360.98564736629 * d
	gmstRad := gmstDeg * math.Pi / 180.0
	cosG, sinG := math.Cos(gmstRad), math.Sin(gmstRad)

	// Rotate Moon to ECEF (Geocentric: observer at center of Earth)
	moonECEF := Vector3{
		X: moonPosGeo.X*cosG + moonPosGeo.Y*sinG,
		Y: -moonPosGeo.X*sinG + moonPosGeo.Y*cosG,
		Z: moonPosGeo.Z,
	}

	altitude, _ = a.Manager.GetLocalAltAz(moonECEF, lat, lon)
	arcv = altitude // Simplified for Geocentric as Sun alt is 0 at center of Earth boundary? 
	// Actually for KHGT we use fixed 5/8 so ArcV is not strictly needed but we'll return altitude

	// Width: Simplified SD for geocentric (approx 15.5')
	width = (15.5 / 60.0) * (1.0 - math.Cos(elongation*math.Pi/180.0))

	return altitude, elongation, arcv, width
}

// CalculateTopocentricParamsGlobal returns parameters relative to an observer at (lat, lon)
// Includes both parallax and atmospheric refraction.
func (a *Adapter) CalculateTopocentricParamsGlobal(dt time.Time, lat, lon float64) (altitude, elongation, arcv, width float64) {
	et := TimeToEt(dt)

	// Get positions relative to surface observer (Topocentric)
	topoSun, _ := a.Manager.GetTopocentricPosition(Sun, et, lat, lon)
	topoMoon, _ := a.Manager.GetTopocentricPosition(Moon, et, lat, lon)

	elongation = math.Acos(topoSun.Unit().Dot(topoMoon.Unit())) * (180.0 / math.Pi)
	
	// Local Altitude/Azimuth from topocentric vectors
	altMoon, _ := a.Manager.GetLocalAltAz(topoMoon, lat, lon)
	altSun, _ := a.Manager.GetLocalAltAz(topoSun, lat, lon)
	
	// Apparent moon altitude (refracted)
	altitude = altMoon + ApplyRefraction(altMoon)
	
	// Arc of Vision (Legacy Odeh uses refracted moon alt - refracted sun alt)
	// At sunset, refracted sun altitude is approx -0.583 (0.2665 radius + 0.566 refraction)
	sunApparent := altSun + ApplyRefraction(altSun)
	arcv = altitude - sunApparent

	// Width: Calculate using dynamic Moon Semi-Diameter (SD)
	// SD = 0.27254 * horizontal_parallax
	// Horizontal parallax approx sin-1(Re/d)
	re := 6378.137
	dist := topoMoon.Norm()
	hp := math.Asin(re/dist) * (180.0 / math.Pi)
	sdDegrees := 0.27254 * hp
	width = sdDegrees * (1.0 - math.Cos(elongation*math.Pi/180.0))

	return altitude, elongation, arcv, width
}

func (a *Adapter) GetFajr(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetFajr(date, lat, lon)
}

func (a *Adapter) GetMoonrise(date time.Time, lat, lon float64) (time.Time, error) {
	return a.Manager.GetMoonrise(date, lat, lon)
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
