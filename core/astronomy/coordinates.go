package astronomy

import "math"

// GetTopocentricPosition: Hitung posisi benda langit dari sudut pandang pengamat di permukaan bumi
func (em *EphemerisManager) GetTopocentricPosition(target string, et, lat, lon float64) (Vector3, error) {
	spiceMu.Lock()
	defer spiceMu.Unlock()

	// Math Geodetic ke Rectangular (Bisa di luar Lock sebenarnya, tapi biar aman satu blok)
	re := 6378.137
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

	posGeo, err := getPositionInternal(target, "IAU_EARTH", "LT+S", "EARTH", et)
	if err != nil {
		return Vector3{}, err
	}

	return Vector3{
		X: posGeo.X - obsX,
		Y: posGeo.Y - obsY,
		Z: posGeo.Z - obsZ,
	}, nil
}

// GetLocalAltAz: Transformasi vektor toposentrik ke Altitude dan Azimuth lokal
func (em *EphemerisManager) GetLocalAltAz(pos Vector3, lat, lon float64) (float64, float64) {
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	sinLat, cosLat := math.Sin(latRad), math.Cos(latRad)
	sinLon, cosLon := math.Sin(lonRad), math.Cos(lonRad)

	// Transformasi ke Toposentrik (North-East-Down)
	x := -sinLat*cosLon*pos.X - sinLat*sinLon*pos.Y + cosLat*pos.Z
	y := -sinLon*pos.X + cosLon*pos.Y
	z := cosLat*cosLon*pos.X + cosLat*sinLon*pos.Y + sinLat*pos.Z

	hyp := math.Sqrt(x*x + y*y + z*z)
	if hyp < 1e-9 {
		return -90, 0 // Fallback untuk posisi nol (pusat bumi)
	}

	ratio := z / hyp
	if ratio > 1.0 {
		ratio = 1.0
	} else if ratio < -1.0 {
		ratio = -1.0
	}

	geoAlt := math.Asin(ratio) * 180.0 / math.Pi
	az := math.Atan2(y, x) * 180.0 / math.Pi
	if az < 0 {
		az += 360.0
	}
	return geoAlt, az
}
