"use client";

import HijriInfoCard from "./HijriInfoCard";
import HijriPrediction from "./HijriPrediction";

import type { HijriDate as HijriDateType, Method } from "@/types/hijri";

interface HijriVisibility {
  moon_altitude: number;
  elongation: number;
  moon_age: number;
  is_visible: boolean;
}

interface HijriPredictionData {
  today?: HijriDateType;
  estimated_end_of_month?: HijriDateType | null;
  message?: string;
  visibility?: HijriVisibility;
}

interface HijriDateProps {
  hijriDate: HijriDateType | null;
  endMonthInfo: HijriPredictionData | null;
  weton: string | null;
  loading: boolean;
  error: string | null;
  method: Method;
}

export default function HijriDate({
  hijriDate,
  endMonthInfo,
  weton,
  loading,
  error,
  method,
}: HijriDateProps) {
  // ==========================
  // LOADING
  // ==========================
  if (loading) {
    return (
      <div className="flex items-center justify-center mt-4">
        <span className="animate-spin mr-2">⏳</span>
        <span className="font-semibold">📍 Menghitung tanggal Hijriyah…</span>
      </div>
    );
  }

  // ==========================
  // ERROR
  // ==========================
  if (error) {
    return (
      <div className="mt-4 p-3 rounded bg-red-100 text-red-700 text-sm">
        {error}
      </div>
    );
  }

  // ==========================
  // EMPTY (defensive)
  // ==========================
  if (!hijriDate) {
    return null;
  }

  // ==========================
  // RESULT
  // ==========================
  return (
    <div>
      <HijriInfoCard hijriDate={hijriDate} weton={weton} method={method} />

      {endMonthInfo && <HijriPrediction prediction={endMonthInfo} />}
    </div>
  );
}
