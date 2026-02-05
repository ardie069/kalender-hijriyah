"use client";

import { useMemo } from "react";
import { useTheme } from "@/context/theme-context";
import type { HijriDate } from "@/types/hijri";

interface HijriVisibility {
  moon_altitude: number;
  elongation: number;
  moon_age: number;
  is_visible: boolean;
}

interface HijriPredictionData {
  today?: HijriDate;
  estimated_end_of_month?: HijriDate | null;
  message?: string;
  visibility?: HijriVisibility;
}

interface HijriPredictionProps {
  prediction: HijriPredictionData | null;
}

const MONTHS = [
  "Muharam",
  "Safar",
  "Rabiulawal",
  "Rabiulakhir",
  "Jumadilawal",
  "Jumadilakhir",
  "Rajab",
  "Syakban",
  "Ramadan",
  "Syawal",
  "Zulkaidah",
  "Zulhijah",
];

function formatHijri(date: HijriDate): string {
  return `${date.day} ${MONTHS[date.month - 1]} ${date.year} H`;
}

export default function HijriPrediction({ prediction }: HijriPredictionProps) {
  const { darkMode } = useTheme();

  const themeClass = useMemo(
    () =>
      darkMode
        ? "bg-white/[0.03] border-white/10"
        : "bg-black/[0.03] border-black/5",
    [darkMode],
  );

  if (!prediction) return null;

  return (
    <div
      className={`p-6 rounded-4xl border transition-all duration-500 ${themeClass}`}
    >
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-bold uppercase tracking-widest opacity-60">
          🌙 Info Akhir Bulan
        </h3>
        {prediction.visibility && (
          <span
            className={`text-[10px] px-3 py-1 rounded-full font-black uppercase tracking-tighter ${
              prediction.visibility.is_visible
                ? "bg-emerald-500/20 text-emerald-500"
                : "bg-red-500/20 text-red-500"
            }`}
          >
            {prediction.visibility.is_visible ? "Hilal Terlihat" : "Istikmal"}
          </span>
        )}
      </div>

      {prediction.message && (
        <p className="text-sm leading-relaxed opacity-80 mb-4 italic">
          <q>{prediction.message}</q>
        </p>
      )}

      {prediction.estimated_end_of_month && (
        <div className="mb-6 p-4 rounded-2xl bg-emerald-500/5 border border-emerald-500/10">
          <p className="text-[10px] uppercase opacity-50 font-bold mb-1">
            Estimasi Tanggal:
          </p>
          <p className="text-lg font-bold text-emerald-500">
            {formatHijri(prediction.estimated_end_of_month)}
          </p>
        </div>
      )}

      {prediction.visibility && (
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1">
            <p className="text-[9px] uppercase opacity-40 font-bold">
              Altitude
            </p>
            <p className="font-mono text-sm font-bold">
              {prediction.visibility.moon_altitude.toFixed(2)}°
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-[9px] uppercase opacity-40 font-bold">
              Elongasi
            </p>
            <p className="font-mono text-sm font-bold">
              {prediction.visibility.elongation.toFixed(2)}°
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-[9px] uppercase opacity-40 font-bold">
              Usia Bulan
            </p>
            <p className="font-mono text-sm font-bold">
              {prediction.visibility.moon_age.toFixed(1)}j
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-[9px] uppercase opacity-40 font-bold">
              Kriteria
            </p>
            <p
              className={`text-[10px] font-bold ${prediction.visibility.is_visible ? "text-emerald-500" : "text-red-500"}`}
            >
              {prediction.visibility.is_visible ? "✓ Lolos" : "✕ Tidak Lolos"}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
