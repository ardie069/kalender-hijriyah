"use client";

import HijriInfoCard from "./HijriInfoCard";
import HijriPrediction from "./HijriPrediction";
import type {
  HijriDate as HijriDateType,
  Method,
  HijriEndMonthResponse as HijriPredictionType,
  HijriExplanation,
} from "@/types/hijri";

interface HijriDateProps {
  hijriDate: HijriDateType | null;
  endMonthInfo: HijriPredictionType | null;
  weton: string | null;
  loading: boolean;
  error: string | null;
  method: Method;
  explanation?: HijriExplanation | null;
}

export default function HijriDate({
  hijriDate,
  endMonthInfo,
  weton,
  loading,
  error,
  method,
  explanation,
}: HijriDateProps) {
  if (loading) {
    return (
      <div className="mt-6 space-y-6 animate-pulse px-2">
        <div className="h-48 bg-gray-200 dark:bg-gray-800 rounded-[2.5rem]" />
        <div className="h-20 bg-gray-100 dark:bg-gray-900 rounded-3xl w-full" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="mt-6 p-5 rounded-4xl bg-red-500/10 border border-red-500/20 text-red-500 text-sm flex items-start gap-4">
        <span className="text-xl">⚠️</span>
        <div className="space-y-1">
          <p className="font-bold">Koneksi Terputus</p>
          <p className="opacity-80 leading-relaxed">{error}</p>
        </div>
      </div>
    );
  }

  if (!hijriDate) return null;

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-6 duration-700 ease-out px-1">
      {/* Kartu Informasi Utama */}
      <HijriInfoCard hijriDate={hijriDate} weton={weton} method={method} />

      {/* Reasoning Section (The Scientific Soul) */}
      {explanation && explanation.reasoning.length > 0 && (
        <div className="p-5 rounded-4xl bg-black/3 dark:bg-white/3 border border-current/5 text-left transition-all hover:bg-black/5 dark:hover:bg-white/5">
          <div className="flex items-center justify-between mb-3">
            <p className="text-[10px] font-bold opacity-40 uppercase tracking-widest">
              Analisis Logika & Sains
            </p>
            {explanation.astronomical_data && (
              <span
                className={`text-[10px] px-2 py-0.5 rounded-full font-bold uppercase ${explanation.astronomical_data.is_visible ? "bg-emerald-500/20 text-emerald-600" : "bg-orange-500/20 text-orange-600"}`}
              >
                {explanation.astronomical_data.is_visible
                  ? "Hilal Terlihat"
                  : "Istikmal"}
              </span>
            )}
          </div>

          <div className="space-y-2">
            {explanation.reasoning.map((note, i) => (
              <p
                key={i}
                className="text-xs opacity-70 italic leading-relaxed flex gap-2"
              >
                <span className="text-emerald-500 shrink-0">●</span> {note}
              </p>
            ))}
          </div>

          {/* Menampilkan angka astronomis jika ada */}
          {explanation.astronomical_data && (
            <div className="mt-4 pt-3 border-t border-current/5 grid grid-cols-2 gap-4">
              <div>
                <p className="text-[9px] opacity-40 uppercase font-bold">
                  Tinggi Hilal
                </p>
                <p className="text-sm font-mono font-bold text-emerald-500">
                  {explanation.astronomical_data.moon_altitude.toFixed(2)}°
                </p>
              </div>
              <div>
                <p className="text-[9px] opacity-40 uppercase font-bold">
                  Status Visibilitas
                </p>
                <p className="text-sm font-bold opacity-80">
                  {explanation.astronomical_data.is_visible
                    ? "Memenuhi Syarat"
                    : "Dibawah Kriteria"}
                </p>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Prediksi Akhir Bulan */}
      {endMonthInfo && (
        <div className="pt-2">
          <HijriPrediction
            prediction={{
              ...endMonthInfo,
              message: endMonthInfo.message ?? undefined,
            }}
          />
        </div>
      )}
    </div>
  );
}
