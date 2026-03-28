package calendar

import (
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"github.com/ardie069/kalender-hijriyah/core/models"
	"math"
	"sync"
	"time"
)

var (
	vectorCache   = make(map[string]*astronomy.DayVectors)
	vectorCacheMu sync.Mutex
)

// KHGT Parameters berdasarkan Muhammadiyah
const (
	GlobalLonStart = -180.0
	GlobalLonEnd   = 180.0
	GlobalLatStart = 65.0
	GlobalLatEnd   = -65.0

	AmericaLonStart = -170.0
	AmericaLonEnd   = -35.0
	AmericaLatStart = 60.0
	AmericaLatEnd   = -55.0

	NZLat = -41.2866
	NZLon = 174.7756

	MinElongation = 8.0
	MinAltitude   = 5.0
)

func (l *Logic) ScanGlobalKHGT(targetDateUTC time.Time, ijtimaTimeUTC time.Time) *models.KHGTResult {
	deadlineUTC := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 23, 59, 59, 999999999, time.UTC)

	// 1. Get or Pre-calculate geocentric vectors (IAU_EARTH frame)
	dateKey := targetDateUTC.Format("2006-01-02")
	vectorCacheMu.Lock()
	vectors, exists := vectorCache[dateKey]
	if !exists {
		vectors = l.Astro.Manager.PrecalculateDayVectors(targetDateUTC)
		vectorCache[dateKey] = vectors
	}
	vectorCacheMu.Unlock()

	// 2. Pre-calculate Fajr NZ for America Exception
	nextDay := targetDateUTC.AddDate(0, 0, 1)
	fajrNZ, _ := l.Astro.GetFajr(nextDay, NZLat, NZLon)
	isIjtimaBeforeFajrNZ := ijtimaTimeUTC.Before(fajrNZ)

	result := &models.KHGTResult{
		Date:               targetDateUTC,
		IsGlobalValid:      false,
		IsAmericaException: false,
		ValidLocations:     []models.HilalPrediction{},
	}

	var bestPred *models.HilalPrediction
	var bestScore float64 = -999.0

	// 3. Grid Scan
	step := 5.0
	for ln := GlobalLonStart; ln <= GlobalLonEnd; ln += step {
		for lat := GlobalLatStart; lat >= GlobalLatEnd; lat -= step {
			approxNoonUTC := 12.0 - (ln / 15.0)
			low := approxNoonUTC
			high := approxNoonUTC + 12.0

			baseTime := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 0, 0, 0, 0, time.UTC)
			var preciseSunset time.Time
			for i := 0; i < 10; i++ {
				mid := (low + high) / 2
				t := baseTime.Add(time.Duration(mid * float64(time.Hour)))
				sunVec := vectors.InterpolateVector("SUN", t)
				alt, _ := l.Astro.Manager.GetLocalAltAz(sunVec, lat, ln)
				if alt > -0.833 {
					low = mid
				} else {
					high = mid
				}
			}
			preciseSunset = baseTime.Add(time.Duration(((low + high) / 2) * float64(time.Hour)))

			if preciseSunset.Before(ijtimaTimeUTC) {
				continue
			}

			moonVec := vectors.InterpolateVector("MOON", preciseSunset)
			sunVec := vectors.InterpolateVector("SUN", preciseSunset)
			altGeo, _ := l.Astro.Manager.GetLocalAltAz(moonVec, lat, ln)

			uSun := sunVec.Unit()
			uMoon := moonVec.Unit()
			cosElong := uSun.X*uMoon.X + uSun.Y*uMoon.Y + uSun.Z*uMoon.Z
			if cosElong > 1.0 {
				cosElong = 1.0
			}
			if cosElong < -1.0 {
				cosElong = -1.0
			}
			elongGeo := math.Acos(cosElong) * 180.0 / math.Pi

			isValid := false
			var isAmEx bool
			if altGeo >= MinAltitude && elongGeo >= MinElongation {
				if preciseSunset.Before(deadlineUTC) || preciseSunset.Equal(deadlineUTC) {
					isValid = true
				} else {
					inAmerica := ln >= AmericaLonStart && ln <= AmericaLonEnd && lat >= AmericaLatEnd && lat <= AmericaLatStart
					if inAmerica && isIjtimaBeforeFajrNZ {
						isValid = true
						isAmEx = true
					}
				}
			}

			if altGeo > -20 {
				score := altGeo + elongGeo

				isBetter := false
				if bestPred == nil {
					isBetter = true
				} else if isValid && !bestPred.IsNewMonth {
					isBetter = true // Prioritas tinggi jika valid
				} else if isValid == bestPred.IsNewMonth {
					// Jika keduanya sama-sama valid atau sama-sama invalid, cek score
					if score > bestScore {
						isBetter = true
					} else if score == bestScore {
						// Tie breaker deterministik berdasarkan koordinat
						if ln > bestPred.Location.Lon {
							isBetter = true
						} else if ln == bestPred.Location.Lon && lat > bestPred.Location.Lat {
							isBetter = true
						}
					}
				}

				if isBetter {
					bestScore = score
					bestPred = &models.HilalPrediction{
						CheckDateUTC: preciseSunset,
						IjtimaTime:   ijtimaTimeUTC,
						Altitude:     altGeo,
						Elongation:   elongGeo,
						AgeHours:     preciseSunset.Sub(ijtimaTimeUTC).Hours(),
						Location:     &models.LocationInfo{Lat: lat, Lon: ln},
						IsNewMonth:   isValid,
					}
				}

				if isValid {
					result.IsGlobalValid = true
					if isAmEx {
						result.IsAmericaException = true
					}
				}
			}
		}
	}
	result.BestVisibility = bestPred
	return result
}
