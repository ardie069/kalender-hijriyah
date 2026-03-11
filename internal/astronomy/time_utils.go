package astronomy

import "time"

// TimeToEt: Versi Matematika Murni (TIDAK BUTUH MUTEX)
// Sangat disarankan untuk digunakan di dalam loop scanning agar tidak rebutan Mutex
func TimeToEt(t time.Time) float64 {
	// Konversi UTC ke TDB (pendekatan J2000)
	// J2000 Epoch: 946728000 Unix
	// Delta T 2026: ~69.184 detik
	const unixJ2000 = 946728000
	const deltaT = 69.184
	return float64(t.Unix()) - unixJ2000 + deltaT
}

// Et2Utc: Konversi Ephemeris Time kembali ke UTC (inverse dari TimeToEt)
func Et2Utc(et float64) time.Time {
	const unixJ2000 = 946728000
	const deltaT = 69.184
	unix := int64(et + unixJ2000 - deltaT)
	return time.Unix(unix, 0).UTC()
}
