import SunCalc from 'suncalc';
import { solar, moonposition, conjunction } from 'astronomia';
import { DateTime } from 'luxon';

function angularSeparation(ra1, dec1, ra2, dec2) {
    const rad = Math.PI / 180;
    const cosTheta = Math.sin(dec1 * rad) * Math.sin(dec2 * rad) +
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

    while ((right - left) > precision) {
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
    let hijriDay = Math.floor(remainingDays - (29.5 * (hijriMonth - 1))) + 1;

    if (hijriMonth > 12) {
        hijriMonth = 1;
        hijriYear++;
    }

    return { year: hijriYear, month: hijriMonth, day: hijriDay };
}

export function getHijriDate(lat, lon, method, timezone, jd = null) {
    const nowLocal = DateTime.local().setZone(timezone);
    if (!jd) {
        jd = getJulianDate(nowLocal.toUTC());
    }

    // Koordinat Arab Saudi untuk metode global
    const saudiLat = 21.422487;
    const saudiLon = 39.826206;
    const saudiZone = "Asia/Riyadh";

    let sunsetRefLat = method === "global" ? saudiLat : lat;
    let sunsetRefLon = method === "global" ? saudiLon : lon;
    let sunsetRefZone = method === "global" ? saudiZone : timezone;

    // Ambil waktu matahari terbenam berdasarkan metode
    const sunsetTodayUTC = SunCalc.getTimes(nowLocal.toJSDate(), sunsetRefLat, sunsetRefLon).sunset;
    const sunsetLocal = DateTime.fromJSDate(sunsetTodayUTC).setZone(sunsetRefZone);
    const sunsetJD = getJulianDate(sunsetLocal.toUTC());

    // Hitung waktu konjungsi
    const conjunctionJD = getConjunctionTime(sunsetJD);

    let effectiveJD = jd;

    // Jika metode global, gunakan referensi Arab Saudi
    if (method === "global" && conjunctionJD < sunsetJD) {
        effectiveJD = sunsetJD + 1;
    } else if (nowLocal >= sunsetLocal) {
        if (method === "hisab" && conjunctionJD < sunsetJD) {
            effectiveJD += 1;
        } else if (method === "rukyat") {
            let moonAltitude = getMoonAltitude(sunsetLocal.toJSDate(), lat, lon);
            let elongation = getElongation(moonposition.position(sunsetJD), solar.trueEquatorial(sunsetJD));
            let moonAgeHours = (sunsetJD - conjunctionJD) * 24;

            if (moonAltitude >= 3 && elongation >= 6.4 && moonAgeHours >= 8) {
                effectiveJD += 1;
            }
        }
    }

    return julianToHijri(effectiveJD);
}

export function predictEndOfMonth(lat, lon, method, timezone) {
    function julianToDate(jd) {
        const J2000 = 2451545.0;
        const SECONDS_IN_DAY = 86400;
        const MILLISECONDS_IN_SECOND = 1000;

        let timestamp = (jd - J2000) * SECONDS_IN_DAY * MILLISECONDS_IN_SECOND + Date.UTC(2000, 0, 1, 12);
        return new Date(timestamp);
    }

    const nowLuxon = DateTime.now().setZone(timezone);
    const nowUTC = DateTime.utc();
    const jd = getJulianDate(nowUTC);

    const hijriToday = getHijriDate(lat, lon, method, timezone, jd);
    const hijriTarget = { day: 29, month: hijriToday.month, year: hijriToday.year };
    const daysTo29 = hijriToday.day <= 29 ? 29 - hijriToday.day : 0;
    const estimatedDate = nowLuxon.plus({ days: daysTo29 }).startOf('day');

    // Gunakan referensi Arab Saudi untuk metode global
    const sunsetRefLat = method === "global" ? 21.422487 : lat;
    const sunsetRefLon = method === "global" ? 39.826206 : lon;
    const sunsetRefZone = method === "global" ? "Asia/Riyadh" : timezone;

    const sunsetTargetUTC = SunCalc.getTimes(estimatedDate.toJSDate(), sunsetRefLat, sunsetRefLon).sunset;
    const sunsetTarget = DateTime.fromJSDate(sunsetTargetUTC).setZone(sunsetRefZone);
    const sunsetJD = getJulianDate(sunsetTarget.toUTC());

    // Hitung konjungsi
    const conjunctionJD = getConjunctionTime(sunsetJD);
    const conjunctionTime = conjunctionJD
        ? DateTime.fromJSDate(julianToDate(conjunctionJD)).setZone(sunsetRefZone)
        : null;

    const toleranceJD = 0.02083; // 30 menit dalam hari
    const conjunctionBeforeSunset = conjunctionJD !== null && conjunctionJD < (sunsetJD - toleranceJD);

    // Hitung posisi bulan dan matahari
    const moonAltitude = getMoonAltitude(sunsetTarget.toJSDate(), sunsetRefLat, sunsetRefLon);
    const sunAltitude = getSunAltitude(sunsetTarget.toJSDate(), sunsetRefLat, sunsetRefLon);
    const moonPos = getMoonPosition(sunsetJD) || { ra: 0, dec: 0 };
    const sunPos = getSunPosition(sunsetJD) || { ra: 0, dec: 0 };
    const elongation = getElongation(moonPos, sunPos) || NaN;
    const moonAgeHours = conjunctionJD !== null ? (sunsetJD - conjunctionJD) * 24 : NaN;

    const memenuhiGlobal = method === "global" || (method === "hisab" && conjunctionBeforeSunset);
    const memenuhiRukyat =
        method === "rukyat" &&
        moonAltitude >= 3 &&
        elongation >= 6.4 &&
        moonAgeHours >= 8;

    let isEndOfMonth = "30 hari";
    let hijriEnd = { day: 30, month: hijriTarget.month, year: hijriTarget.year };

    if (!conjunctionBeforeSunset) {
        isEndOfMonth = "30 hari";
    } else if (memenuhiGlobal || memenuhiRukyat) {
        hijriEnd = { day: 1, month: hijriTarget.month + 1, year: hijriTarget.year };
        isEndOfMonth = "29 hari";
    }

    if (hijriEnd.month > 12) {
        hijriEnd.month = 1;
        hijriEnd.year++;
    }

    const nextHijriStartDate = sunsetTarget.plus({ days: 1 }).toISODate();

    return {
        today: hijriToday,
        estimatedEndOfMonth: {
            hijri: hijriEnd,
            explanation: `Bulan ini berjumlah ${isEndOfMonth} berdasarkan metode ${method}.`,
            moonAltitude: isNaN(moonAltitude) ? "Tidak tersedia" : moonAltitude.toFixed(2),
            sunAltitude: isNaN(sunAltitude) ? "Tidak tersedia" : sunAltitude.toFixed(2),
            elongation: isNaN(elongation) ? "Tidak tersedia" : elongation.toFixed(2),
            moonAge: isNaN(moonAgeHours) ? "Tidak diketahui" : moonAgeHours.toFixed(2),
            conjunctionBeforeSunset: conjunctionBeforeSunset,
            conjunctionTime: conjunctionTime
                ? conjunctionTime.toFormat("yyyy-MM-dd HH:mm:ss")
                : "Tidak diketahui",
        },
        estimatedStartOfMonth: {
            hijri: { day: 1, month: hijriEnd.month, year: hijriEnd.year },
            gregorian: nextHijriStartDate,
        },
    };
}
