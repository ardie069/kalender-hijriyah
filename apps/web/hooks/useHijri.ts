"use client";

import { useEffect, useState, useCallback } from "react";
import type {
  Method,
  HijriDate,
  HijriEndMonthResponse,
  HijriExplanation,
} from "@/types/hijri";
import { fetchHijriDate, fetchHijriEndMonth } from "@/lib/api/hijri";
import { isLocationInJava, getWeton } from "@/lib/utils/weton";

interface UseHijriResult {
  hijriDate: HijriDate | null;
  explanation: HijriExplanation | null;
  endMonthInfo: HijriEndMonthResponse | null;
  weton: string | null;
  loading: boolean;
  error: string | null;
  reload: () => void;
}

export function useHijri(method: Method, timezone: string): UseHijriResult {
  const [lat, setLat] = useState<number | null>(null);
  const [lon, setLon] = useState<number | null>(null);

  const [hijriDate, setHijriDate] = useState<HijriDate | null>(null);
  const [explanation, setExplanation] = useState<HijriExplanation | null>(null);
  const [endMonthInfo, setEndMonthInfo] =
    useState<HijriEndMonthResponse | null>(null);

  const [weton, setWeton] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  const loadHijri = useCallback(async () => {
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
      setExplanation(dateRes.explanation ?? null);
      setEndMonthInfo(endRes ?? null);

      if (isLocationInJava(currentLat, currentLon)) {
        const now = new Date();
        const dayName = now.toLocaleDateString("id-ID", { weekday: "long" });
        const fixDay = dayName === "Minggu" ? "Ahad" : dayName;

        setWeton(`${fixDay} ${getWeton(now)}`);
      } else {
        setWeton(null);
      }
    } catch (err) {
      setError("Gagal mensinkronkan data dengan posisi benda langit.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [lat, lon, method, timezone, resolveLocation]);

  useEffect(() => {
    loadHijri();
  }, [loadHijri]);

  return {
    hijriDate,
    explanation,
    endMonthInfo,
    weton,
    loading,
    error,
    reload: loadHijri,
  };
}
