import express from 'express';
import { getHijriDate, predictEndOfMonth } from '../hijriCalculator.js';
import serverless from 'serverless-http';

const app = express();

const DEFAULT_LAT = 0;
const DEFAULT_LON = 0;
const HIJRI_METHOD = 'global';
const TIMEZONE = 'UTC';

app.get('/api/hijri-date', (req, res) => {
  const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
  const lon = parseFloat(req.query.lon) || DEFAULT_LON;
  const method = req.query.method || HIJRI_METHOD;
  const timezone = req.query.timezone || TIMEZONE;

  try {
    const hijriDate = getHijriDate(lat, lon, method, timezone);
    res.json({ hijriDate });
  } catch (error) {
    res.status(500).json({ error: "Terjadi kesalahan pada server" });
  }
});

app.get('/api/hijri-end-month', (req, res) => {
  const lat = parseFloat(req.query.lat) || DEFAULT_LAT;
  const lon = parseFloat(req.query.lon) || DEFAULT_LON;
  const method = req.query.method || HIJRI_METHOD;
  const timezone = req.query.timezone || TIMEZONE;

  try {
    const prediction = predictEndOfMonth(lat, lon, method, timezone);
    res.json(prediction);
  } catch (error) {
    res.status(500).json({ error: "Terjadi kesalahan saat memprediksi akhir bulan" });
  }
});

export default serverless(app);
