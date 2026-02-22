"use client";

import HijriInfoCard from "./HijriInfoCard";
import HijriPrediction from "./HijriPrediction";
import { formatGeneratedTime } from "@/lib/utils/timezone";
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
  generatedAt?: string;
}

export default function HijriDate({
  hijriDate,
  endMonthInfo,
  weton,
  loading,
  error,
  method,
  explanation,
  generatedAt,
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
      {explanation?.reasoning?.length ? (
        <div className="p-5 rounded-4xl bg-black/3 dark:bg-white/3 border border-current/5 text-left transition-all hover:bg-black/5 dark:hover:bg-white/5">
          <div className="flex items-center justify-between mb-3">
            <p className="text-[10px] font-bold opacity-40 uppercase tracking-widest">
              Analisis Logika & Sains
            </p>
            {explanation.decision === "new_month" && (
              <span className="text-[10px] px-2 py-0.5 rounded-full font-bold uppercase bg-emerald-500/20 text-emerald-600">
                Bulan Baru
              </span>
            )}

            {explanation.decision === "istikmal_30" && (
              <span className="text-[10px] px-2 py-0.5 rounded-full font-bold uppercase bg-orange-500/20 text-orange-600">
                Istikmal
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

          {explanation.astronomical_data?.global_visible !== undefined && (
            <div className="mt-2 text-[10px] opacity-50">
              Status Global:{" "}
              {explanation.astronomical_data.global_visible
                ? "Memenuhi Global"
                : "Belum Global"}
            </div>
          )}

          {/* Menampilkan angka astronomis jika ada */}
          {explanation.astronomical_data && (
            <div className="mt-4 pt-3 border-t border-current/5 grid grid-cols-2 gap-4">
              <div>
                <p className="text-[9px] opacity-40 uppercase font-bold">
                  Tinggi Hilal
                </p>
                <p className="text-sm font-mono font-bold text-emerald-500">
                  {typeof explanation.astronomical_data.moon_altitude ===
                  "number"
                    ? `${explanation.astronomical_data.moon_altitude.toFixed(2)}°`
                    : "—"}
                </p>
              </div>

              <div>
                <p className="text-[9px] opacity-40 uppercase font-bold">
                  Elongasi
                </p>
                <p className="text-sm font-mono font-bold opacity-80">
                  {typeof explanation.astronomical_data.elongation === "number"
                    ? `${explanation.astronomical_data.elongation.toFixed(2)}°`
                    : "—"}
                </p>
              </div>
            </div>
          )}
        </div>
      ) : null}

      {/* Prediksi Akhir Bulan */}
      {endMonthInfo && (
        <div className="pt-2">
          <HijriPrediction
            prediction={{
              ...endMonthInfo,
              method: method,
              message: endMonthInfo.message ?? undefined,
            }}
          />
        </div>
      )}

      {/* Footer: Timestamp (The Final Touch) */}
      {generatedAt && (
        <div className="pt-4 pb-2 text-center">
          <p className="text-[10px] font-medium opacity-30 flex items-center justify-center gap-1.5 uppercase tracking-[0.2em]">
            <span className="relative flex h-1.5 w-1.5">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500"></span>
            </span>
            Sistem diperbarui: {formatGeneratedTime(generatedAt)}
          </p>
        </div>
      )}
    </div>
  );
}
