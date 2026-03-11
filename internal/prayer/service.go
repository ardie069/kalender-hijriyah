package prayer

import (
	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"time"
)

type PrayerService struct {
	Astro *astronomy.Adapter
}

func NewPrayerService(astro *astronomy.Adapter) *PrayerService {
	return &PrayerService{
		Astro: astro,
	}
}

type PrayerTimes struct {
	Fajr       time.Time `json:"fajr"`
	Sunrise    time.Time `json:"sunrise"`
	Dhuhr      time.Time `json:"dhuhr"`
	Asr        time.Time `json:"asr"`
	Maghrib    time.Time `json:"maghrib"`
	Isha       time.Time `json:"isha"`
	Midnight   time.Time `json:"midnight"`
	ThirdNight time.Time `json:"third_night"`
}

func (s *PrayerService) GetPrayerTimes(date time.Time, lat, lon float64) PrayerTimes {
	dhuhr := s.Astro.Manager.GetSolarTransit(date, lat, lon)

	maghrib, _ := s.Astro.Manager.GetSunset(date, lat, lon)
	sunrise, _ := s.Astro.Manager.GetSunrise(date, lat, lon)

	fajr, _ := s.Astro.Manager.GetTimeByAltitude(date, lat, lon, -20.0, true)
	isha, _ := s.Astro.Manager.GetTimeByAltitude(date, lat, lon, -18.0, false)

	asr := s.Astro.Manager.GetAsrTime(dhuhr, lat, lon)

	nextFajr, _ := s.Astro.Manager.GetTimeByAltitude(date.AddDate(0, 0, 1), lat, lon, -20.0, true)

	nightDuration := nextFajr.Sub(maghrib)
	midnight := maghrib.Add(nightDuration / 2)
	thirdNight := maghrib.Add(nightDuration * 2 / 3)

	return PrayerTimes{
		Fajr:       fajr,
		Sunrise:    sunrise,
		Dhuhr:      dhuhr,
		Asr:        asr,
		Maghrib:    maghrib,
		Isha:       isha,
		Midnight:   midnight,
		ThirdNight: thirdNight,
	}
}
