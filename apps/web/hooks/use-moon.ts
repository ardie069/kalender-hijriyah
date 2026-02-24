import { useEffect, useState, useCallback } from "react";
import { MoonDataResponse } from "@/types/moon";
import { fetchMoonInfo } from "@/lib/api/moon";

interface UseMoonResult {
  data: MoonDataResponse | null;
  loading: boolean;
  error: string | null;
  lat: number | null;
  lon: number | null;
  reload: () => void;
}

export function useMoon(): UseMoonResult {
  const [data, setData] = useState<MoonDataResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lat, setLat] = useState<number | null>(null);
  const [lon, setLon] = useState<number | null>(null);

  const resolveLocation = useCallback((): Promise<{
    lat: number;
    lon: number;
  }> => {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error("Geolocation not supported"));
        return;
      }
      navigator.geolocation.getCurrentPosition(
        (pos) =>
          resolve({ lat: pos.coords.latitude, lon: pos.coords.longitude }),
        reject,
      );
    });
  }, []);

  const loadMoon = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const loc = await resolveLocation();
      setLat(loc.lat);
      setLon(loc.lon);
      const res = await fetchMoonInfo(loc.lat, loc.lon);
      setData(res);
    } catch (err) {
      setError("Gagal mensinkronkan data visualisasi bulan.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [resolveLocation]);

  useEffect(() => {
    loadMoon();
  }, [loadMoon]);

  return { data, loading, error, lat, lon, reload: loadMoon };
}
