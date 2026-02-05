"use client";

import { useState } from "react";
import Clock from "@/components/Clock";
import Method from "@/components/Method";
import HijriDate from "@/components/HijriDate";
import { useHijri } from "@/hooks/useHijri";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";
import type { Method as HijriMethod } from "@/types/hijri";

export default function HomePage() {
  const { darkMode } = useTheme();
  const mounted = useMounted();
  const [selectedMethod, setSelectedMethod] = useState<HijriMethod>("global");

  const isDark = mounted && darkMode;

  const bgClass = isDark
    ? "bg-gray-950 text-white"
    : "bg-gray-50 text-gray-800";

  const userTimezone = mounted
    ? Intl.DateTimeFormat().resolvedOptions().timeZone
    : "UTC";

  const { hijriDate, explanation, endMonthInfo, weton, loading, error } =
    useHijri(selectedMethod, userTimezone);

  return (
    <main
      className={`${bgClass} min-h-[calc(100vh-64px)] transition-colors duration-500 px-4 py-6 md:py-12 flex flex-col items-center justify-start`}
    >
      <div className="w-full max-w-6xl">
        <header className="mb-8 md:mb-12 text-center md:text-left px-2">
          <h1 className="text-3xl md:text-5xl font-black mb-3 tracking-tight flex items-center justify-center md:justify-start gap-3">
            <span className="text-4xl md:text-6xl">🕌</span>
            <span>Kalender Hijriyah</span>
          </h1>
          <p className="text-sm md:text-lg opacity-50 max-w-md mx-auto md:mx-0 font-medium leading-relaxed">
            Integrasi sains astronomi dan kearifan lokal dalam satu genggaman
            digital.
          </p>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 md:gap-8 items-start">
          <div className="lg:col-span-7 lg:row-start-1 order-1 animate-in fade-in slide-in-from-left-8 duration-700">
            <Clock userTimezone={userTimezone} />
          </div>

          <aside className="lg:col-span-5 lg:row-span-2 order-2 lg:order-2 space-y-6">
            <div
              className={`p-8 rounded-[3rem] border transition-all duration-500 shadow-xl 
              ${
                isDark
                  ? "bg-gray-900 border-gray-800 shadow-emerald-500/5"
                  : "bg-white border-gray-200/50 shadow-gray-200/30"
              }`}
            >
              <Method value={selectedMethod} onChange={setSelectedMethod} />

              <div className="mt-8 pt-8 border-t border-gray-100/80 dark:border-gray-800">
                <p className="text-[10px] font-bold uppercase tracking-[0.2em] opacity-30 mb-4">
                  Sistem Informasi
                </p>
                <div className="grid grid-cols-2 gap-4">
                  {[
                    { label: "Zona Waktu", value: userTimezone },
                    {
                      label: "Standardisasi",
                      value: selectedMethod.toUpperCase(),
                    },
                  ].map((info, i) => (
                    <div
                      key={i}
                      className="p-4 rounded-2xl bg-gray-50/50 dark:bg-white/5 border border-gray-100 dark:border-transparent"
                    >
                      <p className="text-[9px] opacity-40 uppercase font-black mb-1">
                        {info.label}
                      </p>
                      <p className="text-xs font-mono font-bold truncate opacity-80">
                        {info.value}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="hidden lg:block p-8 rounded-[3rem] bg-emerald-500/5 dark:bg-emerald-500/10 border border-emerald-500/10 dark:border-emerald-500/20">
              <p className="text-sm leading-relaxed text-emerald-800 dark:text-emerald-400 font-medium italic opacity-80">
                <q>
                  Waktu bukan sekadar angka, melainkan gerak materi di ruang
                  angkasa yang kita tangkap melalui logika.
                </q>
              </p>
            </div>
          </aside>

          <div className="lg:col-span-7 lg:row-start-2 order-3 lg:order-3 animate-in fade-in slide-in-from-left-10 duration-1000">
            <HijriDate
              hijriDate={hijriDate}
              explanation={explanation}
              endMonthInfo={endMonthInfo}
              weton={weton}
              loading={loading}
              error={error}
              method={selectedMethod}
            />
          </div>
        </div>

        <footer className="mt-12 py-12 border-t border-gray-100 dark:border-gray-800 text-center md:text-left">
          <div className="flex flex-col md:flex-row justify-between items-center gap-6">
            <p className="text-[10px] md:text-xs opacity-30 uppercase tracking-[0.3em] font-bold">
              © 2026 Kalender Hijriyah Digital • Build with Logic
            </p>
            <div className="flex gap-6 opacity-30 text-[10px] uppercase font-bold tracking-widest">
              <span>Next.js 15</span>
              <span>FastAPI</span>
              <span>Skyfield API</span>
            </div>
          </div>
        </footer>
      </div>
    </main>
  );
}
