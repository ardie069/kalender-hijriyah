import express from "express";
import { getHijriDate, predictEndOfMonth } from "../server/hijriCalculator.js";

const app = express();

app.get("/api/hijri-date", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || "global";
    const timezone = req.query.timezone || "UTC";

    const hijriDate = getHijriDate(lat, lon, method, timezone);
    res.json({ hijriDate });
});

app.get("/api/hijri-end-month", (req, res) => {
    const lat = parseFloat(req.query.lat) || 0;
    const lon = parseFloat(req.query.lon) || 0;
    const method = req.query.method || "global";
    const timezone = req.query.timezone || "UTC";

    const endOfMonthPrediction = predictEndOfMonth(lat, lon, method, timezone);
    res.json(endOfMonthPrediction);
});

// Ekspor handler supaya bisa dipakai di Vercel
export default app;
