"use client";

import { useMemo } from "react";
import { useTheme } from "@/context/theme-context";
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
  const { darkMode } = useTheme();

  // Dialektika Warna: Pake background yang lebih deep buat Dark Mode
  const containerClass = useMemo(
    () =>
      darkMode
        ? "bg-gray-900 border-gray-800 text-white shadow-2xl shadow-emerald-500/10"
        : "bg-white border-gray-100 text-black shadow-xl shadow-black/5",
    [darkMode],
  );

  const methodLabel = useMemo(() => {
    const labels: Record<Method, string> = {
      global: "Global Standard",
      hisab: "Hisab Wujudul Hilal",
      rukyat: "Rukyat Hilal Indonesia",
    };
    return labels[method] || method;
  }, [method]);

  return (
    <div
      className={`relative overflow-hidden mt-6 p-6 rounded-4xl border transition-all duration-500 ${containerClass}`}
    >
      {/* Dekorasi halus buat kesan modern */}
      <div className="absolute -top-10 -right-10 w-32 h-32 bg-emerald-500/10 rounded-full blur-3xl" />

      <div className="relative z-10">
        <header className="flex justify-between items-start mb-4">
          <div>
            <p className="text-[10px] font-bold uppercase tracking-[0.2em] opacity-50">
              Kalender Hijriyah
            </p>
            <h2 className="text-3xl font-black tracking-tight mt-1">
              {hijriDate.day} {MONTHS[hijriDate.month - 1]}
            </h2>
            <p className="text-xl font-medium opacity-80">
              {hijriDate.year} Hijriyah
            </p>
          </div>
          <div className="bg-emerald-500 text-white p-2 rounded-2xl shadow-lg shadow-emerald-500/20">
            🕌
          </div>
        </header>

        {/* Weton & Metadata */}
        <div className="flex flex-wrap items-center gap-2 mt-6">
          {weton && (
            <span className="px-3 py-1.5 rounded-full bg-emerald-500/10 text-emerald-500 text-xs font-bold border border-emerald-500/20">
              ✨ {weton}
            </span>
          )}
          <span className="px-3 py-1.5 rounded-full bg-black/5 dark:bg-white/5 text-[10px] font-medium opacity-60">
            {methodLabel}
          </span>
        </div>

        {/* Footer info kecil */}
        <p className="mt-4 text-[9px] opacity-40 leading-relaxed italic">
          * Perubahan hari Hijriyah & Weton terjadi setiap masuk waktu Maghrib.
        </p>
      </div>
    </div>
  );
}
