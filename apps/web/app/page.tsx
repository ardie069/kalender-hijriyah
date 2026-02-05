"use client";

import { useEffect, useState } from "react";
import Clock from "@/components/Clock";
import Method from "@/components/Method";
import HijriDate from "@/components/HijriDate";
import { useHijri } from "@/hooks/useHijri";
import { useTheme } from "@/context/theme-context";
import type { Method as HijriMethod } from "@/types/hijri";

export default function HomePage() {
  const { darkMode } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const bgClass =
    mounted && darkMode ? "bg-gray-950 text-white" : "bg-gray-50 text-black";
  const cardClass =
    mounted && darkMode
      ? "bg-gray-900 border-gray-800 shadow-2xl"
      : "bg-white border-gray-100 shadow-xl";

  const userTimezone = mounted
    ? Intl.DateTimeFormat().resolvedOptions().timeZone
    : "UTC";

  const [selectedMethod, setSelectedMethod] = useState<HijriMethod>("global");

  const { hijriDate, endMonthInfo, weton, loading, error } = useHijri(
    selectedMethod,
    userTimezone,
  );

  return (
    <main
      className={`${bgClass} min-h-[calc(100vh-64px)] transition-colors duration-500 px-4 py-8 md:py-12 flex flex-col items-center`}
    >
      <div
        className={`${cardClass} w-full max-w-md rounded-[2.5rem] border p-6 md:p-8 text-center transition-all duration-500`}
      >
        <header className="mb-8">
          <h1 className="text-2xl md:text-3xl font-extrabold mb-3 tracking-tight">
            🕌 Kalender Hijriyah
          </h1>
          <p className="text-sm md:text-base opacity-70 leading-relaxed">
            Kalender berbasis digital
          </p>
        </header>

        <section className="space-y-6">
          {/* Komponen waktu yang dinamis */}
          <Clock userTimezone={userTimezone} />

          {/* Form kontrol metode */}
          <div className="text-left bg-black/5 dark:bg-white/5 p-4 rounded-2xl">
            <Method value={selectedMethod} onChange={setSelectedMethod} />
          </div>

          {/* Display data utama */}
          <div className="relative pt-4 border-t border-current/10">
            <HijriDate
              hijriDate={hijriDate}
              endMonthInfo={endMonthInfo}
              weton={weton}
              loading={loading}
              error={error}
              method={selectedMethod}
            />
          </div>
        </section>

        <footer className="mt-8 pt-6 border-t border-current/5">
          <p className="text-[10px] md:text-xs opacity-40 uppercase tracking-widest">
            Metode: {selectedMethod} • {userTimezone}
          </p>
        </footer>
      </div>
    </main>
  );
}
