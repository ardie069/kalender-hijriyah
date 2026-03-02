import { useState, useCallback, useEffect } from "react";
import type { DateSystem, UnifiedMonthData } from "@/types/calendar";
import type { Method } from "@/types/hijri";
import {
  fetchYearlyCalendar,
  fetchGregorianCalendar,
} from "@/lib/api/calendar";
import { fetchHijriDate } from "@/lib/api/hijri";

export interface TodayInfo {
  month: number;
  day: number;
  year: number;
}

export function useCalendar(
  year: number,
  dateSystem: DateSystem,
  method: Method,
  lat: number | null,
  lon: number | null,
) {
  const [months, setMonths] = useState<UnifiedMonthData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [today, setToday] = useState<TodayInfo | null>(null);

  const loadCalendar = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      if (dateSystem === "hijri") {
        if (lat === null || lon === null) return;

        // Fetch calendar + current hijri date in parallel
        const [calRes, hijriRes] = await Promise.all([
          fetchYearlyCalendar(year, lat, lon, method),
          fetchHijriDate(lat, lon, method, "Asia/Jakarta"),
        ]);

        setMonths(
          calRes.months.map((m) => ({
            month_id: m.month_id,
            month_name: m.month_name,
            total_days: m.total_days,
            day_1_weekday: m.day_1_weekday,
            start_gregorian: m.start_gregorian,
          })),
        );

        // Set today dari hijri-date endpoint
        const hd = hijriRes.hijri_date;
        setToday({ year: hd.year, month: hd.month, day: hd.day });
      } else {
        const res = await fetchGregorianCalendar(year);
        setMonths(
          res.months.map((m) => ({
            month_id: m.month_id,
            month_name: m.month_name,
            total_days: m.total_days,
            day_1_weekday: m.day_1_weekday,
          })),
        );

        // Today untuk Masehi — dari client
        const now = new Date();
        setToday({
          year: now.getFullYear(),
          month: now.getMonth() + 1,
          day: now.getDate(),
        });
      }
    } catch {
      setError("Gagal memuat data kalender.");
    } finally {
      setLoading(false);
    }
  }, [year, dateSystem, method, lat, lon]);

  useEffect(() => {
    loadCalendar();
  }, [loadCalendar]);

  return { months, loading, error, today, refresh: loadCalendar };
}
