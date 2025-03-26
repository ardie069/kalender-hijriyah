import { getHijriDate, predictEndOfMonth } from "../server/hijriCalculator.js";

export default function handler(req, res) {
    const { query } = req;
    const lat = parseFloat(query.lat) || 0;
    const lon = parseFloat(query.lon) || 0;
    const method = query.method || "global";
    const timezone = query.timezone || "UTC";

    if (req.url.includes("/hijri-date")) {
        const hijriDate = getHijriDate(lat, lon, method, timezone);
        return res.json({ hijriDate });
    }

    if (req.url.includes("/hijri-end-month")) {
        const endOfMonthPrediction = predictEndOfMonth(lat, lon, method, timezone);
        return res.json(endOfMonthPrediction);
    }

    res.status(404).json({ error: "Not Found" });
}
