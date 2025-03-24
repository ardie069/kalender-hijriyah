import { getHijriDate } from "../server/hijriCalculator.js";

export default function handler(req, res) {
    console.log("📡 API Request:", req.query); // ✅ Debug request

    const { lat, lon, method, timezone } = req.query;

    if (!lat || !lon || !method) {
        console.error("❌ Missing Parameters:", { lat, lon, method, timezone });
        return res.status(400).json({ error: "Latitude, Longitude, dan Method diperlukan!" });
    }

    try {
        console.log("✅ Memproses Hijri Date...");
        const hijriDate = getHijriDate(parseFloat(lat), parseFloat(lon), method);

        console.log("✅ Sukses! Hijri Date:", hijriDate);
        res.status(200).json({ hijriDate, timezone });
    } catch (error) {
        console.error("❌ API Error:", error);
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah." });
    }
}
