export function isLocationInJava(lat, lon) {
  const minLat = -9.5,
    maxLat = -5.5;
  const minLon = 105.5,
    maxLon = 114.0;
  return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon;
}

export function getWeton(gregDate = new Date()) {
  const pasaran = ["Legi", "Pahing", "Pon", "Wage", "Kliwon"];
  const baseDate = new Date(1633, 6, 8);
  const diffDays = Math.floor((gregDate - baseDate) / (1000 * 60 * 60 * 24));
  const pasaranIndex = diffDays % 5;
  return pasaran[pasaranIndex];
}
