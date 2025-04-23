import express from 'express';
import path from 'path';
import { fileURLToPath } from 'url';
import { getHijriDate, predictEndOfMonth } from '../backend/hijriCalculator.js';

const app = express();

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

app.use(express.json());

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
    console.log("📅 Hijri date:", hijriDate);
    res.json({ hijriDate });
  } catch (error) {
    console.error("❌ Gagal mengambil tanggal Hijriyah:", error);
    res.status(500).json({ error: "Terjadi kesalahan pada server" });
  }
});

app.get("/api/hijri-end-month", (req, res) => {
  const lat = parseFloat(req.query.lat) || 0;
  const lon = parseFloat(req.query.lon) || 0;
  const method = req.query.method || "global";
  const timezone = req.query.timezone || "UTC";

  try {
    const prediction = predictEndOfMonth(lat, lon, method, timezone);
    res.json(prediction);
  } catch (error) {
    console.error("❌ Gagal prediksi akhir bulan:", error);
    res.status(500).json({ error: "Terjadi kesalahan saat memprediksi akhir bulan" });
  }
});

app.use(express.static(path.join(__dirname, '../frontend/dist')));
app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, '../frontend/dist/index.html'));
});

export default app;
