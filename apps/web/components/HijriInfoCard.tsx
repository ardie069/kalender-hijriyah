"use client";

import { useMemo } from "react";
import type { HijriDate, Method } from "@/types/hijri";

interface HijriInfoCardProps {
  hijriDate: HijriDate;
  weton: string | null;
  method: Method;
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

export default function HijriInfoCard({
  hijriDate,
  weton,
  method,
}: HijriInfoCardProps) {
  const methodLabel = useMemo(() => {
    const labels: Record<Method, string> = {
      umm_al_qura: "Umm al-Qura",
      local_hisab: "Hisab Lokal",
      local_rukyat: "Rukyat Lokal",
      ughc: "KHGT (Global)",
    };
    return labels[method];
  }, [method]);

  return (
    <div className="group relative animate-in fade-in zoom-in duration-700">
      {/* Glow Aura di belakang kartu - Dialektika Cahaya */}
      <div className="absolute -inset-1 bg-linear-to-tr from-primary/30 to-emerald-400/20 rounded-3xl blur-2xl opacity-20 group-hover:opacity-40 transition-opacity duration-1000"></div>

      <div className="relative overflow-hidden bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl p-6 sm:p-10 md:p-12 rounded-[2rem] sm:rounded-[2.8rem] border border-white/40 dark:border-white/5 shadow-card transition-all duration-500 hover:border-primary/20">
        {/* Dekorasi Ornamen: Tekstur Geometris Halus */}
        <div className="absolute top-0 right-0 p-8 opacity-5 dark:opacity-10 pointer-events-none">
          <svg
            width="120"
            height="120"
            viewBox="0 0 100 100"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <circle
              cx="50"
              cy="50"
              r="48"
              stroke="currentColor"
              strokeWidth="0.5"
              strokeDasharray="4 4"
            />
            <circle
              cx="50"
              cy="50"
              r="30"
              stroke="currentColor"
              strokeWidth="1"
            />
          </svg>
        </div>

        <div className="relative z-10 flex flex-col items-center">
          <header className="mb-6 sm:mb-10 w-full text-center">
            <div className="inline-flex items-center gap-3 px-4 py-1.5 rounded-full bg-primary/5 border border-primary/10 mb-6">
              <span className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse"></span>
              <p className="text-[10px] font-black uppercase tracking-[0.4em] text-primary">
                Tanggal Hijriyah Saat Ini
              </p>
            </div>

            <div className="flex flex-col items-center gap-1">
              {/* Angka Tanggal: Materi Utama yang Bold */}
              <h2 className="text-6xl sm:text-8xl md:text-9xl font-black tracking-tighter text-gray-900 dark:text-white tabular-nums drop-shadow-sm">
                {hijriDate.day}
              </h2>
              <div className="space-y-0">
                <p className="text-2xl sm:text-3xl md:text-4xl font-black text-gray-800 dark:text-gray-200 tracking-tight">
                  {MONTHS[hijriDate.month - 1]}
                </p>
                <p className="text-lg font-bold opacity-30 italic tracking-[0.2em] mt-1">
                  {hijriDate.year} H
                </p>
              </div>
            </div>
          </header>

          {/* Footer Card: Meta-Data Analysis */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 w-full pt-10 border-t border-gray-100 dark:border-white/5">
            <div className="flex flex-col items-center sm:items-start gap-1">
              <span className="text-[9px] font-black opacity-30 uppercase tracking-widest">
                Tradisi & Kearifan
              </span>
              <div className="flex items-center gap-2">
                <span className="text-lg">✨</span>
                <span className="font-bold text-gray-800 dark:text-gray-200">
                  {weton || "Data Nihil"}
                </span>
              </div>
            </div>

            <div className="flex flex-col items-center sm:items-end gap-1">
              <span className="text-[9px] font-black opacity-30 uppercase tracking-widest">
                Sistem Standardisasi
              </span>
              <div className="flex items-center gap-2">
                <span className="px-3 py-1 rounded-lg bg-gray-900 dark:bg-primary text-white dark:text-black text-[10px] font-black uppercase tracking-tighter">
                  {methodLabel}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
