import { MoonDataResponse } from "@/types/moon";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api";

export async function fetchMoonInfo(
  lat: number,
  lon: number,
): Promise<MoonDataResponse> {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
  });

  const res = await fetch(`${API_BASE}/moon/info?${params.toString()}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Gagal mengambil data bulan: ${res.status}`);
  }

  return res.json();
}
