import express from 'express';
import dotenv from 'dotenv';
import { getHijriDate } from './hijriCalculator.js';

dotenv.config();

const app = express();
const PORT = 3000;
const DEFAULT_LAT = parseFloat(process.env.DEFAULT_LAT) || 0;
const DEFAULT_LON = parseFloat(process.env.DEFAULT_LON) || 0;
const HIJRI_METHOD = process.env.HIJRI_METHOD || 'global';

app.use(express.static('public'));

app.get("/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || "global"; // Default ke metode "global"

    if (isNaN(lat) || isNaN(lon)) {
        return res.status(400).json({ error: "Lokasi tidak valid" });
    }

    const hijriDate = getHijriDate(lat, lon, method); // Gunakan metode yang dipilih
    res.json({ hijriDate });
});

app.listen(PORT, () => {
    console.log(`Server berjalan di http://localhost:${PORT}/`);
});
