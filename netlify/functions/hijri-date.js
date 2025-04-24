const { getHijriDate } = require('../../apps/backend/hijriCalculator');

exports.handler = async (event, context) => {
    const lat = parseFloat(event.queryStringParameters.lat) || 0;
    const lon = parseFloat(event.queryStringParameters.lon) || 0;
    const method = event.queryStringParameters.method || 'global';
    const timezone = event.queryStringParameters.timezone || 'UTC';

    if (isNaN(lat) || isNaN(lon)) {
        return {
            statusCode: 400,
            body: JSON.stringify({ error: "Lokasi tidak valid" }),
        };
    }

    try {
        const hijriDate = getHijriDate(lat, lon, method, timezone);
        return {
            statusCode: 200,
            body: JSON.stringify({ hijriDate }),
        };
    } catch (error) {
        return {
            statusCode: 500,
            body: JSON.stringify({ error: "Terjadi kesalahan pada server" }),
        };
    }
};
