import express from 'express';
import serverless from 'serverless-http';
import { getHijriDate, predictEndOfMonth } from '../backend/hijriCalculator.js';

const app = express();

app.use(express.json());

app.get("/api/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || 'global';
    const timezone = req.query.timezone || 'UTC';

    try {
        const hijriDate = getHijriDate(lat, lon, method, timezone);
        res.json({ hijriDate });
    } catch (err) {
        res.status(500).json({ error: "Server error" });
    }
});

app.get("/api/hijri-end-month", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || 'global';
    const timezone = req.query.timezone || 'UTC';

    try {
        const prediction = predictEndOfMonth(lat, lon, method, timezone);
        res.json(prediction);
    } catch (err) {
        res.status(500).json({ error: "Server error" });
    }
});

export default serverless(app);
