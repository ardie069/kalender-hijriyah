import express from 'express';
import dotenv from 'dotenv';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { getHijriDate, predictEndOfMonth } from './hijriCalculator.js';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;
const DEFAULT_LAT = parseFloat(process.env.DEFAULT_LAT) || 0;
const DEFAULT_LON = parseFloat(process.env.DEFAULT_LON) || 0;
const HIJRI_METHOD = process.env.HIJRI_METHOD || 'global';
const TIMEZONE = process.env.TIMEZONE || 'UTC';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Middleware
app.use(express.json());
app.use(cors({
    origin: 'http://127.0.0.1:5173', // Sesuaikan dengan origin frontend Anda
    methods: ['GET', 'POST'],
    allowedHeaders: ['Content-Type'],
}));

app.use((req, res, next) => {
    res.setHeader('Cache-Control', 'no-store, no-cache, must-revalidate, proxy-revalidate');
    res.setHeader('Pragma', 'no-cache');
    res.setHeader('Expires', '0');
    next();
});

// Endpoint: Tanggal Hijriyah hari ini
app.get("/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
    const lon = parseFloat(req.query.lon) || DEFAULT_LON;
    const method = req.query.method || HIJRI_METHOD;
    const timezone = req.query.timezone || TIMEZONE;

    if (isNaN(lat) || isNaN(lon)) {
        return res.status(400).json({ error: "Lokasi tidak valid" });
    }

    try {
        const hijriDate = getHijriDate(lat, lon, method, timezone);
        res.json({ hijriDate });
    } catch (error) {
        res.status(500).json({ error: "Terjadi kesalahan pada server" });
    }
});

// Endpoint: Prediksi akhir bulan Hijriyah
app.get("/hijri-end-month", (req, res) => {
    const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
    const lon = parseFloat(req.query.lon) || DEFAULT_LON;
    const method = req.query.method || HIJRI_METHOD;
    const timezone = req.query.timezone || TIMEZONE;

    try {
        const prediction = predictEndOfMonth(lat, lon, method, timezone);
        res.json(prediction);
    } catch (error) {
        console.error("❌ Gagal prediksi akhir bulan:", error);
        res.status(500).json({ error: "Terjadi kesalahan saat memprediksi akhir bulan" });
    }
});

// Sajikan file statis frontend (jika dibuild)
app.use(express.static(path.join(__dirname, '../frontend/dist')));

// Catch-all: fallback ke index.html untuk SPA
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../frontend/dist', 'index.html'));
});

// Menjalankan server
app.listen(PORT, () => {
    console.log(`🚀 Server berjalan di http://localhost:${PORT}/`);
});
