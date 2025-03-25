import SunCalc from 'suncalc';
import { julian, moonposition, solar } from 'astronomia';
import { DateTime } from 'luxon';

/**
 * Konversi Gregorian ke Julian Day
 */
function gregorianToJulian(year, month, day) {
    return julian.CalendarGregorianToJD(year, month, day);
}

/**
 * Menghitung fase bulan (konjungsi)
 */
function getMoonPhase(jd) {
    return moonposition.phase(jd);
}

/**
 * Menghitung posisi Bulan
 */
function getMoonPosition(jd) {
    return moonposition.position(jd);
}

/**
 * Menghitung posisi Matahari
 */
function getSunPosition(jd) {
    return solar.apparentVSOP87(jd);
}

/**
 * Konversi Julian Day ke Hijriyah
 */
function julianToHijri(jd) {
    const jdEpoch = 1948439.5; // Julian Day untuk 1 Muharram 1H (16 Juli 622M)
    const daysSinceEpoch = Math.floor(jd - jdEpoch);
    const hijriYear = Math.floor((30 * daysSinceEpoch + 10646) / 10631);

    // Menghitung ulang Julian Day untuk 1 Muharram tahun tersebut
    const startYearJD = jdEpoch + Math.floor((10631 * hijriYear) / 30);
    let remainingDays = Math.floor(jd - startYearJD);

    if (remainingDays < 0) remainingDays += 354;

    let hijriMonth = Math.ceil(remainingDays / 29.5);
    hijriMonth = Math.max(1, Math.min(12, hijriMonth));

    let hijriDay = Math.floor(remainingDays - ((hijriMonth - 1) * 29.5)) + 1;
    hijriDay = Math.max(1, Math.min(30, hijriDay));

    return { year: hijriYear, month: hijriMonth, day: hijriDay };
}

/**
 * Ambil tanggal Hijriyah berdasarkan lokasi pengguna dan metode perhitungan
 */
function getHijriDate(lat, lon, method, timezone) {
    const nowUTC = DateTime.utc();
    const sunsetTimeUTC = SunCalc.getTimes(nowUTC.toJSDate(), lat, lon).sunset;
    const sunsetLuxon = DateTime.fromJSDate(sunsetTimeUTC, { zone: "utc" }).setZone(timezone);
    const nowLuxon = nowUTC.setZone(timezone);

    console.log(`⏳ Sekarang UTC    : ${nowUTC.toISO()}`);
    console.log(`⏳ Sekarang (Local): ${nowLuxon.toISO()}`);
    console.log(`🌅 Matahari terbenam (Local): ${sunsetLuxon.toISO()}`);

    // Gunakan hari berikutnya jika sudah melewati matahari terbenam
    const effectiveDate = nowLuxon >= sunsetLuxon ? nowLuxon.plus({ days: 1 }) : nowLuxon;
    let jd = gregorianToJulian(effectiveDate.year, effectiveDate.month, effectiveDate.day);

    console.log(`📅 Julian Day: ${jd}`);

    let hijri = julianToHijri(jd);
    console.log(`🕌 Hijriyah awal: ${JSON.stringify(hijri)}`);

    if (hijri.day === 29) {
        const moonPos = getMoonPosition(jd);
        const sunPos = getSunPosition(jd);
        const moonPhase = getMoonPhase(jd);
        const conjunction = moonPhase < 0.5;
        const moonAltitude = moonPos.dec;
        const sunAltitude = sunPos.dec;
        const elongation = Math.abs(moonPos.ra - sunPos.ra);

        console.log(`🌙 Posisi Bulan: altitude=${moonAltitude}, elongation=${elongation}`);
        console.log(`☀️ Posisi Matahari: altitude=${sunAltitude}`);
        console.log(`🔭 Konjungsi terjadi: ${conjunction}`);

        if ((method === 'global' || method === 'hisab') && conjunction && moonAltitude > sunAltitude) {
            hijri.day = 1;
            hijri.month++;
        } else if (method === 'rukyat' && conjunction && moonAltitude > sunAltitude) {
            hijri.day = 1;
            hijri.month++;
        } else if (method === 'imkanur-rukyat' && conjunction && moonAltitude >= 3 && elongation >= 6.4) {
            hijri.day = 1;
            hijri.month++;
        } else {
            hijri.day = 30;
        }

        // Pastikan bulan tidak melebihi 12
        if (hijri.month > 12) {
            hijri.month = 1;
            hijri.year++;
        }
    }

    console.log(`📅 Hijriyah akhir: ${JSON.stringify(hijri)}`);
    return { 
        ...hijri, 
        method, 
        timezone, 
        localTime: nowLuxon.toISO(),
        sunsetTime: sunsetLuxon.toISO()
    };
}

export { getHijriDate };