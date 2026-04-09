package astronomy

import (
	"fmt"
	"math"
)

// GetTopocentricPosition: Hitung posisi benda langit dari sudut pandang pengamat di permukaan bumi
func (em *EphemerisManager) GetTopocentricPosition(target string, et, lat, lon float64) (Vector3, error) {
	// 1. Ambil posisi benda dalam frame inertial J2000 (Selalu tersedia)
	posInertial, err := getPosition(target, "J2000", "LT+S", "EARTH", et)
	if err != nil {
		return Vector3{}, fmt.Errorf("J2000 lookup failed: %v", err)
	}

	// 2. Hitung rotasi Bumi (GMST) untuk mengubah J2000 ke ECEF (IAU_EARTH)
	// T = centuries since J2000.0
	d := et / 86400.0
	t := d / 36525.0
	
	// Greenwich Mean Sidereal Time (IAU 1982 formula)
	gmstDeg := 280.46061837 + 360.98564736629*d + 0.000387933*t*t - t*t*t/38710000.0
	
	// Normalize to [0, 360)
	for gmstDeg < 0 { gmstDeg += 360.0 }
	for gmstDeg >= 360 { gmstDeg -= 360.0 }
	
	gmstRad := gmstDeg * math.Pi / 180.0

	// Rotasi J2000 -> ECEF (Rotasi di sumbu Z)
	cosG, sinG := math.Cos(gmstRad), math.Sin(gmstRad)
	posECEF := Vector3{
		X: posInertial.X*cosG + posInertial.Y*sinG,
		Y: -posInertial.X*sinG + posInertial.Y*cosG,
		Z: posInertial.Z,
	}

	// 3. Hitung posisi pengamat (ECEF)
	re := 6378.137 // km
	f := 1.0 / 298.257223563
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0
	cosLat := math.Cos(latRad)
	sinLat := math.Sin(latRad)
	
	N := re / math.Sqrt(1-f*(2-f)*sinLat*sinLat)
	altKM := 0.45 
	obsX := (N + altKM) * cosLat * math.Cos(lonRad)
	obsY := (N + altKM) * cosLat * math.Sin(lonRad)
	obsZ := (N*(1-f)*(1-f) + altKM) * sinLat

	return Vector3{
		X: posECEF.X - obsX,
		Y: posECEF.Y - obsY,
		Z: posECEF.Z - obsZ,
	}, nil
}

// GetLocalAltAz: Transformasi vektor toposentrik ke Altitude dan Azimuth lokal (Toposentrik)
func (em *EphemerisManager) GetLocalAltAz(pos Vector3, lat, lon float64) (float64, float64) {
	x, y, z := ecefToHorizon(pos, lat, lon)

	dist := pos.Norm()
	if dist < 1e-9 {
		return -90, 0
	}

	alt := math.Asin(z/dist) * 180.0 / math.Pi
	az := math.Atan2(y, x) * 180.0 / math.Pi
	if az < 0 {
		az += 360.0
	}
	return alt, az
}

// GetGeocentricAltAz: Transformasi geometri dari IAU_EARTH ke horizon lokal (Geosentris)
// Digunakan untuk kriteria KHGT/Muhammadiyah di mana posisi bulan geosentris
// diproyeksikan ke horizon pengamat tanpa koreksi paralaks.
func (em *EphemerisManager) GetGeocentricAltAz(pos Vector3, lat, lon float64) (float64, float64) {
	x, y, z := ecefToHorizon(pos, lat, lon)

	// Jarak (dist) adalah norm dari vektor posisi (rotasi mempertahankan magnitude)
	dist := pos.Norm()
	if dist < 1e-9 {
		return -90, 0
	}

	// Hitung Altitude Geosentris
	geoAlt := math.Asin(z/dist) * 180.0 / math.Pi

	// Hitung Azimuth (North-based: North=0, East=90)
	az := math.Atan2(y, x) * 180.0 / math.Pi
	if az < 0 {
		az += 360.0
	}

	return geoAlt, az
}

// ecefToHorizon mengubah vektor dalam frame IAU_EARTH (ECEF) ke sistem koordinat lokal (NEU: North-East-Up)
func ecefToHorizon(pos Vector3, lat, lon float64) (x, y, z float64) {
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	sinLat, cosLat := math.Sin(latRad), math.Cos(latRad)
	sinLon, cosLon := math.Sin(lonRad), math.Cos(lonRad)

	// Matrix Rotasi ECEF -> Horizontal (North-East-Up)
	// x = North
	x = -sinLat*cosLon*pos.X - sinLat*sinLon*pos.Y + cosLat*pos.Z
	// y = East
	y = -sinLon*pos.X + cosLon*pos.Y
	// z = Up (Zenith)
	z = cosLat*cosLon*pos.X + cosLat*sinLon*pos.Y + sinLat*pos.Z

	return x, y, z
}