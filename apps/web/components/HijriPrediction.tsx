"use client";

import { useMemo } from "react";
import { useTheme } from "@/context/theme-context";
import type { HijriDate, Method, HijriVisibility } from "@/types/hijri"; // Gunakan type dari file utama

interface HijriPredictionData {
  method: Method; // Tambahkan method di sini
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

  const isGlobal = prediction.method === "global";

  return (
    <div
      className={`p-6 rounded-4xl border transition-all duration-500 ${themeClass}`}
    >
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-bold uppercase tracking-widest opacity-60">
          🌙 Info Akhir Bulan
        </h3>

        {/* Badge Status: Hanya muncul jika BUKAN global */}
        {!isGlobal && prediction.visibility && (
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
            Estimasi Tanggal Baru:
          </p>
          <p className="text-lg font-bold text-emerald-500">
            {formatHijri(prediction.estimated_end_of_month)}
          </p>
          {isGlobal && (
            <p className="text-[10px] mt-2 opacity-40 italic">
              *Berdasarkan ketetapan kalender aritmatika tetap.
            </p>
          )}
        </div>
      )}

      {/* Grid Statistik: Sembunyikan total jika metode global */}
      {!isGlobal && prediction.visibility && (
        <div className="grid grid-cols-2 gap-4 animate-in fade-in zoom-in duration-500">
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
              className={`text-[10px] font-bold ${
                prediction.visibility.is_visible
                  ? "text-emerald-500"
                  : "text-red-500"
              }`}
            >
              {prediction.visibility.is_visible ? "✓ Lolos" : "✕ Tidak Lolos"}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
