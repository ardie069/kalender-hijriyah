import SunCalc from "suncalc";
import { getHijriDate } from "../server/hijriCalculator.js";

export default function handler(req, res) {
    const { lat, lon, method, timezone } = req.query;

    if (!lat || !lon || !method || !timezone) {
        return res.status(400).json({ error: "Latitude, Longitude, Method, dan Timezone diperlukan!" });
    }

    try {
        console.log("📡 API Request:", { lat, lon, method, timezone });

        // Konversi lat & lon ke float
        const latitude = parseFloat(lat);
        const longitude = parseFloat(lon);

        // Ambil waktu sekarang dalam UTC, lalu konversi ke waktu lokal pengguna
        const nowUTC = new Date();
        const nowLocal = new Date(nowUTC.toLocaleString("en-US", { timeZone: timezone }));

        // Hitung waktu matahari terbenam di lokasi pengguna
        const sunsetTime = SunCalc.getTimes(nowLocal, latitude, longitude).sunset;
        
        console.log(`⏳ Sekarang UTC: ${nowUTC.toISOString()}`);
        console.log(`🌅 Matahari terbenam di lokasi pengguna: ${sunsetTime.toISOString()}`);

        // Pastikan hijriDate dihitung berdasarkan waktu setelah matahari terbenam
        const effectiveTime = nowLocal >= sunsetTime ? new Date(nowLocal.getTime() + 86400000) : nowLocal;
        
        // Panggil fungsi getHijriDate dengan waktu yang sudah disesuaikan
        const hijriDate = getHijriDate(latitude, longitude, method, timezone, effectiveTime);

        res.status(200).json({ hijriDate, localTime: nowLocal.toISOString(), timezone });
    } catch (error) {
        console.error("❌ API Error:", error);
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah." });
    }
}
