"use client";

import type { HijriDate, Method, HijriAstronomicalData } from "@/types/hijri";
import { formatFullDate } from "@/lib/utils/timezone";
import { formatCoordinates } from "@/lib/utils/maps";

interface HijriPredictionData {
  method: Method;
  today?: HijriDate;
  estimated_next_month_1?: HijriDate | null;
  estimated_gregorian?: string | null;
  message?: string;
  visibility?: HijriAstronomicalData;
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
  if (
    !prediction ||
    !prediction.estimated_next_month_1 ||
    !prediction.estimated_gregorian
  ) {
    return null;
  }

  const isArithmetic = prediction.method === "umm_al_qura";
  const vis = prediction.visibility;
  const userTimezone =
    typeof window !== "undefined"
      ? Intl.DateTimeFormat().resolvedOptions().timeZone
      : "UTC";

  return (
    <div className="group relative animate-in fade-in slide-in-from-bottom-10 duration-1000">
      {/* Aura Glow - Proyeksi Masa Depan */}
      <div className="absolute -inset-1 bg-linear-to-tr from-indigo-500/10 via-primary/5 to-emerald-500/10 rounded-3xl blur-2xl opacity-50 group-hover:opacity-100 transition duration-1000"></div>

      <div className="relative overflow-hidden bg-white/40 dark:bg-card-dark/40 backdrop-blur-3xl p-5 sm:p-8 md:p-10 rounded-2xl border border-white/40 dark:border-white/5 shadow-soft transition-all duration-500">
        {/* Header: Status Proyeksi */}
        <header className="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 sm:mb-10 gap-4">
          <div className="flex items-center gap-3">
            <div className="relative flex h-3 w-3">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-3 w-3 bg-indigo-500"></span>
            </div>
            <h3 className="text-[10px] font-black uppercase tracking-[0.4em] text-gray-400 dark:text-gray-500">
              Lunar Projection Module
            </h3>
          </div>

          {!isArithmetic && vis && (
            <div
              className={`px-4 py-1.5 rounded-full border text-[9px] font-black uppercase tracking-widest transition-all duration-500
              ${
                vis.is_visible
                  ? "bg-primary text-white border-primary shadow-lg shadow-primary/20"
                  : "bg-orange-500/10 text-orange-600 border-orange-500/20"
              }`}
            >
              {vis.is_visible
                ? "Visibility Verified"
                : "Cycle Extended (Istikmal)"}
            </div>
          )}
        </header>

        {/* Prediction Message: Analisis Naratif */}
        {prediction.message && (
          <div className="mb-6 sm:mb-10 p-4 sm:p-6 bg-white/50 dark:bg-white/3 rounded-2xl sm:rounded-3xl border border-gray-100 dark:border-white/5 italic">
            <p className="text-sm leading-relaxed text-gray-600 dark:text-gray-400 font-medium">
              <span className="text-primary font-black not-italic mr-2">
                LOG:
              </span>
              {prediction.message}
            </p>
          </div>
        )}

        {/* Main Result: Target Date (Materi Utama) */}
        <section className="mb-6 sm:mb-10 relative overflow-hidden p-4 sm:p-6 md:p-8 rounded-[1.5rem] sm:rounded-[2rem] bg-linear-to-br from-primary/10 to-emerald-500/5 border border-primary/10 shadow-inner group/result">
          <div className="absolute top-0 right-0 p-6 opacity-10 dark:opacity-20 transition-transform duration-3000 group-hover/result:rotate-180">
            <svg
              width="80"
              height="80"
              viewBox="0 0 100 100"
              fill="none"
              stroke="currentColor"
            >
              <circle
                cx="50"
                cy="50"
                r="45"
                strokeWidth="0.5"
                strokeDasharray="5 5"
              />
              <path
                d="M50 5L50 15M50 85L50 95M5 50L15 50M85 50L95 50"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </div>

          <div className="relative z-10">
            <p className="text-[9px] uppercase opacity-50 font-black mb-3 tracking-[0.3em] text-primary">
              Estimated 1 {MONTHS[prediction.estimated_next_month_1.month - 1]}
            </p>
            <h4 className="text-2xl sm:text-3xl md:text-4xl font-black text-primary tracking-tighter mb-2">
              {formatHijri(prediction.estimated_next_month_1)}
            </h4>
            <div className="flex items-center gap-2">
              <div className="h-px w-4 bg-primary/30"></div>
              <p className="text-sm font-bold text-gray-700 dark:text-gray-300">
                {formatFullDate(prediction.estimated_gregorian, userTimezone)}
              </p>
            </div>
          </div>
        </section>

        {/* Telemetry Grid: Bukti Material */}
        {!isArithmetic && vis && (
          <div className="pt-8 border-t border-gray-100 dark:border-white/5">
            <div className="grid grid-cols-2 gap-8 mb-8">
              <Stat
                label="Moon Altitude"
                value={`${vis.moon_altitude?.toFixed(2)}°`}
                isPrimary
              />
              <Stat
                label="Elongation"
                value={`${vis.elongation?.toFixed(2)}°`}
              />
            </div>

            {vis.lat !== undefined && (
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4 p-5 rounded-2xl bg-gray-50/50 dark:bg-white/5 border border-gray-100 dark:border-transparent">
                <div className="space-y-1">
                  <p className="text-[8px] font-black opacity-30 uppercase tracking-[0.3em]">
                    Projection Reference
                  </p>
                  <p className="text-[11px] font-mono font-bold text-gray-600 dark:text-gray-400">
                    {formatCoordinates(vis.lat, vis.lon || 0)}
                  </p>
                </div>
                <div className="px-3 py-1 bg-primary/10 rounded-lg">
                  <p className="text-[9px] font-black text-primary uppercase tracking-widest">
                    {vis.is_visible ? "✓ Criteria Met" : "✕ Below Limit"}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Footer: Tech Hint */}
        <p className="mt-8 text-[8px] opacity-20 text-center uppercase tracking-[0.5em] font-black">
          Automated Astronomical Simulation Output
        </p>
      </div>
    </div>
  );
}

function Stat({
  label,
  value,
  isPrimary = false,
}: {
  label: string;
  value: string;
  isPrimary?: boolean;
}) {
  return (
    <div className="flex flex-col gap-1">
      <p className="text-[9px] font-black opacity-30 uppercase tracking-[0.2em]">
        {label}
      </p>
      <p
        className={`text-xl sm:text-2xl font-black tabular-nums tracking-tighter ${isPrimary ? "text-primary" : "text-gray-900 dark:text-white"}`}
      >
        {value}
      </p>
    </div>
  );
}
