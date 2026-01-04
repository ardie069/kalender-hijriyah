import type {
  HijriDateResponse,
  HijriEndMonthResponse,
  Method,
} from "@/types/hijri";

const API_BASE = import.meta.env.VITE_API_BASE_URL;

async function request<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) throw new Error("API error");
  return res.json();
}

export function fetchHijriDate(
  lat: number,
  lon: number,
  method: Method,
  timezone: string
) {
  return request<HijriDateResponse>(
    `${API_BASE}/hijri-date?lat=${lat}&lon=${lon}&method=${method}&timezone=${timezone}`
  );
}

export function fetchHijriEndMonth(
  lat: number,
  lon: number,
  method: Method,
  timezone: string
) {
  return request<HijriEndMonthResponse>(
    `${API_BASE}/hijri-end-month?lat=${lat}&lon=${lon}&method=${method}&timezone=${timezone}`
  );
}
