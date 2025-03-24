export default function handler(req, res) {
    const { lat, lon, method } = req.query;

    if (!lat || !lon || !method) {
        return res.status(400).json({ error: "Latitude, Longitude, dan Method diperlukan!" });
    }

    try {
        const hijriDate = getHijriDate(parseFloat(lat), parseFloat(lon), method);
        res.status(200).json({ hijriDate });
    } catch (error) {
        console.error("Error:", error);
        res.status(500).json({ error: "Gagal menghitung tanggal Hijriyah." });
    }
}
