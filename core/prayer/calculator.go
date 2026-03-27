package prayer

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/models"
)

// Calculator handles prayer time computation using astronomical data.
type Calculator struct {
	Astro *astronomy.Adapter
}

// NewCalculator creates a new prayer Calculator.
func NewCalculator(astro *astronomy.Adapter) *Calculator {
	return &Calculator{Astro: astro}
}

// GetPrayerTimes calculates prayer times based on the given method configuration.
func (c *Calculator) GetPrayerTimes(date time.Time, lat, lon float64, cfg PrayerConfig) (models.PrayerTimes, error) {
	// Kalkulasi secara sequential karena SPICE bridge menggunakan global mutex (spiceMu).
	// Parallelizing calls that wait for the same mutex adds overhead without speedup.

	dhuhr, err := c.Astro.Manager.GetSolarTransit(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	maghrib, err := c.Astro.Manager.GetSunset(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	fajr, err := c.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.FajrAngle, true, "SUN")
	if err != nil {
		return models.PrayerTimes{}, err
	}

	sunrise, err := c.Astro.Manager.GetSunrise(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	asr, err := c.Astro.Manager.GetAsrTimeWithFactor(dhuhr, lat, lon, int(cfg.Madhab))
	if err != nil {
		return models.PrayerTimes{}, err
	}

	var isha time.Time
	if cfg.IshaOffsetMin > 0 {
		offset := cfg.IshaOffsetMin
		// Cek apakah sedang bulan Ramadan (menggunakan calendar logic yang stabil)
		hDate := calendar.GetTabularHijri(date)
		if cfg.IshaOffsetRamadanMin > 0 && hDate.Month == 9 {
			offset = cfg.IshaOffsetRamadanMin
		}
		isha = maghrib.Add(time.Duration(offset) * time.Minute)
	} else {
		isha, err = c.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.IshaAngle, false, "SUN")
		if err != nil {
			return models.PrayerTimes{}, err
		}
	}

	nextFajr, err := c.Astro.Manager.GetTimeByAltitude(date.AddDate(0, 0, 1), lat, lon, cfg.FajrAngle, true, "SUN")
	if err != nil {
		// Jika gagal cari besok subuh, fallback ke 24 jam setelah subuh hari ini (estimasi kasar)
		nextFajr = fajr.Add(24 * time.Hour)
	}

	nightDuration := nextFajr.Sub(maghrib)
	midnight := maghrib.Add(nightDuration / 2)
	thirdNight := maghrib.Add(nightDuration * 2 / 3)

	// High Latitude Adjustment
	if cfg.HighLatMethod != HIGH_LAT_NONE && math.Abs(lat) >= cfg.HighLatThreshold {
		fajr, isha = adjustHighLatitude(cfg, fajr, isha, sunrise, maghrib, nightDuration)
		// Recalculate night-derived times after adjustment
		nightDuration = nextFajr.Sub(maghrib)
		midnight = maghrib.Add(nightDuration / 2)
		thirdNight = maghrib.Add(nightDuration * 2 / 3)
	}

	return models.PrayerTimes{
		Fajr:       fajr,
		Sunrise:    sunrise,
		Dhuhr:      dhuhr,
		Asr:        asr,
		Maghrib:    maghrib,
		Isha:       isha,
		Midnight:   midnight,
		ThirdNight: thirdNight,
	}, nil
}

func adjustHighLatitude(cfg PrayerConfig, fajr, isha, sunrise, maghrib time.Time, nightDur time.Duration) (time.Time, time.Time) {
	switch cfg.HighLatMethod {
	case HIGH_LAT_MIDDLE_OF_NIGHT:
		half := nightDur / 2
		fajr = sunrise.Add(-half)
		isha = maghrib.Add(half)

	case HIGH_LAT_SEVENTH_OF_NIGHT:
		seventh := nightDur / 7
		fajr = sunrise.Add(-seventh)
		isha = maghrib.Add(seventh)

	case HIGH_LAT_ANGLE_BASED:
		fajrPortion := math.Abs(cfg.FajrAngle) / 60.0
		fajr = sunrise.Add(-time.Duration(float64(nightDur) * fajrPortion))

		if cfg.IshaOffsetMin <= 0 {
			ishaPortion := math.Abs(cfg.IshaAngle) / 60.0
			isha = maghrib.Add(time.Duration(float64(nightDur) * ishaPortion))
		}
	}
	return fajr, isha
}
