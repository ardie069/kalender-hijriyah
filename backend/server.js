import express from 'express';
import dotenv from 'dotenv';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { getHijriDate, predictEndOfMonth } from './hijriCalculator.js';

dotenv.config();

const app = express();

const PORT = 3000;
const DEFAULT_LAT = parseFloat(process.env.DEFAULT_LAT) || 0;
const DEFAULT_LON = parseFloat(process.env.DEFAULT_LON) || 0;
const HIJRI_METHOD = process.env.HIJRI_METHOD || 'global';
const TIMEZONE = process.env.TIMEZONE || 'UTC';

app.use(cors({
    origin: ['http://localhost:5173', 'https://kalender-hijriyah.vercel.app'], // Sesuaikan dengan origin frontend kamu
    methods: ['GET', 'POST'],
    allowedHeaders: ['Content-Type'],
}));

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Menambahkan cache control untuk menghindari masalah dengan cache browser
app.use((req, res, next) => {
    res.setHeader('Cache-Control', 'no-store, no-cache, must-revalidate, proxy-revalidate');
    res.setHeader('Pragma', 'no-cache');
    res.setHeader('Expires', '0');
    next();
});

app.use(express.static(path.join(__dirname, '../frontend/dist')));

app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../frontend/dist', 'index.html'));
});

app.get("/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || "global";
    const timezone = req.query.timezone || "UTC";

    if (isNaN(lat) || isNaN(lon)) {
        return res.status(400).json({ error: "Lokasi tidak valid" });
    }

    const hijriDate = getHijriDate(lat, lon, method, timezone);
    res.json({ hijriDate });
});

app.get("/hijri-end-month", (req, res) => {
    const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
    const lon = parseFloat(req.query.lon) || DEFAULT_LON;
    const method = req.query.method || HIJRI_METHOD;
    const timezone = req.query.timezone || TIMEZONE;

    const prediction = predictEndOfMonth(lat, lon, method, timezone);
    res.json(prediction);
});

app.listen(PORT, () => {
    console.log(`Server berjalan di http://localhost:${PORT}/`);
});
