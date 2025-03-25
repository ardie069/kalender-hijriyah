import SunCalc from "suncalc";
import { getHijriDate } from "../server/hijriCalculator.js";
import { DateTime } from "luxon";

export default function handler(req, res) {
    try {
        const { lat, lon, method, timezone } = req.query;

        // Validasi parameter
        if (!lat || !lon || !method || !timezone) {
            return res.status(400).json({ error: "Latitude, Longitude, Method, dan Timezone diperlukan!" });
        }

        const latitude = parseFloat(lat);
        const longitude = parseFloat(lon);

        if (isNaN(latitude) || isNaN(longitude)) {
            return res.status(400).json({ error: "Latitude dan Longitude harus berupa angka yang valid!" });
        }

        console.log("📡 API Request:", { latitude, longitude, method, timezone });

        // Hitung waktu matahari terbenam di lokasi pengguna dalam UTC
        const nowUTC = DateTime.utc();
        const sunsetUTC = SunCalc.getTimes(nowUTC.toJSDate(), latitude, longitude).sunset;
        const sunsetLocal = DateTime.fromJSDate(sunsetUTC).setZone(timezone);
        const nowLocal = nowUTC.setZone(timezone);

        console.log(`⏳ Sekarang UTC: ${nowUTC.toISO()}`);
        console.log(`⏳ Sekarang Lokal (${timezone}): ${nowLocal.toISO()}`);
        console.log(`🌅 Matahari terbenam (${timezone}): ${sunsetLocal.toISO()}`);

        // Ambil tanggal Hijriyah berdasarkan parameter
        const hijriDate = getHijriDate(latitude, longitude, method, timezone);

        res.status(200).json({
            hijriDate,
            localTime: nowLocal.toISO(),
            sunsetTime: sunsetLocal.toISO(),
            timezone
        });
    } catch (error) {
        console.error("❌ API Error:", error);
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah. Silakan coba lagi." });
    }
}
