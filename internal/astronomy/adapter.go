package astronomy

import (
	"math"
	"sync"
	"time"
)

// MoonTelemetry nampung data hasil perhitungan bulan
type MoonTelemetry struct {
	Altitude     float64   `json:"altitude"`
	Azimuth      float64   `json:"azimuth"`
	Elongation   float64   `json:"elongation"`
	Illumination float64   `json:"illumination"`
	DistanceKM   float64   `json:"distance_km"`
	Timestamp    time.Time `json:"timestamp"`
}

type Adapter struct {
	Manager *EphemerisManager
}

var instance *Adapter
var once sync.Once

func GetAdapter(manager *EphemerisManager) *Adapter {
	return &Adapter{Manager: manager}
}

func (a *Adapter) GetMoonTelemetry(dt time.Time, lat, lon, elev float64) (MoonTelemetry, error) {
	utcStr := dt.UTC().Format("2006-01-02 15:04:05 UTC")
	et, err := Str2et(utcStr)
	if err != nil {
		return MoonTelemetry{}, err
	}

	// 1. Ambil Posisi J2000 (Semuanya sekarang sudah Vector3)
	sunPos, _ := a.Manager.GetPosition(NAIF_SUN, et)
	moonPos, _ := a.Manager.GetPosition(NAIF_MOON, et)
	earthPos, _ := a.Manager.GetPosition(NAIF_EARTH, et)

	// 2. Hitung Vektor Dasar (Pake method Sub biar rapi, Bang!)
	mVec := moonPos.Sub(earthPos) // Vektor Bumi ke Bulan
	sVec := sunPos.Sub(earthPos)  // Vektor Bumi ke Matahari
	dist := mVec.Norm()

	// 3. Hitung Elongasi (Sudut di Bumi)
	dotElong := sVec.Dot(mVec)
	denomElong := sVec.Norm() * mVec.Norm()
	elongDeg := 0.0
	if denomElong != 0 {
		elongDeg = math.Acos(math.Max(-1, math.Min(1, dotElong/denomElong))) * (180.0 / math.Pi)
	}

	// 4. Hitung Illumination (Sudut di Bulan)
	sm := sunPos.Sub(moonPos)   // Vektor Bulan ke Matahari
	em := earthPos.Sub(moonPos) // Vektor Bulan ke Bumi

	dotPhase := sm.Dot(em)
	denomPhase := sm.Norm() * em.Norm()
	phaseDeg := 0.0
	if denomPhase != 0 {
		phaseDeg = math.Acos(math.Max(-1, math.Min(1, dotPhase/denomPhase))) * (180.0 / math.Pi)
	}
	illum := (1.0 + math.Cos(phaseDeg*(math.Pi/180.0))) / 2.0 * 100.0

	// 5. TRANSFORMASI KE ALTITUDE/AZIMUTH (TOPOCENTRIC)
	// Kita pake koordinat mVec (Bumi ke Bulan)
	ra := math.Atan2(mVec.Y, mVec.X)
	dec := math.Atan2(mVec.Z, math.Sqrt(mVec.X*mVec.X+mVec.Y*mVec.Y))

	// Hitung Local Sidereal Time (LST)
	jd := (float64(dt.UTC().Unix()) / 86400.0) + 2440587.5
	d := jd - 2451545.0
	lst := math.Mod(18.697374558+24.06570982441908*d+lon/15.0, 24.0)
	if lst < 0 {
		lst += 24.0
	}
	ha := (lst * 15.0 * (math.Pi / 180.0)) - ra

	// Rumus Ketinggian (Altitude)
	latRad := lat * (math.Pi / 180.0)
	sinAlt := math.Sin(latRad)*math.Sin(dec) + math.Cos(latRad)*math.Cos(dec)*math.Cos(ha)
	altRad := math.Asin(sinAlt)
	altDeg := altRad * (180.0 / math.Pi)

	// Rumus Arah (Azimuth)
	cosAz := (math.Sin(dec) - math.Sin(latRad)*sinAlt) / (math.Cos(latRad) * math.Cos(altRad))
	azDeg := math.Acos(math.Max(-1, math.Min(1, cosAz))) * (180.0 / math.Pi)
	if math.Sin(ha) > 0 {
		azDeg = 360.0 - azDeg
	}

	return MoonTelemetry{
		Altitude:     altDeg,
		Azimuth:      azDeg,
		Elongation:   elongDeg,
		Illumination: illum,
		DistanceKM:   dist,
		Timestamp:    dt.UTC(),
	}, nil
}
