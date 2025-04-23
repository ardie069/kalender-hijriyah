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

// Konfigurasi CORS untuk mengizinkan permintaan dari frontend lokal
app.use(cors({
    origin: 'http://127.0.0.1:5173', // Pastikan URL frontend yang benar
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

// Menyajikan file statis frontend (apabila dibuild)
app.use(express.static(path.join(__dirname, '../frontend/dist')));

// Menyajikan index.html ketika diakses di luar API
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../frontend/dist', 'index.html'));
});

// Endpoint untuk mendapatkan tanggal Hijriyah
app.get("/api/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || "global";
    const timezone = req.query.timezone || "UTC";

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

// Endpoint untuk memprediksi akhir bulan Hijriyah
app.get("/api/hijri-end-month", (req, res) => {
    const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
    const lon = parseFloat(req.query.lon) || DEFAULT_LON;
    const method = req.query.method || HIJRI_METHOD;
    const timezone = req.query.timezone || TIMEZONE;

    const prediction = predictEndOfMonth(lat, lon, method, timezone);
    res.json(prediction);
});

// Menjalankan server di port yang telah ditentukan
app.listen(PORT, () => {
    console.log(`Server berjalan di http://localhost:${PORT}/`);
});
