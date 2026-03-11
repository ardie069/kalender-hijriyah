package prayer

import (
	"math"
	"sync"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/models"
)

type PrayerService struct {
	Astro *astronomy.Adapter
}

func NewPrayerService(astro *astronomy.Adapter) *PrayerService {
	return &PrayerService{
		Astro: astro,
	}
}

// GetPrayerTimes menghitung waktu sholat berdasarkan konfigurasi metode.
func (s *PrayerService) GetPrayerTimes(date time.Time, lat, lon float64, cfg PrayerConfig) models.PrayerTimes {
	var (
		wg                           sync.WaitGroup
		fajr, sunrise, maghrib, isha time.Time
		dhuhr, asr, nextFajr         time.Time
	)

	// Phase 1: Kalkulasi independen secara paralel
	wg.Add(4)

	go func() {
		defer wg.Done()
		dhuhr = s.Astro.Manager.GetSolarTransit(date, lat, lon)
	}()

	go func() {
		defer wg.Done()
		maghrib, _ = s.Astro.Manager.GetSunset(date, lat, lon)
	}()

	go func() {
		defer wg.Done()
		fajr, _ = s.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.FajrAngle, true)
	}()

	go func() {
		defer wg.Done()
		sunrise, _ = s.Astro.Manager.GetSunrise(date, lat, lon)
	}()

	wg.Wait()

	// Phase 2: Kalkulasi yang butuh hasil Phase 1
	wg.Add(2)

	go func() {
		defer wg.Done()
		asr = s.Astro.Manager.GetAsrTimeWithFactor(dhuhr, lat, lon, int(cfg.Madhab))
	}()

	go func() {
		defer wg.Done()
		// Isya: bisa pakai offset (Umm Al-Qura) atau sudut
		if cfg.IshaOffsetMin > 0 {
			offset := cfg.IshaOffsetMin
			// Cek apakah sedang bulan Ramadan
			if cfg.IshaOffsetRamadanMin > 0 && isRamadan(date) {
				offset = cfg.IshaOffsetRamadanMin
			}
			isha = maghrib.Add(time.Duration(offset) * time.Minute)
		} else {
			isha, _ = s.Astro.Manager.GetTimeByAltitude(date, lat, lon, cfg.IshaAngle, false)
		}
	}()

	wg.Wait()

	// Phase 3: Next Fajr + durasi malam
	nextFajr, _ = s.Astro.Manager.GetTimeByAltitude(date.AddDate(0, 0, 1), lat, lon, cfg.FajrAngle, true)

	nightDuration := nextFajr.Sub(maghrib)
	midnight := maghrib.Add(nightDuration / 2)
	thirdNight := maghrib.Add(nightDuration * 2 / 3)

	// Phase 4: High Latitude Adjustment
	// Diterapkan jika lintang melewati threshold DAN metode bukan NONE
	if cfg.HighLatMethod != HIGH_LAT_NONE && math.Abs(lat) >= cfg.HighLatThreshold {
		fajr, isha = s.adjustHighLatitude(cfg, fajr, isha, sunrise, maghrib, nightDuration)
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
	}
}

// adjustHighLatitude menerapkan penyesuaian lintang tinggi untuk Subuh dan Isya.
func (s *PrayerService) adjustHighLatitude(cfg PrayerConfig, fajr, isha, sunrise, maghrib time.Time, nightDur time.Duration) (time.Time, time.Time) {
	switch cfg.HighLatMethod {

	case HIGH_LAT_MIDDLE_OF_NIGHT:
		// Subuh & Isya masing-masing = ½ durasi malam
		half := nightDur / 2
		fajr = sunrise.Add(-half)
		isha = maghrib.Add(half)

	case HIGH_LAT_SEVENTH_OF_NIGHT:
		// Subuh = Sunrise - 1/7 malam (terakhir)
		// Isya = Maghrib + 1/7 malam (pertama)
		seventh := nightDur / 7
		fajr = sunrise.Add(-seventh)
		isha = maghrib.Add(seventh)

	case HIGH_LAT_ANGLE_BASED:
		// Porsi malam = |angle| / 60
		// Subuh: mundur dari sunrise sebesar (fajrAngle/60 * nightDur)
		// Isya: maju dari maghrib sebesar (ishaAngle/60 * nightDur)
		fajrPortion := math.Abs(cfg.FajrAngle) / 60.0
		fajr = sunrise.Add(-time.Duration(float64(nightDur) * fajrPortion))

		if cfg.IshaOffsetMin <= 0 {
			// Hanya adjust Isha jika bukan offset-based (Umm Al-Qura)
			ishaPortion := math.Abs(cfg.IshaAngle) / 60.0
			isha = maghrib.Add(time.Duration(float64(nightDur) * ishaPortion))
		}
	}

	return fajr, isha
}

// isRamadan cek apakah tanggal Masehi ini kemungkinan jatuh di bulan Ramadan.
// Menggunakan pendekatan tabular Hijriah sederhana.
func isRamadan(t time.Time) bool {
	// Julian Day
	jd := float64(t.Unix())/86400.0 + 2440587.5
	days := int(math.Floor(jd - 1948439.5))
	cycle := days / 10631
	rem := days % 10631

	leapYears := map[int]bool{
		2: true, 5: true, 7: true, 10: true, 13: true,
		16: true, 18: true, 21: true, 24: true, 26: true, 29: true,
	}

	year := cycle * 30
	for i := 1; i <= 30; i++ {
		daysInYear := 354
		if leapYears[i] {
			daysInYear = 355
		}
		if rem < daysInYear {
			year += i
			break
		}
		rem -= daysInYear
	}

	// Hitung bulan
	_ = year
	for m := 1; m <= 12; m++ {
		daysInMonth := 29
		if m%2 != 0 {
			daysInMonth = 30
		}
		if m == 12 && leapYears[year%30] {
			daysInMonth = 30
		}
		if rem < daysInMonth {
			return m == 9 // Bulan 9 = Ramadan
		}
		rem -= daysInMonth
	}

	return false
}
