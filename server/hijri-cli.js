#!/usr/bin/env node

import { DateTime } from "luxon";
import SunCalc from "suncalc";
import readlineSync from "readline-sync";
import { solar, planetposition } from "astronomia";
import vsop87Bearth from "astronomia/data/vsop87Bearth";

// Koordinat Mekkah untuk metode global
const MECCA_LAT = 21.422487;
const MECCA_LON = 39.826206;
const MECCA_ZONE = "Asia/Riyadh";

async function main() {
    console.log("=== CLI Penanggalan Hijriyah ===");

    // Pilih metode (global, hisab, rukyat)
    const method = readlineSync.question("Pilih metode:\n1. Global\n2. Hisab\n3. Rukyat\nMasukkan nomor pilihan (1/2/3): ");

    // Tentukan lokasi berdasarkan metode
    let lat, lon, timezone;
    if (method === "1") {
        lat = MECCA_LAT;
        lon = MECCA_LON;
        timezone = MECCA_ZONE;
        console.log("Metode Global dipilih. Lokasi diatur ke Mekkah.");
    } else {
        lat = parseFloat(readlineSync.question("Masukkan latitude: "));
        lon = parseFloat(readlineSync.question("Masukkan longitude: "));
        timezone = readlineSync.question("Masukkan timezone (misal: Asia/Jakarta): ");
    }

    // Pilih tanggal opsional
    let dateInput = readlineSync.question("Masukkan tanggal (YYYY-MM-DD) atau tekan Enter untuk hari ini: ");
    let now = dateInput ? DateTime.fromISO(dateInput, { zone: timezone }) : DateTime.now().setZone(timezone);

    // Hitung waktu matahari terbenam
    const sunsetTime = SunCalc.getTimes(now.toJSDate(), lat, lon).sunset;
    const sunsetLocal = DateTime.fromJSDate(sunsetTime).setZone(timezone);
    console.log(`Waktu matahari terbenam: ${sunsetLocal.toFormat("yyyy-MM-dd HH:mm:ss")}`);

    // Hitung waktu konjungsi (perkiraan sederhana dengan astronomia)
    const conjunctionTime = findConjunction(now);
    console.log(`Perkiraan waktu konjungsi: ${conjunctionTime.toFormat("yyyy-MM-dd HH:mm:ss")}`);

    // Hitung usia bulan
    const moonAge = now.diff(conjunctionTime, "hours").hours;
    console.log(`Usia Bulan: ${moonAge.toFixed(2)} jam`);

    // Hitung ketinggian bulan saat matahari terbenam
    const moonPos = SunCalc.getMoonPosition(sunsetTime, lat, lon);
    const moonAltitude = moonPos.altitude * (180 / Math.PI);
    console.log(`Ketinggian Bulan: ${moonAltitude.toFixed(2)}°`);

    // Hitung elongasi
    const sunPos = SunCalc.getPosition(sunsetTime, lat, lon);
    const elongation = getElongation(moonPos, sunPos);
    console.log(`Elongasi Bulan: ${elongation.toFixed(2)}°`);
}

// Fungsi untuk menemukan waktu konjungsi sederhana
function findConjunction(now) {
    let date = now;
    let earth = new planetposition.Planet(vsop87Bearth);
    for (let i = 0; i < 30; i++) {
        let sunEq = solar.apparentVSOP87(earth, date.toJSDate()); // Posisi matahari
        let moonPos = SunCalc.getMoonPosition(date.toJSDate(), MECCA_LAT, MECCA_LON); // Posisi bulan
        
        // Menghitung elongasi bulan terhadap matahari (dalam derajat)
        const elongation = Math.abs(moonPos.azimuth - sunEq.ra);

        // Memeriksa apakah bulan sudah dalam fase sabit tua
        if (elongation > 0 && elongation <= 45) {
            // Menghitung ketinggian bulan untuk memastikan bulan berada dekat dengan matahari
            if (Math.abs(moonPos.altitude) < 5) { // Ketinggian bulan dekat dengan matahari
                return date;
            }
        }
        
        // Periksa jika elongasi sudah sangat kecil (konjungsi)
        if (elongation < 0.5) {  // Elongasi < 0.5° menunjukkan konjungsi
            if (Math.abs(moonPos.altitude) < 5) { // Ketinggian bulan sangat rendah (hampir konjungsi)
                return date;
            }
        }
        
        date = date.plus({ hours: 2 });
    }
    return now; // Jika konjungsi tidak ditemukan dalam 30 jam, kembalikan waktu saat ini
}

// Fungsi untuk menghitung elongasi
function getElongation(moonPos, sunPos) {
    return Math.sqrt((moonPos.azimuth - sunPos.azimuth) ** 2 + (moonPos.altitude - sunPos.altitude) ** 2) * (180 / Math.PI);
}

// Jalankan CLI
main();
