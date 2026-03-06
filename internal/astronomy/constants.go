package astronomy

const (
	// NAIF Body IDs
	Sun   = "SUN"
	Earth = "EARTH"
	Moon  = "MOON"

	// Frames
	FrameJ2000      = "J2000"      // Equatorial
	FrameEclipJ2000 = "ECLIPJ2000" // Ecliptic

	// Aberration Corrections
	CorrLT  = "LT"   // Light Time
	CorrLTS = "LT+S" // Light Time + Stellar Aberration

	TimeFormat = "2006-01-02 15:04:05 UTC"
)
