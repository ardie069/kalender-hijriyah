import SunCalc from "suncalc";
import { solar, moonposition } from "astronomia";
import { DateTime } from "luxon";

function angularSeparation(ra1, dec1, ra2, dec2) {
  const rad = Math.PI / 180;
  const cosTheta =
    Math.sin(dec1 * rad) * Math.sin(dec2 * rad) +
    Math.cos(dec1 * rad) * Math.cos(dec2 * rad) * Math.cos((ra1 - ra2) * rad);
  return Math.acos(cosTheta) / rad;
}

function getJulianDate(date) {
  return date.toMillis() / 86400000.0 + 2440587.5;
}

function getElongation(moon, sun) {
  return angularSeparation(moon.ra, moon.dec, sun.ra, sun.dec);
}

function getMoonPosition(jd) {
  let pos = moonposition.position(jd);
  return { ra: pos.ra, dec: pos.dec };
}

function getSunPosition(jd) {
  let pos = solar.trueEquatorial(jd); // lebih akurat dibanding apparentEquatorial
  return { ra: pos.ra, dec: pos.dec };
}

function getConjunctionTime(jd) {
  let step = 0.01; // Perbesar langkah awal untuk cakupan luas (~86 detik)
  let limit = jd + 2; // Pencarian lebih luas dalam 2 hari setelah sunset
  let minDiff = Infinity;
  let bestJD = null;

  while (jd < limit) {
    let moonPos = getMoonPosition(jd);
    let sunPos = getSunPosition(jd);
    let diff = Math.abs(moonPos.ra - sunPos.ra);

    if (diff < minDiff) {
      minDiff = diff;
      bestJD = jd;
    }

    jd += step;
  }

  // **Refinement dengan Pencarian Biner**
  let left = bestJD - step;
  let right = bestJD + step;
  let precision = 0.00001; // Refinement tinggi (~0.86 detik)

  while (right - left > precision) {
    let mid = (left + right) / 2;
    let moonPos = getMoonPosition(mid);
    let sunPos = getSunPosition(mid);
    let diff = Math.abs(moonPos.ra - sunPos.ra);

    if (diff < minDiff) {
      minDiff = diff;
      bestJD = mid;
    }

    if (moonPos.ra > sunPos.ra) {
      right = mid;
    } else {
      left = mid;
    }
  }

  return bestJD;
}

function getMoonAltitude(date, lat, lon) {
  const moonPos = SunCalc.getMoonPosition(date, lat, lon);
  return moonPos.altitude * (180 / Math.PI);
}

function getSunAltitude(date, lat, lon) {
  const sunPos = SunCalc.getPosition(date, lat, lon);
  return sunPos.altitude * (180 / Math.PI);
}

function julianToHijri(jd) {
  const jdEpoch = 1948439.5;
  const daysSinceEpoch = Math.floor(jd - jdEpoch);

  let hijriYear = Math.floor((30 * daysSinceEpoch + 10646) / 10631);
  const startYearJD = jdEpoch + Math.floor((10631 * hijriYear) / 30);
  let remainingDays = Math.floor(jd - startYearJD);

  if (remainingDays < 0) remainingDays += 354;

  let hijriMonth = Math.floor(remainingDays / 29.5) + 1;
  let hijriDay = Math.floor(remainingDays - 29.5 * (hijriMonth - 1)) + 1;

  if (hijriMonth > 12) {
    hijriMonth = 1;
    hijriYear++;
  }

  return { year: hijriYear, month: hijriMonth, day: hijriDay };
}

export function getHijriDate(lat, lon, method, timezone, jd = null) {
  const nowLocal = DateTime.local().setZone(timezone);

  // Lokasi referensi Arab Saudi untuk metode global
  const saudiLat = 21.422487;
  const saudiLon = 39.826206;
  const saudiZone = "Asia/Riyadh";

  // Lokasi acuan matahari terbenam
  const refLat = method === "global" ? saudiLat : lat;
  const refLon = method === "global" ? saudiLon : lon;
  const refZone = method === "global" ? saudiZone : timezone;

  // Hitung waktu matahari terbenam lokal
  const sunsetTodayUTC = SunCalc.getTimes(
    nowLocal.toJSDate(),
    refLat,
    refLon
  ).sunset;
  const sunsetLocal = DateTime.fromJSDate(sunsetTodayUTC).setZone(refZone);
  const sunsetJD = getJulianDate(sunsetLocal.toUTC());

  // Bandingkan dengan waktu sekarang (dalam zona referensi)
  const nowRef = nowLocal.setZone(refZone);
  const afterSunset = nowRef >= sunsetLocal;

  // Gunakan JD saat maghrib, dan tambah 1 jika waktu sudah lewat maghrib
  let effectiveJD = sunsetJD;
  if (afterSunset) {
    effectiveJD += 1;
  }

  // Konversi ke tanggal Hijriyah sementara
  let hijriDate = julianToHijri(effectiveJD);

  // Logika tanggal 29 Hijriyah (penentuan apakah digenapkan jadi 30)
  if (hijriDate.day === 29) {
    const conjunctionJD = getConjunctionTime(effectiveJD - 1);
    let conjunctionValid = false;

    if (method === "global" || method === "hisab") {
      conjunctionValid = conjunctionJD < sunsetJD;
    } else if (method === "rukyat") {
      const moonAltitude = getMoonAltitude(sunsetLocal.toJSDate(), lat, lon);
      const elongation = getElongation(
        moonposition.position(sunsetJD),
        solar.trueEquatorial(sunsetJD)
      );
      const moonAgeHours = (sunsetJD - conjunctionJD) * 24;

      conjunctionValid =
        moonAltitude >= 3 && elongation >= 6.4 && moonAgeHours >= 8;
    }

    // Jika hilal belum dianggap terlihat → genapkan jadi 30
    if (!conjunctionValid && afterSunset) {
      effectiveJD += 1;
    }
  }

  return julianToHijri(effectiveJD);
}

export function predictEndOfMonth(lat, lon, method, timezone) {
  function julianToDate(jd) {
    const J2000 = 2451545.0;
    const SECONDS_IN_DAY = 86400;
    const MILLISECONDS_IN_SECOND = 1000;

    let timestamp =
      (jd - J2000) * SECONDS_IN_DAY * MILLISECONDS_IN_SECOND +
      Date.UTC(2000, 0, 1, 12);
    return new Date(timestamp);
  }

  // Mendapatkan tanggal Hijriyah hari ini
  const nowLuxon = DateTime.now().setZone(timezone);
  const nowUTC = DateTime.utc();
  const jd = getJulianDate(nowUTC);

  const hijriToday = getHijriDate(lat, lon, method, timezone, jd);

  if (hijriToday.day !== 29) {
    return {
      today: hijriToday,
      estimatedEndOfMonth: null,
      message: "Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.",
    };
  }

  // Melanjutkan perhitungan jika tanggal Hijriyah adalah 29
  const hijriTarget = {
    day: 29,
    month: hijriToday.month,
    year: hijriToday.year,
  };
  const daysTo29 = hijriToday.day <= 29 ? 29 - hijriToday.day : 0;
  const estimatedDate = nowLuxon.plus({ days: daysTo29 }).startOf("day");

  // Gunakan referensi Mekkah untuk waktu konjungsi jika metode global
  const refLat = method === "global" ? 21.422487 : lat;
  const refLon = method === "global" ? 39.826206 : lon;
  const refZone = method === "global" ? "Asia/Riyadh" : timezone;

  // Sunset waktu di lokasi pengguna
  const sunsetTargetUTC = SunCalc.getTimes(
    estimatedDate.toJSDate(),
    lat,
    lon
  ).sunset;
  const sunsetTarget = DateTime.fromJSDate(sunsetTargetUTC).setZone(timezone);
  const sunsetJD = getJulianDate(sunsetTarget.toUTC());

  // Sunset waktu di Mekkah untuk metode global
  const sunsetTargetMekkahUTC = SunCalc.getTimes(
    estimatedDate.toJSDate(),
    21.422487,
    39.826206
  ).sunset;
  const sunsetTargetMekkah = DateTime.fromJSDate(sunsetTargetMekkahUTC).setZone(
    "Asia/Riyadh"
  );
  const sunsetJDGlobal = getJulianDate(sunsetTargetMekkah.toUTC());

  // Hitung konjungsi berdasarkan referensi global jika metode global
  const conjunctionJD = getConjunctionTime(
    getJulianDate(
      DateTime.fromJSDate(
        SunCalc.getTimes(estimatedDate.toJSDate(), refLat, refLon).sunset
      )
        .setZone(refZone)
        .toUTC()
    )
  );
  const conjunctionTime = conjunctionJD
    ? DateTime.fromJSDate(julianToDate(conjunctionJD)).setZone(refZone)
    : null;

  const toleranceJD = 0.02083; // 30 menit dalam hari

  // Cek konjungsi sebelum matahari terbenam di Mekkah (untuk global) atau di lokasi pengguna (untuk Hisab)
  const conjunctionBeforeSunsetGlobal =
    conjunctionJD !== null && conjunctionJD < sunsetJDGlobal - toleranceJD;
  const conjunctionBeforeSunsetUser =
    conjunctionJD !== null && conjunctionJD < sunsetJD - toleranceJD;

  // Hitung posisi bulan dan matahari di lokasi pengguna
  const moonAltitude = getMoonAltitude(sunsetTarget.toJSDate(), lat, lon);
  const sunAltitude = getSunAltitude(sunsetTarget.toJSDate(), lat, lon);
  const moonPos = getMoonPosition(sunsetJD) || { ra: 0, dec: 0 };
  const sunPos = getSunPosition(sunsetJD) || { ra: 0, dec: 0 };
  const elongation = getElongation(moonPos, sunPos) || NaN;
  const moonAgeHours =
    conjunctionJD !== null ? (sunsetJD - conjunctionJD) * 24 : NaN;

  // Validasi untuk metode Global dan Rukyat
  const memenuhiGlobal = method === "global" && conjunctionBeforeSunsetGlobal;
  const memenuhiHisab = method === "hisab" && conjunctionBeforeSunsetUser;
  const memenuhiRukyat =
    method === "rukyat" &&
    moonAltitude >= 3 &&
    elongation >= 6.4 &&
    moonAgeHours >= 8;

  let isEndOfMonth = "30 hari";
  let hijriEnd = { day: 30, month: hijriTarget.month, year: hijriTarget.year };

  // Tentukan apakah bulan baru dimulai berdasarkan konjungsi
  if (!conjunctionBeforeSunsetGlobal && !conjunctionBeforeSunsetUser) {
    isEndOfMonth = "30 hari";
  } else if (memenuhiGlobal || memenuhiHisab || memenuhiRukyat) {
    hijriEnd = { day: 1, month: hijriTarget.month + 1, year: hijriTarget.year };
    isEndOfMonth = "29 hari";
  }

  // Jika bulan lebih dari 12 (Zulhijjah), maka harus roll over ke tahun baru
  if (hijriEnd.month > 12) {
    hijriEnd.month = 1;
    hijriEnd.year++;
  }

  // Tentukan tanggal mulai bulan baru
  const nextHijriStartDate = sunsetTarget.plus({ days: 1 }).toISODate();

  // Parsing tanggal Hijriyah ke format Gregorian
  const nextHijriStartDateParsed =
    DateTime.fromISO(nextHijriStartDate).setZone(timezone);

  // Return hasil perhitungan hanya jika tanggal Hijriyah adalah 29
  return {
    today: hijriToday,
    estimatedEndOfMonth: {
      hijri: hijriEnd,
      explanation: `Bulan ini berjumlah ${isEndOfMonth} berdasarkan metode ${method}.`,
      moonAltitude: isNaN(moonAltitude)
        ? "Tidak tersedia"
        : moonAltitude.toFixed(2),
      sunAltitude: isNaN(sunAltitude)
        ? "Tidak tersedia"
        : sunAltitude.toFixed(2),
      elongation: isNaN(elongation) ? "Tidak tersedia" : elongation.toFixed(2),
      moonAge: isNaN(moonAgeHours)
        ? "Tidak diketahui"
        : moonAgeHours.toFixed(2),
      conjunctionBeforeSunsetGlobal: conjunctionBeforeSunsetGlobal,
      conjunctionBeforeSunsetUser: conjunctionBeforeSunsetUser,
      conjunctionTime: conjunctionTime
        ? conjunctionTime.toFormat("yyyy-MM-dd HH:mm:ss")
        : "Tidak diketahui",
    },
    estimatedStartOfMonth:
      hijriEnd.month === 12
        ? {
            hijri: { day: 1, month: 1, year: hijriEnd.year + 1 },
            gregorian: {
              day: nextHijriStartDateParsed.day,
              month: nextHijriStartDateParsed.month,
              year: nextHijriStartDateParsed.year,
            },
          }
        : {
            hijri: { day: 1, month: hijriEnd.month + 1, year: hijriEnd.year },
            gregorian: {
              day: nextHijriStartDateParsed.day,
              month: nextHijriStartDateParsed.month,
              year: nextHijriStartDateParsed.year,
            },
          },
  };
}
