package prayer

import (
	"slices"
	"strings"

	"github.com/ardie069/kalender-hijriyah/internal/prayer/highlatitude"
	"github.com/ardie069/kalender-hijriyah/internal/prayer/madhab"
)

type Config struct {
	Name                 string
	FajrAngle            float64
	IshaAngle            float64
	IshaOffsetMin        float64
	IshaOffsetRamadanMin float64
	Madhab               madhab.Madhab
	HighLatMethod        highlatitude.Method
	HighLatThreshold     float64
}

func DefaultConfig() Config {
	return GetPreset("KEMENAG")
}

func GetPreset(name string) Config {
	name = strings.ToUpper(strings.TrimSpace(name))

	base := Config{
		Madhab:           madhab.Shafii{},
		HighLatMethod:    highlatitude.AngleBased{},
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
		base.Name = "KEMENAG"
		base.FajrAngle = -20.0
		base.IshaAngle = -18.0
	}

	return base
}

func IsValidMethod(name string) bool {
	valid := []string{"KEMENAG", "MWL", "ISNA", "EGYPTIAN", "UMM_AL_QURA", "KARACHI", "TEHRAN", "JAKIM", "MUIS"}
	upper := strings.ToUpper(strings.TrimSpace(name))
	return slices.Contains(valid, upper)
}

func ParseMadhab(s string) madhab.Madhab {
	if strings.ToUpper(strings.TrimSpace(s)) == "HANAFI" {
		return madhab.Hanafi{}
	}
	return madhab.Shafii{}
}

func ParseHighLat(s string) highlatitude.Method {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "MIDDLE_OF_NIGHT":
		return highlatitude.MiddleNight{}
	case "SEVENTH_OF_NIGHT":
		return highlatitude.Seventh{}
	case "ANGLE_BASED":
		return highlatitude.AngleBased{}
	case "NONE":
		return nil
	default:
		return highlatitude.AngleBased{}
	}
}
