import { getHijriDate } from "../server/hijriCalculator.js";

export default function handler(req, res) {
    console.log("🔍 API Dihit:", req.query); // Tambahkan log untuk melihat query parameters

    const { lat, lon, method } = req.query;

    if (!lat || !lon || !method) {
        console.error("❌ Error: Parameter tidak lengkap");
        return res.status(400).json({ error: "Latitude, Longitude, dan Metode diperlukan!" });
    }

    try {
        const parsedLat = parseFloat(lat);
        const parsedLon = parseFloat(lon);

        if (isNaN(parsedLat) || isNaN(parsedLon)) {
            console.error("❌ Error: Latitude atau Longitude tidak valid");
            return res.status(400).json({ error: "Latitude dan Longitude harus berupa angka" });
        }

        console.log("📡 Memanggil getHijriDate...");
        const hijriDate = getHijriDate(parsedLat, parsedLon, method);
        
        console.log("✅ Hasil Hijriyah:", hijriDate);
        res.status(200).json({ hijriDate });

    } catch (error) {
        console.error("🔥 Error dalam getHijriDate:", error);
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah." });
    }
}
