/**
 * Menghasilkan URL Google Maps untuk koordinat tertentu.
 * Digunakan untuk memverifikasi lokasi penampakan hilal (Global Scan).
 */
export const getHilalMapUrl = (lat: number, lon: number): string => {
  // Kita pake format 'z=5' (zoom level) biar kelihatan konteks wilayahnya (negara/benua)
  // 'q=lat,lon' buat naruh pin di titik presisi
  return `https://www.google.com/maps/search/?api=1&query=${lat},${lon}`;
};

/**
 * Memformat koordinat jadi teks yang enak dibaca manusia.
 */
export const formatCoordinates = (lat: number, lon: number): string => {
  const latDir = lat >= 0 ? "LU" : "LS";
  const lonDir = lon >= 0 ? "BT" : "BB";
  return `${Math.abs(lat).toFixed(2)}° ${latDir}, ${Math.abs(lon).toFixed(2)}° ${lonDir}`;
};
