"use client";

import { useEffect, useState } from "react";
import type { Method, HijriDate } from "@/types/hijri";
import { fetchHijriDate, fetchHijriEndMonth } from "@/lib/api/hijri";
import { isLocationInJava, getWeton } from "@/lib/utils/weton";

interface UseHijriResult {
  hijriDate: HijriDate | null;
  endMonthInfo: Record<string, unknown> | null;
  weton: string | null;
  loading: boolean;
  error: string | null;
  reload: () => void;
}

export function useHijri(method: Method, timezone: string): UseHijriResult {
  const [lat, setLat] = useState<number | null>(null);
  const [lon, setLon] = useState<number | null>(null);

  const [hijriDate, setHijriDate] = useState<HijriDate | null>(null);
  const [endMonthInfo, setEndMonthInfo] = useState<Record<string, unknown> | null>(
    null,
  );

  const [weton, setWeton] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Resolve Location
   */
  function resolveLocation(): Promise<{ lat: number; lon: number }> {
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
  }

  /**
   * Load Data
   */
  async function loadHijri() {
    if (!timezone) return;

    setLoading(true);
    setError(null);

    try {
      let currentLat = lat;
      let currentLon = lon;

      if (currentLat === null || currentLon === null) {
        const loc = await resolveLocation();
        currentLat = loc.lat;
        currentLon = loc.lon;
        setLat(loc.lat);
        setLon(loc.lon);
      }

      const [dateRes, endRes] = await Promise.all([
        fetchHijriDate(currentLat, currentLon, method, timezone),
        fetchHijriEndMonth(currentLat, currentLon, method, timezone),
      ]);

      setHijriDate(dateRes.hijri_date);
      setEndMonthInfo(endRes as unknown as Record<string, unknown>);

      if (isLocationInJava(currentLat, currentLon)) {
        const day = new Date().toLocaleDateString("id-ID", {
          weekday: "long",
        });
        setWeton(`${day === "Minggu" ? "Ahad" : day} ${getWeton(new Date())}`);
      } else {
        setWeton(null);
      }
    } catch (err) {
      setError("Gagal mengambil data Hijriyah");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  /**
   * Effect
   */
  useEffect(() => {
    loadHijri();
  }, [method, timezone]);

  return {
    hijriDate,
    endMonthInfo,
    weton,
    loading,
    error,
    reload: loadHijri,
  };
}
