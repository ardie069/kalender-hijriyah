/**
 * Mengecek apakah lokasi berada di Pulau Jawa, Madura
 * (kecuali Jakarta & Banten, sesuai tradisi weton Jawa)
 */
export function isLocationInJava(lat: number, lon: number): boolean {
  // Batas geografis kasar Pulau Jawa
  const JAVA_LAT_MIN = -9.5;
  const JAVA_LAT_MAX = -5.5;
  const JAVA_LON_MIN = 105.5;
  const JAVA_LON_MAX = 114.0;

  // Jakarta & Banten dikecualikan (tradisi weton tidak umum dipakai)
  const JAKARTA_BANTEN_LON_MAX = 107.0;

  const inJava =
    lat >= JAVA_LAT_MIN &&
    lat <= JAVA_LAT_MAX &&
    lon >= JAVA_LON_MIN &&
    lon <= JAVA_LON_MAX;

  const isJakartaOrBanten = lon < JAKARTA_BANTEN_LON_MAX;

  return inJava && !isJakartaOrBanten;
}

/**
 * Menghitung weton Jawa berdasarkan kalender Sultan Agung (1633 M)
 * Basis: 8 Juni 1633 = Jumat Legi
 */
export function getWeton(date: Date = new Date()): string {
  const PASARAN = ["Legi", "Pahing", "Pon", "Wage", "Kliwon"] as const;

  // Basis kalender Jawa: 8 Juni 1633 M (Gregorian)
  const BASE_DATE = new Date(Date.UTC(1633, 6, 8));

  const target = new Date(
    Date.UTC(date.getFullYear(), date.getMonth(), date.getDate())
  );

  const diffDays =
    Math.floor(
      (target.getTime() - BASE_DATE.getTime()) / (1000 * 60 * 60 * 24)
    ) % 5;

  return PASARAN[(diffDays + 5) % 5]!;
}
