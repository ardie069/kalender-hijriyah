import { useEffect, useState, useCallback } from "react";
import { RukyatResponse } from "@/types/moon";
import { fetchRukyat } from "@/lib/api/rukyat";

interface UseRukyatResult {
  data: RukyatResponse | null;
  loading: boolean;
  error: string | null;
  lat: number | null;
  lon: number | null;
  reload: () => void;
}

export function useRukyat(
  region: string = "indonesia",
  mode: "local" | "national" = "local",
): UseRukyatResult {
  const [data, setData] = useState<RukyatResponse | null>(null);
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
          resolve({
            lat: pos.coords.latitude,
            lon: pos.coords.longitude,
          }),
        reject,
      );
    });
  }, []);

  const loadRukyat = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const loc = await resolveLocation();
      setLat(loc.lat);
      setLon(loc.lon);

      const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      const res = await fetchRukyat(loc.lat, loc.lon, timezone, region, mode);

      setData(res);
    } catch (err) {
      console.error(err);
      setError("Gagal memuat data rukyat.");
    } finally {
      setLoading(false);
    }
  }, [resolveLocation, region, mode]);

  useEffect(() => {
    loadRukyat();
  }, [loadRukyat]);

  return {
    data,
    loading,
    error,
    lat,
    lon,
    reload: loadRukyat,
  };
}
