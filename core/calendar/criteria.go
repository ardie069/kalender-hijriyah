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
