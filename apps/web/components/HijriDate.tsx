"use client";

import HijriInfoCard from "./HijriInfoCard";
import HijriPrediction from "./HijriPrediction";
import HijriSkeleton from "./HijriSkeleton";
import { formatGeneratedTime } from "@/lib/utils/timezone";
import { getHilalMapUrl, formatCoordinates } from "@/lib/utils/maps";
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
  if ((loading && !hijriDate) || (!hijriDate && !error)) return <HijriSkeleton />;

  if (error) {
    return (
      <div className="bg-error/10 border border-error/20 rounded-3xl p-6 flex items-start gap-4 animate-in fade-in zoom-in duration-500">
        <div className="bg-error text-error-content rounded-xl p-2 shrink-0 shadow-lg shadow-error/20">
          <span className="text-xl font-bold">⚠️</span>
        </div>
        <div>
          <h3 className="font-black text-error">System Critical Failure</h3>
          <p className="text-error/80 text-sm mt-1 leading-relaxed font-medium">
            {error}
          </p>
        </div>
      </div>
    );
  }

  if (!hijriDate) return null;

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-8 duration-1000 px-1 sm:px-0">
      {/* 1. Main Display: Fokus pada Identitas Waktu */}
      <HijriInfoCard hijriDate={hijriDate} weton={weton} method={method} />

      {/* 2. Reasoning Section: The Logic Stream */}
      {explanation?.reasoning?.length ? (
        <div className="relative group animate-in fade-in slide-in-from-bottom-10 duration-1000 delay-150">
          {/* Efek Pendaran Halus di Belakang (Glow) */}
          <div className="absolute -inset-2 bg-linear-to-br from-primary/5 via-transparent to-emerald-500/5 rounded-5xl blur-2xl opacity-0 group-hover:opacity-100 transition duration-1000"></div>

          <div className="relative overflow-hidden bg-white/40 dark:bg-card-dark/40 backdrop-blur-2xl rounded-4xl border border-white/40 dark:border-white/5 shadow-[0_8px_32px_0_rgba(0,0,0,0.05)] transition-all hover:shadow-primary/5">
            {/* Top Header Section */}
            <div className="p-8 pb-4">
              <header className="flex flex-col sm:flex-row sm:items-end justify-between gap-6 mb-10">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="h-1 w-8 bg-primary rounded-full"></div>
                    <h3 className="text-[10px] font-black text-primary uppercase tracking-[0.4em]">
                      Logic Analysis Report
                    </h3>
                  </div>
                  <p className="text-2xl font-black tracking-tight text-gray-900 dark:text-white">
                    Kenapa {hijriDate.day} {MONTHS[hijriDate.month - 1]}?
                  </p>
                </div>

                <div className="flex flex-wrap gap-2">
                  {explanation.after_sunset && (
                    <div className="flex items-center gap-2 px-4 py-2 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl">
                      <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                      <p className="text-[9px] text-emerald-600 dark:text-emerald-400 font-black uppercase tracking-widest">
                        After Sunset
                      </p>
                    </div>
                  )}
                  <div
                    className={`px-4 py-2 rounded-2xl border font-black text-[9px] uppercase tracking-widest
              ${
                explanation.decision === "new_month"
                  ? "bg-primary text-white border-primary shadow-lg shadow-primary/20"
                  : "bg-orange-500/10 text-orange-600 border-orange-500/20"
              }`}
                  >
                    {explanation.decision === "new_month"
                      ? "New Crescent"
                      : "Istikmal"}
                  </div>
                </div>
              </header>

              {/* Reasoning Flow: Timeline Style */}
              <div className="relative space-y-6 mb-10 ml-2">
                <div className="absolute left-1.75 top-2 bottom-2 w-px bg-gray-200 dark:bg-gray-800"></div>
                {explanation.reasoning.map((note, i) => (
                  <div key={i} className="flex gap-6 items-start group/item">
                    <div className="relative z-10 w-3.5 h-3.5 rounded-full bg-white dark:bg-gray-950 border-2 border-primary group-hover/item:scale-125 transition-transform duration-300 shadow-sm"></div>
                    <p className="text-sm text-gray-600 dark:text-gray-400 font-medium italic leading-relaxed pt-0.5">
                      {note}
                    </p>
                  </div>
                ))}
              </div>
            </div>

            {/* Bottom Data Section: Telemetry Widgets */}
            {explanation.astronomical_data && (
              <div className="bg-gray-50/50 dark:bg-white/2 p-8 pt-10 border-t border-gray-100 dark:border-white/5">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-10 mb-10">
                  <Stat
                    label="Lunar Altitude"
                    value={`${explanation.astronomical_data.moon_altitude?.toFixed(2)}°`}
                    isPrimary
                  />
                  <Stat
                    label="Elongation Angle"
                    value={`${explanation.astronomical_data.elongation?.toFixed(2)}°`}
                  />
                </div>

                {/* Reference Map Action Card */}
                {explanation.astronomical_data.lat !== undefined && (
                  <div className="relative overflow-hidden p-6 rounded-3xl bg-white dark:bg-gray-900 border border-gray-100 dark:border-white/5 group/map shadow-soft">
                    {/* Background Hint (Simulasi Map) */}
                    <div className="absolute inset-0 opacity-[0.03] dark:opacity-[0.07] pointer-events-none grayscale group-hover/map:scale-110 transition-transform duration-2000">
                      <svg viewBox="0 0 100 100" className="w-full h-full">
                        <path d="M0 0h100v100H0z" />
                        <path
                          d="M10 10l80 80M10 90L90 10"
                          stroke="currentColor"
                          strokeWidth="0.5"
                        />
                      </svg>
                    </div>

                    <div className="relative flex flex-col sm:flex-row items-center justify-between gap-6">
                      <div className="text-center sm:text-left space-y-1">
                        <p className="text-[9px] font-black text-gray-400 uppercase tracking-[0.3em]">
                          Telemetry Source
                        </p>
                        <p className="text-sm font-mono font-bold text-gray-900 dark:text-gray-100">
                          {explanation.astronomical_data.location_name ||
                            formatCoordinates(
                              explanation.astronomical_data.lat,
                              explanation.astronomical_data.lon!,
                            )}
                        </p>
                      </div>
                      <a
                        href={getHilalMapUrl(
                          explanation.astronomical_data.lat,
                          explanation.astronomical_data.lon!,
                        )}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="w-full sm:w-auto px-8 py-4 rounded-2xl bg-gray-900 dark:bg-primary text-white dark:text-black font-black text-xs tracking-widest hover:scale-[1.05] active:scale-95 transition-all shadow-xl shadow-gray-900/10 dark:shadow-primary/20 flex items-center justify-center gap-3"
                      >
                        📍 GLOBAL TRACKER
                      </a>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      ) : null}

      {/* 3. Predictive Analysis Section: The Future Scan */}
      {endMonthInfo && (
        <section className="relative group animate-in slide-in-from-bottom-6 duration-700 delay-300">
          {/* Aksentuasi Cahaya di Belakang Card */}
          <div className="absolute -inset-1 bg-linear-to-r from-primary/20 to-emerald-500/10 rounded-4xl blur-xl opacity-50 group-hover:opacity-100 transition duration-1000"></div>

          <div className="relative">
            <HijriPrediction
              prediction={{
                ...endMonthInfo,
                method,
                estimated_next_month_1: endMonthInfo.estimated_next_month_1,
                estimated_gregorian: endMonthInfo.estimated_gregorian,
                message: endMonthInfo.message ?? undefined,
                visibility: endMonthInfo.visibility,
              }}
            />
          </div>
        </section>
      )}

      {/* 4. Live Synchronization Status: Floating System Health */}
      {generatedAt && (
        <div className="flex justify-center items-center py-4">
          <div className="group flex items-center gap-4 px-6 py-3 bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl rounded-full border border-white/20 dark:border-white/5 shadow-soft hover:shadow-primary/10 transition-all duration-500">
            <div className="relative flex h-3 w-3">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
              <span className="relative inline-flex rounded-full h-3 w-3 bg-primary shadow-[0_0_10px_rgba(16,185,129,0.5)]"></span>
            </div>

            <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3">
              <p className="text-[10px] font-black dark:text-primary-content text-base-content/60 uppercase tracking-[0.25em]">
                System Telemetry
              </p>
              <div className="hidden sm:block h-3 w-px bg-base-content/10"></div>
              <p className="text-[10px] font-bold text-primary uppercase tracking-widest">
                Sync OK: {formatGeneratedTime(generatedAt)}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Stat Component: Refactored for more 'Airy' feel
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
      <div className="flex items-center gap-2">
        <div
          className={`w-1 h-3 rounded-full ${isPrimary ? "bg-primary" : "bg-base-content/20"}`}
        ></div>
        <p className="text-[9px] font-black opacity-40 uppercase tracking-[0.2em]">
          {label}
        </p>
      </div>
      <p
        className={`text-2xl sm:text-4xl font-black tabular-nums tracking-tighter transition-all ${isPrimary ? "text-primary drop-shadow-[0_0_15px_rgba(16,185,129,0.2)]" : "text-base-content"}`}
      >
        {value}
      </p>
    </div>
  );
}
