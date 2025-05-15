export function isLocationInJava(lat, lon) {
  const inJava = lat >= -9.5 && lat <= -5.5 && lon >= 105.5 && lon <= 114.0;
  const isJakartaOrBanten = lon < 107.0;
  return inJava && !isJakartaOrBanten;
}

export function getWeton(gregDate = new Date()) {
  const pasaran = ["Legi", "Pahing", "Pon", "Wage", "Kliwon"];
  const baseDate = new Date(1633, 6, 8);
  const diffDays = Math.floor((gregDate - baseDate) / (1000 * 60 * 60 * 24));
  const pasaranIndex = diffDays % 5;
  return pasaran[pasaranIndex];
}
