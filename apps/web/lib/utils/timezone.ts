/**
 * Memformat waktu ISO menjadi format Indonesia yang manusiawi
 * Contoh: 21 Feb 2026, 17:21 WIB
 */
export function formatGeneratedTime(dateStr: string | undefined): string {
  if (!dateStr) return "--:--";

  const date = new Date(dateStr);

  return new Intl.DateTimeFormat("id-ID", {
    day: "numeric",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZoneName: "short",
  }).format(date);
}

/**
 * Helper untuk menampilkan hari dan tanggal
 */
export function formatFullDate(dateString: string, timeZone: string): string {
  const date = new Date(dateString);

  return new Intl.DateTimeFormat("id-ID", {
    weekday: "long",
    day: "2-digit",
    month: "long",
    year: "numeric",
    timeZone,
  }).format(date);
}
