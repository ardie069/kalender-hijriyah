package prayer

import (
	"math"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/models"
	"github.com/ardie069/kalender-hijriyah/internal/prayer/highlatitude"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Service struct {
	Astro *cspice.Adapter
}

func NewService(astro *cspice.Adapter) *Service {
	return &Service{Astro: astro}
}

func (s *Service) GetPrayerTimes(date time.Time, lat, lon float64, cfg Config) (models.PrayerTimes, error) {
	dhuhr, err := s.Astro.Manager.GetSolarTransit(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	maghrib, err := s.Astro.Manager.GetSunset(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	fajr, err := s.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.FajrAngle, true, "SUN")
	if err != nil {
		return models.PrayerTimes{}, err
	}

	sunrise, err := s.Astro.Manager.GetSunrise(date, lat, lon)
	if err != nil {
		return models.PrayerTimes{}, err
	}

	asr, err := s.Astro.Manager.GetAsrTimeWithFactor(dhuhr, lat, lon, cfg.Madhab.ShadowFactor())
	if err != nil {
		return models.PrayerTimes{}, err
	}

	var isha time.Time
	if cfg.IshaOffsetMin > 0 {
		offset := cfg.IshaOffsetMin
		hDate := calendar.GetTabularHijri(date)
		if cfg.IshaOffsetRamadanMin > 0 && hDate.Month == 9 {
			offset = cfg.IshaOffsetRamadanMin
		}
		isha = maghrib.Add(time.Duration(offset) * time.Minute)
	} else {
		isha, err = s.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.IshaAngle, false, "SUN")
		if err != nil {
			return models.PrayerTimes{}, err
		}
	}

	nextFajr, err := s.Astro.Manager.GetTimeByAltitude(date.AddDate(0, 0, 1), lat, lon, cfg.FajrAngle, true, "SUN")
	if err != nil {
		nextFajr = fajr.Add(24 * time.Hour)
	}

	nightDuration := nextFajr.Sub(maghrib)
	midnight := maghrib.Add(nightDuration / 2)
	thirdNight := maghrib.Add(nightDuration * 2 / 3)

	if cfg.HighLatMethod != nil && math.Abs(lat) >= cfg.HighLatThreshold {
		params := highlatitude.AdjustParams{
			Fajr:          fajr,
			Isha:          isha,
			Sunrise:       sunrise,
			Maghrib:       maghrib,
			NightDuration: nightDuration,
			FajrAngle:     cfg.FajrAngle,
			IshaAngle:     cfg.IshaAngle,
			IshaOffsetMin: cfg.IshaOffsetMin,
		}
		fajr, isha = cfg.HighLatMethod.Adjust(params)
		
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
