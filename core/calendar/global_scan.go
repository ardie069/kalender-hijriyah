package calendar

import (
	"github.com/ardie069/kalender-hijriyah/core/models"
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"math"
	"sync"
	"sync/atomic"
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

	AmericaLonStart = -168.0
	AmericaLonEnd   = -30.0
	AmericaLatStart = 60.0
	AmericaLatEnd   = -55.0

	NZLat = -41.2866
	NZLon = 174.7756

	MinElongation = 8.0
	MinAltitude   = 5.0
)

func (l *Logic) ScanGlobalKHGT(targetDateUTC time.Time, ijtimaTimeUTC time.Time) (bool, *models.HilalPrediction) {
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
	var mu sync.Mutex
	var foundGlobal int32
	var wg sync.WaitGroup

	// 3. Grid Scan (Pure Go)
	step := 15.0
	for lon := GlobalLonStart; lon <= GlobalLonEnd; lon += step {
		wg.Add(1)
		go func(ln float64) {
			defer wg.Done()
			for lat := GlobalLatStart; lat >= GlobalLatEnd; lat -= step {
				if atomic.LoadInt32(&foundGlobal) == 1 {
					return
				}

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
				if cosElong > 1.0 { cosElong = 1.0 }
				if cosElong < -1.0 { cosElong = -1.0 }
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
					mu.Lock()
					// bestPred should reflect the most "visible" point, but we prioritize valid ones
					if score > bestScore || (isValid && bestPred != nil && !bestPred.IsNewMonth) {
						bestScore = score
						bestPred = &models.HilalPrediction{
							CheckDateUTC:      preciseSunset,
							IjtimaTime:        ijtimaTimeUTC,
							Altitude:          altGeo,
							Elongation:        elongGeo,
							AgeHours:          preciseSunset.Sub(ijtimaTimeUTC).Hours(),
							Location:          &models.LocationInfo{Lat: lat, Lon: ln},
							IsNewMonth:        isValid,
						}
					}
					if isValid {
						result.IsGlobalValid = true
						if isAmEx { 
							result.IsAmericaException = true 
						} else {
							atomic.StoreInt32(&foundGlobal, 1)
						}
					}
					mu.Unlock()
				}
			}
		}(lon)
	}

	wg.Wait()
	return result.IsGlobalValid, bestPred
}

func (l *Logic) evaluateGeocentricPoint(
	targetDateUTC, ijtimaTimeUTC, deadlineUTC time.Time,
	lat, lon float64,
	americaExceptionFlag *bool,
) (bool, *models.HilalPrediction) {

	// Hitung sunset di lokasi ini
	sunsetUTC, err := l.Astro.GetSunset(targetDateUTC, lat, lon)
	if err != nil {
		return false, nil
	}

	if sunsetUTC.Before(ijtimaTimeUTC) {
		return false, nil
	}

	altGeo, elongGeo := l.Astro.CalculateGeocentricParamsGlobal(targetDateUTC, lat, lon)

	ageHours := sunsetUTC.Sub(ijtimaTimeUTC).Hours()
	pred := &models.HilalPrediction{
		CheckDateUTC:      sunsetUTC,
		IjtimaTime:        ijtimaTimeUTC,
		Altitude:          altGeo,
		Elongation:        elongGeo,
		AgeHours:          ageHours,
		Location:          &models.LocationInfo{Lat: lat, Lon: lon},
		IsNewMonth:        false,
	}

	if altGeo >= MinAltitude && elongGeo >= MinElongation {
		if sunsetUTC.Before(deadlineUTC) || sunsetUTC.Equal(deadlineUTC) {
			pred.IsNewMonth = true
			return true, pred
		}

		isAmericaLand := (lon >= AmericaLonStart && lon <= AmericaLonEnd &&
			lat >= AmericaLatEnd && lat <= AmericaLatStart)

		if isAmericaLand {
			nextDay := targetDateUTC.AddDate(0, 0, 1)
			fajrNZ, err := l.Astro.GetFajr(nextDay, NZLat, NZLon)
			if err == nil && ijtimaTimeUTC.Before(fajrNZ) {
				pred.IsNewMonth = true
				*americaExceptionFlag = true
				return true, pred
			}
		}

		return false, pred
	}

	return false, pred
}
