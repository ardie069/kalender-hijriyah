import SunCalc from 'suncalc';
import { solar, moonposition } from 'astronomia';
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
    let pos = solar.apparentEquatorial(jd);
    return { ra: pos.ra, dec: pos.dec };
}

function getConjunctionTime(jd) {
    let step = 0.1;
    let limit = jd + 1;
    let lastDiff = null;

    while (jd < limit) {
        let moonPos = getMoonPosition(jd);
        let sunPos = getSunPosition(jd);
        let diff = Math.abs(moonPos.ra - sunPos.ra);

        if (lastDiff !== null && diff > lastDiff) {
            return jd - step;
        }

        lastDiff = diff;
        jd += step;
    }

    return null;
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

    const yesterdayLocal = nowLocal.minus({ days: 1 });
    
    const sunsetTodayUTC = SunCalc.getTimes(nowLocal.toJSDate(), lat, lon).sunset;
    const sunsetYesterdayUTC = SunCalc.getTimes(yesterdayLocal.toJSDate(), lat, lon).sunset;

    const sunsetToday = DateTime.fromJSDate(sunsetTodayUTC).setZone(timezone);
    const sunsetYesterday = DateTime.fromJSDate(sunsetYesterdayUTC).setZone(timezone);

    let lastSunset = nowLocal >= sunsetToday ? sunsetToday : sunsetYesterday;

    let effectiveDate = nowLocal;
    if (nowLocal >= lastSunset) {
        effectiveDate = nowLocal.plus({ days: 1 }).startOf('day');
    }

    if (!jd) {
        jd = getJulianDate(effectiveDate);
    }

    let hijri = julianToHijri(jd);
    return {
        ...hijri,
        method,
        timezone,
        localTime: nowLocal.toISO(),
        sunsetYesterday: sunsetYesterday.toISO(),
        sunsetToday: sunsetToday.toISO(),
        effectiveDate: effectiveDate.toISO()
    };
}

export function predictEndOfMonth(lat, lon, method, timezone, targetHijriDay = 29) {
    const nowLuxon = DateTime.now().setZone(timezone);
    const hijriToday = getHijriDate(lat, lon, method, timezone);

    // Estimasi untuk tanggal 29
    const targetDate = nowLuxon.plus({ days: targetHijriDay - hijriToday.day });
    const jd = getJulianDate(targetDate);

    const moonAltitude = getMoonAltitude(targetDate.toJSDate(), lat, lon) ?? "Tidak tersedia";
    const sunAltitude = getSunAltitude(targetDate.toJSDate(), lat, lon) ?? "Tidak tersedia";
    const moonPos = getMoonPosition(jd) || { ra: 0, dec: 0 };
    const sunPos = getSunPosition(jd) || { ra: 0, dec: 0 };
    const conjunctionJD = getConjunctionTime(jd) || null;

    const conjunctionOccurred = conjunctionJD !== null && jd >= conjunctionJD;
    const elongation = getElongation(moonPos, sunPos) || "Tidak tersedia";
    const moonAgeHours = conjunctionOccurred ? (jd - conjunctionJD) * 24 : null;

    const memenuhiGlobal = method === 'global' || method === 'hisab';
    const memenuhiRukyat = method === 'rukyat' && moonAltitude > sunAltitude;
    const memenuhiImkanurRukyat = method === 'imkanur-rukyat' && moonAltitude >= 3 && elongation >= 6.4 && moonAgeHours !== null && moonAgeHours >= 8;

    let isEndOfMonth = "30 hari";
    let hijriEnd = { ...hijriToday };

    if (conjunctionOccurred && (memenuhiGlobal || memenuhiRukyat || memenuhiImkanurRukyat)) {
        hijriEnd.day = 1;
        hijriEnd.month++;
        isEndOfMonth = "29 hari";
    } else {
        hijriEnd.day = 30;
    }

    if (hijriEnd.month > 12) {
        hijriEnd.month = 1;
        hijriEnd.year++;
    }

    return {
        today: hijriToday,
        estimatedEndOfMonth: {
            hijri: hijriEnd,
            explanation: `Bulan ini berjumlah ${isEndOfMonth} berdasarkan metode ${method}.`,
            moonAltitude: moonAltitude !== "Tidak tersedia" ? moonAltitude.toFixed(2) : "Tidak tersedia",
            sunAltitude: sunAltitude !== "Tidak tersedia" ? sunAltitude.toFixed(2) : "Tidak tersedia",
            elongation: elongation !== "Tidak tersedia" ? elongation.toFixed(2) : "Tidak tersedia",
            moonAge: moonAgeHours !== null ? moonAgeHours.toFixed(2) : "Tidak diketahui",
            conjunction: conjunctionOccurred
        }
    };
}
