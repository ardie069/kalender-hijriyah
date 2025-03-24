import SunCalc from 'suncalc';
import { julian, moonposition, solar } from 'astronomia';

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

    if (remainingDays < 0) {
        remainingDays += 354; // Koreksi jika kurang dari nol
    }

    let hijriMonth = Math.ceil(remainingDays / 29.5);
    hijriMonth = Math.max(1, Math.min(12, hijriMonth));

    let hijriDay = Math.floor(remainingDays - ((hijriMonth - 1) * 29.5)) + 1;
    hijriDay = Math.max(1, Math.min(30, hijriDay));

    return { year: hijriYear, month: hijriMonth, day: hijriDay };
}

/**
 * Ambil tanggal Hijriyah berdasarkan lokasi pengguna dan metode perhitungan
 */
function getHijriDate(lat, lon, method = 'global') {
    const now = new Date();
    const sunsetTime = SunCalc.getTimes(now, lat, lon).sunset;

    console.log(`⏳ Sekarang: ${now.toISOString()}`);
    console.log(`🌅 Matahari terbenam: ${sunsetTime.toISOString()}`);

    let jd = gregorianToJulian(now.getFullYear(), now.getMonth() + 1, now.getDate());

    if (now.getTime() >= sunsetTime.getTime()) {
        jd += 1; // Geser ke hari berikutnya jika sudah melewati matahari terbenam
        console.log(`🔄 Sudah melewati matahari terbenam, Julian Day bertambah: ${jd}`);
    } else {
        console.log(`⏳ Belum melewati matahari terbenam, Julian Day tetap: ${jd}`);
    }

    let hijri = julianToHijri(jd);
    console.log(`📅 Julian Day: ${jd}`);
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

        if (method === 'global' || method === 'hisab') {
            if (conjunction && moonAltitude > sunAltitude) {
                hijri.day = 1;
                hijri.month++;
            } else {
                hijri.day = 30;
            }
        } else if (method === 'rukyat') {
            if (conjunction && moonAltitude > sunAltitude) {
                hijri.day = 1;
                hijri.month++;
            } else {
                hijri.day = 30;
            }
        } else if (method === 'imkanur-rukyat') {
            if (conjunction && moonAltitude >= 3 && elongation >= 6.4) {
                hijri.day = 1;
                hijri.month++;
            } else {
                hijri.day = 30;
            }
        }

        // Pastikan bulan tidak melebihi 12
        if (hijri.month > 12) {
            hijri.month = 1;
            hijri.year++;
        }
    }

    console.log(`📅 Hijriyah akhir: ${JSON.stringify(hijri)}`);
    return { ...hijri, method };
}

export { getHijriDate };
