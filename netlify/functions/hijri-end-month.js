import { predictEndOfMonth } from "../../apps/backend/hijriCalculator";

exports.handler = async (event, context) => {
    const lat = parseFloat(event.queryStringParameters.lat) || 0;
    const lon = parseFloat(event.queryStringParameters.lon) || 0;
    const method = event.queryStringParameters.method || 'global';
    const timezone = event.queryStringParameters.timezone || 'UTC';

    try {
        const prediction = predictEndOfMonth(lat, lon, method, timezone);
        return {
            statusCode: 200,
            body: JSON.stringify(prediction),
        };
    } catch (error) {
        console.error("❌ Gagal prediksi akhir bulan:", error);
        return {
            statusCode: 500,
            body: JSON.stringify({ error: "Terjadi kesalahan saat memprediksi akhir bulan" }),
        };
    }
};
