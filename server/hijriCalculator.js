import SunCalc from 'suncalc';
import { solar, moonposition, angle } from 'astronomia';
import { DateTime } from 'luxon';

function angularSeparation(ra1, dec1, ra2, dec2) {
    const rad = Math.PI / 180;
    const cosTheta = Math.sin(dec1 * rad) * Math.sin(dec2 * rad) +
        Math.cos(dec1 * rad) * Math.cos(dec2 * rad) * Math.cos((ra1 - ra2) * rad);
    return Math.acos(cosTheta) / rad; // Hasil dalam derajat
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
    let step = 0.1; // Gunakan langkah 0.1 hari (~2.4 jam) untuk mendeteksi lebih akurat
    let limit = jd + 1; // Cari dalam rentang 1 hari
    let lastDiff = null;

    while (jd < limit) {
        let moonPos = getMoonPosition(jd);
        let sunPos = getSunPosition(jd);
        let diff = Math.abs(moonPos.ra - sunPos.ra);

        if (lastDiff !== null && diff > lastDiff) {
            console.log("✅ INFO: Konjungsi diperkirakan terjadi di JD:", jd - step);
            return jd - step;
        }

        lastDiff = diff;
        jd += step;
    }

    console.log("⚠️ WARNING: Konjungsi tidak ditemukan dalam rentang waktu.");
    return null;
}

function getMoonAltitude(date, lat, lon) {
    const moonPos = SunCalc.getMoonPosition(date, lat, lon);
    return moonPos.altitude * (180 / Math.PI); // Konversi dari radian ke derajat
}

function getSunAltitude(date, lat, lon) {
    const sunPos = SunCalc.getPosition(date, lat, lon);
    return sunPos.altitude * (180 / Math.PI);
}

// Fungsi untuk mengonversi tanggal Julian ke Hijriyah
function julianToHijri(jd) {
    const jdEpoch = 1948439.5; // Epoch kalender Hijriyah
    const daysSinceEpoch = Math.floor(jd - jdEpoch);

    let hijriYear = Math.floor((30 * daysSinceEpoch + 10646) / 10631);
    const startYearJD = jdEpoch + Math.floor((10631 * hijriYear) / 30);
    let remainingDays = Math.floor(jd - startYearJD);

    if (remainingDays < 0) remainingDays += 354; // Jika tahun sebelumnya

    let hijriMonth = Math.floor(remainingDays / 29.5) + 1;
    let hijriDay = Math.floor(remainingDays - (29.5 * (hijriMonth - 1))) + 1;

    // Koreksi jika bulan lebih dari 12
    if (hijriMonth > 12) {
        hijriMonth = 1;
        hijriYear++;
    }

    return { year: hijriYear, month: hijriMonth, day: hijriDay };
}

// Fungsi untuk mendapatkan tanggal Hijriyah berdasarkan lokasi dan metode perhitungan
export function getHijriDate(lat, lon, method, timezone, jd = null) {
    const nowUTC = DateTime.utc();
    const sunsetTimeUTC = SunCalc.getTimes(nowUTC.toJSDate(), lat, lon).sunset;
    const sunsetLuxon = DateTime.fromJSDate(sunsetTimeUTC, { zone: "utc" }).setZone(timezone);
    const nowLuxon = nowUTC.setZone(timezone);
    const effectiveDate = nowLuxon >= sunsetLuxon ? nowLuxon.plus({ days: 1 }) : nowLuxon;

    if (!jd) {
        jd = getJulianDate(effectiveDate);
    }

    let hijri = julianToHijri(jd);
    return { ...hijri, method, timezone, localTime: nowLuxon.toISO(), sunsetTime: sunsetLuxon.toISO() };
}

export function predictEndOfMonth(lat, lon, method, timezone, targetHijriDay = 29) {
    const nowLuxon = DateTime.now().setZone(timezone);
    const hijri = getHijriDate(lat, lon, method, timezone);

    const targetDate = nowLuxon.plus({ days: targetHijriDay - hijri.day });
    const jd = getJulianDate(targetDate);

    console.log("📅 DEBUG: Prediksi untuk Hijri Day", targetHijriDay);
    console.log("📅 DEBUG: Julian Date yang digunakan", jd);
    console.log("📅 DEBUG: Tanggal yang dianalisis:", targetDate.toISO());

    let isEndOfMonth = "30 hari";
    const moonAltitude = getMoonAltitude(targetDate.toJSDate(), lat, lon) ?? "Tidak tersedia";
    const sunAltitude = getSunAltitude(targetDate.toJSDate(), lat, lon) ?? "Tidak tersedia";

    const moonPos = getMoonPosition(jd) || { ra: 0, dec: 0 };
    const sunPos = getSunPosition(jd) || { ra: 0, dec: 0 };
    const conjunctionJD = getConjunctionTime(jd) || null;

    const conjunctionOccurred = conjunctionJD !== null && jd >= conjunctionJD;
    const elongation = getElongation(moonPos, sunPos) || "Tidak tersedia";
    const moonAgeHours = conjunctionOccurred ? (jd - conjunctionJD) * 24 : null;

    if (conjunctionOccurred) {
        console.log("⏳ INFO: Umur Bulan:", moonAgeHours.toFixed(2), "jam");
    } else {
        console.log("⏳ INFO: Konjungsi belum terjadi, umur bulan tidak diketahui.");
    }

    const memenuhiGlobal = method === 'global' || method === 'hisab';
    const memenuhiRukyat = method === 'rukyat' && moonAltitude > sunAltitude;
    const memenuhiImkanurRukyat = method === 'imkanur-rukyat' && moonAltitude >= 3 && elongation >= 6.4 && moonAgeHours !== null && moonAgeHours >= 8;

    if (conjunctionOccurred && (memenuhiGlobal || memenuhiRukyat || memenuhiImkanurRukyat)) {
        hijri.day = 1;
        hijri.month++;
        isEndOfMonth = "29 hari";
    } else {
        hijri.day = 30;
    }

    if (hijri.month > 12) {
        hijri.month = 1;
        hijri.year++;
    }

    console.log("🔍 Hasil Akhir:", { hijri });
    console.log("🌍 INFO: Lokasi -> Lat:", lat, "Lon:", lon);
    console.log("📅 INFO: Julian Date yang dihitung:", jd);
    console.log("🌙 INFO: Posisi Bulan -> RA:", moonPos.ra, "Dec:", moonPos.dec);
    console.log("☀️ INFO: Posisi Matahari -> RA:", sunPos.ra, "Dec:", sunPos.dec);
    console.log("🌙 INFO: Ketinggian Bulan (derajat):", moonAltitude);
    console.log("☀️ INFO: Ketinggian Matahari (derajat):", sunAltitude);
    console.log("🔭 INFO: Elongasi Bulan:", elongation, "derajat");
    console.log("⏳ INFO: Umur Bulan:", moonAgeHours !== null ? moonAgeHours.toFixed(2) + " jam" : "Tidak diketahui");
    console.log("📌 DEBUG: Metode yang dipakai:", method);
    console.log("📌 DEBUG: Sudah terjadi konjungsi?", conjunctionOccurred);
    console.log("📌 DEBUG: Apakah elongasi memenuhi syarat? (≥ 6.4°)", elongation >= 6.4);
    console.log("📌 DEBUG: Apakah ketinggian bulan memenuhi syarat? (≥ 3°)", moonAltitude >= 3);
    console.log("📌 DEBUG: Apakah umur bulan memenuhi syarat? (≥ 8 jam)", moonAgeHours !== null && moonAgeHours >= 8);

    return {
        hijri,
        explanation: `Bulan ini berjumlah ${isEndOfMonth} berdasarkan metode ${method}.`,
        method,
        timezone,
        localTime: targetDate.toISO(),
        moonAltitude: moonAltitude !== "Tidak tersedia" ? moonAltitude.toFixed(2) : "Tidak tersedia",
        sunAltitude: sunAltitude !== "Tidak tersedia" ? sunAltitude.toFixed(2) : "Tidak tersedia",
        elongation: elongation !== "Tidak tersedia" ? elongation.toFixed(2) : "Tidak tersedia",
        moonAge: moonAgeHours !== null ? moonAgeHours.toFixed(2) : "Tidak diketahui",
        conjunction: conjunctionOccurred
    };
}
