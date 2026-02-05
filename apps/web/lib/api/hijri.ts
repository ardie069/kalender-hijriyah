import type {
  HijriDateResponse,
  HijriEndMonthResponse,
  Method,
} from "@/types/hijri";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

if (!API_BASE) {
  throw new Error("NEXT_PUBLIC_API_BASE_URL is not defined");
}

/**
 * Generic fetch helper
 * Sengaja sederhana, biar gampang debug
 */
async function request<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API error ${res.status}: ${text}`);
  }

  return res.json();
}

/**
 * /hijri-date
 */
export function fetchHijriDate(
  lat: number,
  lon: number,
  method: Method,
  timezone: string,
) {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  return request<HijriDateResponse>(`/hijri-date?${params.toString()}`);
}

/**
 * /hijri-end-month
 */
export function fetchHijriEndMonth(
  lat: number,
  lon: number,
  method: Method,
  timezone: string,
) {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  return request<HijriEndMonthResponse>(
    `/hijri-end-month?${params.toString()}`,
  );
}
