package calendar

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
)

// KHGT Parameters
const (
	AmericasLonStart = -30.0  // Mulai dari pesisir Timur (Greenland/Brasil)
	AmericasLonEnd   = -165.0 // Berakhir di Alaska/Hawaii
	NZLat            = -41.0  // Wellington, Selandia Baru
	NZLon            = 174.0  // Garis penanggalan internasional
)

func (l *Logic) ScanGlobalUGHC(targetDateUTC time.Time, ijtimaTime time.Time) bool {
	// 1. DEADLINE: 24:00 UTC
	deadline := time.Date(targetDateUTC.Year(), targetDateUTC.Month(), targetDateUTC.Day(), 23, 59, 59, 0, time.UTC)

	// A. TEST SEBELUM 24:00 UTC (SELURUH DUNIA)
	// Kita bisa membatasi pencarian di limit Barat (Benua Amerika) karena jika hilal
	// paling timur saja wujud, barat pasti wujud. Posisi paling optimal wujudul hilal
	// sebelum 24:00 UTC adalah di pesisir Barat Amerika.

	foundBeforeMidnight := false

	// Scan parameter 1: Terpenuhi di mana saja di dunia sebelum 24.00 UTC
	// Karena elongasi bertambah seriring waktu (dan time berjalan ke barat),
	// kondisi paling maksimal sebelum 24:00 UTC selalu berada di titik bujur terbarat (Americas)
	for lon := AmericasLonStart; lon >= AmericasLonEnd; lon -= 1.0 {
		for lat := 60.0; lat >= -60.0; lat -= 10.0 {
			sunsetTime, err := l.Astro.GetSunset(targetDateUTC, lat, lon)

			if err != nil || sunsetTime.After(deadline) {
				sunsetTime, _ = l.Astro.GetSunset(targetDateUTC.AddDate(0, 0, -1), lat, lon)
			}

			if sunsetTime.Before(ijtimaTime) || sunsetTime.After(deadline) {
				continue
			}

			et, _ := astronomy.Str2et(sunsetTime.Format(astronomy.TimeFormat))
			sunPos, _ := l.Manager.GetGeocentric(astronomy.Sun, et, "IAU_EARTH")
			moonPos, _ := l.Manager.GetGeocentric(astronomy.Moon, et, "IAU_EARTH")

			alt, elong := l.Astro.CalculateGeocentricParams(sunPos, moonPos, lat, lon)

			if IsTurkey2016(alt, elong) {
				foundBeforeMidnight = true
				break
			}
		}
		if foundBeforeMidnight {
			break
		}
	}

	if foundBeforeMidnight {
		return true // Langsung mulai bulan baru jika syarat utama terpenuhi
	}

	// B. TEST PASCA 24:00 UTC (Tengah malam terlewati)
	// Apabila kriteria Imkanu Rukyat terpenuhi SETELAH 24:00 UTC, bulan baru tetap dimulai APABILA:
	// Butir 1: IJTIMA Selandia Baru sebelum fajar
	// Butir 2: Kriteria terpenuhi di benua Amerika (yang mana sunsetnya memang SETELAH 24:00 UTC!)

	foundAfterMidnightInAmericas := false
	for lon := AmericasLonStart; lon >= AmericasLonEnd; lon -= 1.0 {
		for lat := 60.0; lat >= -60.0; lat -= 10.0 {
			// Kita ambil sunset HARI INI yang justru melewati 24:00 UTC (artinya jatuh di UTC besoknya)
			// Ini normal untuk wilayah Amerika saat targetDateUTC siang/sore
			sunsetTime, err := l.Astro.GetSunset(targetDateUTC, lat, lon)
			if err != nil {
				continue
			}

			// Kita hanya peduli yang sunsetnya SETELAH 24:00 UTC
			if !sunsetTime.After(deadline) {
				// Kalau sebelum deadline, harusnya udah kena break di check A.
				// Tapi buat make sure, kita lanjut ambil sunset besoknya
				sunsetTime, _ = l.Astro.GetSunset(targetDateUTC.AddDate(0, 0, 1), lat, lon)
			}

			if sunsetTime.Before(ijtimaTime) {
				continue
			}

			et, _ := astronomy.Str2et(sunsetTime.Format(astronomy.TimeFormat))
			sunPos, _ := l.Manager.GetGeocentric(astronomy.Sun, et, "IAU_EARTH")
			moonPos, _ := l.Manager.GetGeocentric(astronomy.Moon, et, "IAU_EARTH")

			alt, elong := l.Astro.CalculateGeocentricParams(sunPos, moonPos, lat, lon)

			if IsTurkey2016(alt, elong) {
				foundAfterMidnightInAmericas = true
				break
			}
		}
		if foundAfterMidnightInAmericas {
			break
		}
	}

	// Eksekusi Butir Pengecualian Pasca-Tengah Malam
	if foundAfterMidnightInAmericas {
		// Butir 1: Terjadi di Benua Amerika (Syaratnya secara inheren terpenuhi oleh `foundAfterMidnightInAmericas`)
		// Namun juga ada Butir: "telah terpenuhi di tempat mana pun ... DAN ijtimak Selandia Baru sebelum fajar"
		// Mengacu rumusan Muhammadiyah KHGT:
		// Jika terpenuhi di Amerika, kita BISA MULAI bulan baru asalkan:
		// 1. Secara umum Ijtima di Selandia Baru jatuh sebelum Fajar.
		fajrNZ, _ := l.Astro.GetFajr(targetDateUTC, NZLat, NZLon)

		if ijtimaTime.Before(fajrNZ) {
			return true // Butir 1 terpenuhi
		}

		// Atau klausa Amerika Serikat khusus (Butir 2 Terjadi di wilayah Amerika)
		// Wait, parameter aslinya "Parameter di atas pada butir 1 terjadi di wilayah daratan Benua Amerika."
		// Ini artinya jika TERCAPAI di BENUA AMERIKA, maka SAH, tidak peduli NZ Fajr!
		// Mari kita baca prompt:
		// "Apabila parameter di atas telah terpenuhi di suatu tempat mana pun di dunia dan ijtimak di Selandia Baru terjadi sebelum fajar."
		// "Parameter di atas pada butir 1 terjadi di wilayah daratan Benua Amerika."
		// Interpretasi:
		// Syarat A: Tempat mana pun & Ijtima NZ sebelum Fajr.
		// ATAU Syarat B: Terjadi di Amerika.
		return true // Karena ditemukan di daratan Benua Amerika
	}

	return false
}
