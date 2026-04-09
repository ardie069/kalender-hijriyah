package calendar

import "time"

// IsMABIMS: Syarat Alt >= 3 derajat & Elong >= 6.4 derajat
func IsMABIMS(altitude, elongation float64) bool {
	return altitude >= 3.0 && elongation >= 6.4
}

// IsTurkey2016: Syarat Alt >= 5 derajat & Elong >= 8 derajat
func IsTurkey2016(altitude, elongation float64) bool {
	return altitude >= 5.0 && elongation >= 8.0
}

// IsUmmAlQura: Syarat Ijtima qoblal ghurub & Moonset ba'da ghurub (Makkah)
func IsUmmAlQura(ijtimaTime, sunsetTime, moonsetTime time.Time) bool {
	return ijtimaTime.Before(sunsetTime) && moonsetTime.After(sunsetTime)
}

// IsWujudulHilal: Syarat Ijtima qoblal sunset & Alt > 0 derajat
func IsWujudulHilal(altitude float64, ijtimaHappened bool) bool {
	return ijtimaHappened && altitude > 0.0
}

func CheckLagTime(sunset, moonset time.Time) (bool, float64) {
	diff := moonset.Sub(sunset).Minutes()
	return diff > 0, diff
}

// GetYallopCategory returns visibility category (A-F) based on Yallop (1998)
// Uses geocentric altitude (h) and geocentric elongation (L)
func GetYallopCategory(h, L float64) string {
	q := (h - (-0.014*L*L + 0.1*L + 2.53)) / 10.0
	switch {
	case q > 0.422:
		return "A" // Easily visible
	case q > 0.222:
		return "B" // Visible under perfect conditions
	case q > 0.022:
		return "C" // May need optical aid
	case q > -0.178:
		return "D" // Visible with optical aid only
	case q > -0.278:
		return "E" // Not visible even with optical aid
	default:
		return "F" // Below Danjon limit
	}
}

// GetOdehCategory returns visibility category (A-D) based on Odeh (2006)
// Uses arc of vision (ArcV) and topocentric crescent width (w) in arcminutes
func GetOdehCategory(arcV, w float64) string {
	// Formula from legacy ahc.hilal:
	// V = ArcV - (-0.1018*w^3 + 0.7319*w^2 - 6.3226*w + 7.1651)
	v := arcV - (7.1651 - 6.3226*w + 0.7319*w*w - 0.1018*w*w*w)
	switch {
	case v >= 5.65:
		return "A" // Easily visible
	case v >= 2.0:
		return "B" // Visible under perfect conditions
	case v >= -0.96:
		return "C" // Visible with optical aid only
	default:
		return "D" // Impossible
	}
}

