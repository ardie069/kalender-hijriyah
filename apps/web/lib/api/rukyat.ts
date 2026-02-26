import { RukyatResponse } from "@/types/moon";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api";

export async function fetchRukyat(
  lat: number,
  lon: number,
  timezone: string,
  region: string,
  mode: "local" | "national" = "local",
): Promise<RukyatResponse> {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    timezone,
    region,
    mode,
  });

  const res = await fetch(`${API_BASE}/rukyat/evaluate?${params.toString()}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Gagal mengambil data rukyat: ${res.status}`);
  }

  return res.json();
}
