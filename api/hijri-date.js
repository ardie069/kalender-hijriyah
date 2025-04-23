import { getHijriDate, predictEndOfMonth } from "../backend/hijriCalculator";
import Cors from "cors";

// Inisialisasi CORS
const cors = Cors({
  methods: ['GET', 'POST'],
  allowedHeaders: ['Content-Type'],
  origin: 'http://localhost:5173', // Atur origin yang sesuai
});

// Menggunakan middleware CORS di API Route
function runMiddleware(req, res, fn) {
  return new Promise((resolve, reject) => {
    fn(req, res, (result) => {
      if (result instanceof Error) {
        return reject(result);
      }
      return resolve(result);
    });
  });
}

export default async function handler(req, res) {
  // Jalankan middleware CORS
  await runMiddleware(req, res, cors);

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
