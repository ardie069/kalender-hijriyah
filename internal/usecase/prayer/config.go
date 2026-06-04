package prayer

import "slices"

import "strings"

// ====== Madhab (Asr Calculation) ======

type Madhab int

const (
	SHAFII Madhab = 1 // Shadow factor = 1 (Syafi'i, Maliki, Hanbali)
	HANAFI Madhab = 2 // Shadow factor = 2 (Hanafi)
)

// ====== High Latitude Adjustment ======

type HighLatMethod int

const (
	HIGH_LAT_NONE             HighLatMethod = iota // Tidak ada adjustment
	HIGH_LAT_MIDDLE_OF_NIGHT                       // Fajr & Isha = ½ durasi malam
	HIGH_LAT_SEVENTH_OF_NIGHT                      // Fajr = 6/7 malam, Isha = 1/7 malam
	HIGH_LAT_ANGLE_BASED                           // Porsi = angle/60 dari durasi malam
)

// ====== Prayer Calculation Config ======

type PrayerConfig struct {
	Name      string  // Nama metode (e.g. "KEMENAG")
	FajrAngle float64 // Sudut depresi matahari untuk Subuh (negatif, misal -20.0)
	IshaAngle float64 // Sudut depresi matahari untuk Isya (negatif, misal -18.0)

	// Jika > 0, Isya dihitung sebagai Maghrib + offset (menit), abaikan IshaAngle.
	// Digunakan oleh Umm Al-Qura.
	IshaOffsetMin float64

	// Jika > 0, offset khusus Ramadan (menit). Umm Al-Qura: 120 min.
	IshaOffsetRamadanMin float64

	Madhab           Madhab        // Syafi'i atau Hanafi
	HighLatMethod    HighLatMethod // Metode penyesuaian lintang tinggi
	HighLatThreshold float64       // Lintang absolut di mana adjustment aktif (default: 48.5)
}

// DefaultConfig mengembalikan konfigurasi default (Kemenag/Syafi'i)
func DefaultConfig() PrayerConfig {
	return GetPreset("KEMENAG")
}

// GetPreset mengembalikan PrayerConfig berdasarkan nama metode
func GetPreset(name string) PrayerConfig {
	name = strings.ToUpper(strings.TrimSpace(name))

	base := PrayerConfig{
		Madhab:           SHAFII,
		HighLatMethod:    HIGH_LAT_ANGLE_BASED,
		HighLatThreshold: 48.5,
	}

	switch name {
	case "KEMENAG":
		base.Name = "KEMENAG"
		base.FajrAngle = -20.0
		base.IshaAngle = -18.0

	case "MWL":
		base.Name = "MWL"
		base.FajrAngle = -18.0
		base.IshaAngle = -17.0

	case "ISNA":
		base.Name = "ISNA"
		base.FajrAngle = -15.0
		base.IshaAngle = -15.0

	case "EGYPTIAN":
		base.Name = "EGYPTIAN"
		base.FajrAngle = -19.5
		base.IshaAngle = -17.5

	case "UMM_AL_QURA":
		base.Name = "UMM_AL_QURA"
		base.FajrAngle = -18.5
		base.IshaOffsetMin = 90.0
		base.IshaOffsetRamadanMin = 120.0

	case "KARACHI":
		base.Name = "KARACHI"
		base.FajrAngle = -18.0
		base.IshaAngle = -18.0

	case "TEHRAN":
		base.Name = "TEHRAN"
		base.FajrAngle = -17.7
		base.IshaAngle = -14.0

	case "JAKIM":
		base.Name = "JAKIM"
		base.FajrAngle = -20.0
		base.IshaAngle = -18.0

	case "MUIS":
		base.Name = "MUIS"
		base.FajrAngle = -20.0
		base.IshaAngle = -18.0

	default:
		// Fallback: pakai Kemenag
		base.Name = "KEMENAG"
		base.FajrAngle = -20.0
		base.IshaAngle = -18.0
	}

	return base
}

// IsValidMethod mengecek apakah nama metode valid
func IsValidMethod(name string) bool {
	valid := []string{"KEMENAG", "MWL", "ISNA", "EGYPTIAN", "UMM_AL_QURA", "KARACHI", "TEHRAN", "JAKIM", "MUIS"}
	upper := strings.ToUpper(strings.TrimSpace(name))
	return slices.Contains(valid, upper)
}

// ParseMadhab mengubah string ke Madhab
func ParseMadhab(s string) Madhab {
	if strings.ToUpper(strings.TrimSpace(s)) == "HANAFI" {
		return HANAFI
	}
	return SHAFII
}

// MadhabName mengembalikan nama string dari Madhab
func MadhabName(m Madhab) string {
	if m == HANAFI {
		return "Hanafi"
	}
	return "Syafi'i"
}

// ParseHighLat mengubah string ke HighLatMethod
func ParseHighLat(s string) HighLatMethod {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "MIDDLE_OF_NIGHT":
		return HIGH_LAT_MIDDLE_OF_NIGHT
	case "SEVENTH_OF_NIGHT":
		return HIGH_LAT_SEVENTH_OF_NIGHT
	case "ANGLE_BASED":
		return HIGH_LAT_ANGLE_BASED
	case "NONE":
		return HIGH_LAT_NONE
	default:
		return HIGH_LAT_ANGLE_BASED
	}
}

// HighLatName mengembalikan nama string dari HighLatMethod
func HighLatName(h HighLatMethod) string {
	switch h {
	case HIGH_LAT_MIDDLE_OF_NIGHT:
		return "MIDDLE_OF_NIGHT"
	case HIGH_LAT_SEVENTH_OF_NIGHT:
		return "SEVENTH_OF_NIGHT"
	case HIGH_LAT_ANGLE_BASED:
		return "ANGLE_BASED"
	default:
		return "NONE"
	}
}
