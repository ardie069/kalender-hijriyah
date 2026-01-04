import type { HijriDateResponse, Method } from "@/types/hijri";

const API_BASE = import.meta.env.VITE_API_BASE_URL;

if (!API_BASE) {
  throw new Error("VITE_API_BASE_URL is not defined");
}

export async function fetchHijriDate(
  lat: number,
  lon: number,
  method: Method,
  timezone: string
): Promise<HijriDateResponse> {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  const res = await fetch(`${API_BASE}/hijri-date?${params.toString()}`);

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Hijri API error (${res.status}): ${text}`);
  }

  return res.json();
}
