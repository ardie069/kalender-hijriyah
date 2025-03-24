import { getHijriDate } from "../server/hijriCalculator";

export default function handler(req, res) {
    const { lat, lon, method, timezone } = req.query; // 🟢 Ambil timezone

    if (!lat || !lon || !method) {
        return res.status(400).json({ error: "Latitude, Longitude, dan Method diperlukan!" });
    }

    try {
        console.log("📡 API Request:", { lat, lon, method, timezone }); // 🔍 Debug log
        const hijriDate = getHijriDate(parseFloat(lat), parseFloat(lon), method);

        res.status(200).json({ hijriDate, timezone }); // 🟢 Tambahkan timezone ke response
    } catch (error) {
        console.error("❌ API Error:", error); // 🔍 Debug error
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah." });
    }
}
