"use client";

import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function CalendarPage() {
  const { darkMode } = useTheme();
  const mounted = useMounted();

  const isDark = mounted && darkMode;

  const bgClass = isDark
    ? "bg-gray-950 text-white"
    : "bg-gray-50 text-gray-800";
  const cardClass = isDark
    ? "bg-gray-900 border-gray-800 shadow-2xl shadow-emerald-500/5"
    : "bg-white border-gray-100 shadow-xl shadow-gray-200/40";

  return (
    <main
      className={`${bgClass} min-h-[calc(100vh-64px)] transition-colors duration-500 px-4 py-8 md:py-12`}
    >
      <div className="max-w-6xl mx-auto">
        {/* Header Section */}
        <header className="flex flex-col md:flex-row md:items-end justify-between mb-10 gap-6 px-4">
          <div className="animate-in fade-in slide-in-from-left-4 duration-500">
            <h1 className="text-4xl md:text-6xl font-black tracking-tighter flex items-center gap-4">
              <span className="drop-shadow-sm">📅</span> Kalender
            </h1>
            <p className="opacity-40 mt-2 font-bold uppercase text-[10px] md:text-xs tracking-[0.3em]">
              Syakban — Ramadan 1447 H
            </p>
          </div>

          {/* Navigation Controls */}
          <div className="flex items-center gap-2">
            <button
              className={`${cardClass} p-3 rounded-2xl border flex items-center justify-center hover:bg-emerald-500/5 transition-all active:scale-95`}
            >
              <span className="opacity-60 text-xs">←</span>
            </button>
            <button
              className={`${cardClass} px-6 py-3 rounded-2xl border font-black text-[10px] uppercase tracking-widest bg-emerald-500/3 border-emerald-500/20 text-emerald-600`}
            >
              Hari Ini
            </button>
            <button
              className={`${cardClass} p-3 rounded-2xl border flex items-center justify-center hover:bg-emerald-500/5 transition-all active:scale-95`}
            >
              <span className="opacity-60 text-xs">→</span>
            </button>
          </div>
        </header>

        {/* Container Kalender Tanpa Garis Dalam */}
        <div
          className={`${cardClass} rounded-[3rem] border p-2 md:p-4 animate-in fade-in zoom-in-95 duration-700`}
        >
          {/* Header Hari */}
          <div className="grid grid-cols-7 mb-2">
            {["Ahad", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"].map(
              (day) => (
                <div
                  key={day}
                  className="py-4 text-center text-[9px] md:text-[10px] font-black uppercase tracking-[0.2em] opacity-30"
                >
                  {day}
                </div>
              ),
            )}
          </div>

          {/* Grid Tanggal Seamless */}
          <div className="grid grid-cols-7 gap-1 md:gap-2">
            {Array.from({ length: 35 }).map((_, i) => {
              const hijriDay = (i % 30) + 1;
              const isRamadan = i >= 29;

              return (
                <div
                  key={i}
                  className="group relative min-h-20 md:min-h-32.5 p-3 md:p-5 flex flex-col items-start rounded-4xl hover:bg-emerald-500/4 dark:hover:bg-emerald-500/8 transition-all duration-300 cursor-pointer"
                >
                  {/* Penanda Aktif saat Hover (Opsional: Border muncul tipis) */}
                  <div className="absolute inset-0 rounded-4xl border border-transparent group-hover:border-emerald-500/10 transition-all" />

                  {/* Angka Hijriyah */}
                  <span
                    className={`text-2xl md:text-4xl font-black tabular-nums transition-all duration-300 ${
                      isRamadan
                        ? "text-emerald-500 drop-shadow-[0_0_10px_rgba(16,185,129,0.2)]"
                        : "opacity-80 group-hover:text-emerald-500 group-hover:opacity-100"
                    }`}
                  >
                    {hijriDay}
                  </span>

                  {/* Angka Masehi */}
                  <div className="mt-1 text-[9px] md:text-xs font-bold opacity-20 group-hover:opacity-50 transition-opacity uppercase tracking-tighter">
                    {(i % 28) + 1} Feb
                  </div>

                  {/* Event Marker */}
                  {isRamadan && hijriDay === 1 && (
                    <div className="mt-auto w-full">
                      <div className="bg-emerald-500/10 dark:bg-emerald-500/20 text-emerald-600 dark:text-emerald-400 py-1.5 px-2 rounded-xl text-[7px] md:text-[9px] font-black uppercase tracking-tighter text-center">
                        🌙 1 Ramadan
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        {/* Footer Info */}
        <footer className="mt-10 px-6 flex flex-col md:flex-row justify-between items-center gap-6">
          <p className="text-[10px] md:text-xs opacity-30 italic max-w-sm text-center md:text-left leading-relaxed">
            * Penentuan bulan berbasis visibilitas hilal kriteria Wujudul Hilal.
            Realita lapangan tetap mengikuti ketetapan sidang isbat pemerintah.
          </p>
          <div className="flex gap-6 items-center">
            <div className="flex items-center gap-2.5 text-[9px] font-black uppercase tracking-widest opacity-30">
              <span className="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" />{" "}
              Hari Besar
            </div>
            <div className="flex items-center gap-2.5 text-[9px] font-black uppercase tracking-widest opacity-30">
              <span className="w-2 h-2 rounded-full bg-gray-300 dark:bg-gray-700" />{" "}
              Libur Nasional
            </div>
          </div>
        </footer>
      </div>
    </main>
  );
}
