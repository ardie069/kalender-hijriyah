"use client";

import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function CalendarPage() {
  const { darkMode } = useTheme();
  const mounted = useMounted();

  const bgClass =
    mounted && darkMode ? "bg-gray-950 text-white" : "bg-gray-50 text-black";
  const cardClass =
    mounted && darkMode
      ? "bg-gray-900 border-gray-800"
      : "bg-white border-gray-100";

  return (
    <main
      className={`${bgClass} min-h-[calc(100vh-64px)] transition-colors duration-500 px-4 py-8`}
    >
      <div className="max-w-5xl mx-auto">
        <header className="flex flex-col md:flex-row md:items-end justify-between mb-8 gap-4">
          <div>
            <h1 className="text-3xl font-extrabold flex items-center gap-3">
              <span>📅</span> Kalender Hijriyah
            </h1>
            <p className="opacity-70 mt-1">Syakban - Ramadan 1447 H</p>
          </div>

          <div className="flex gap-2">
            <button
              className={`${cardClass} px-4 py-2 rounded-xl border text-sm hover:scale-105 transition-transform`}
            >
              &larr; Prev
            </button>
            <button
              className={`${cardClass} px-4 py-2 rounded-xl border text-sm hover:scale-105 transition-transform`}
            >
              Next &rarr;
            </button>
          </div>
        </header>

        {/* Container Kalender */}
        <div
          className={`${cardClass} rounded-4xl border shadow-2xl overflow-hidden`}
        >
          {/* Header Hari (Desktop) */}
          <div className="hidden md:grid grid-cols-7 bg-current/5 border-b border-current/10">
            {["Ahad", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"].map(
              (day) => (
                <div
                  key={day}
                  className="py-4 text-center text-xs font-bold uppercase tracking-widest opacity-60"
                >
                  {day}
                </div>
              ),
            )}
          </div>

          {/* Grid Tanggal */}
          <div className="grid grid-cols-7 md:divide-x md:divide-y divide-current/5">
            {/* Contoh render kotak tanggal */}
            {Array.from({ length: 30 }).map((_, i) => (
              <div
                key={i}
                className="min-h-20 md:min-h-30 p-2 md:p-4 hover:bg-emerald-500/5 transition-colors cursor-pointer group"
              >
                <span className="text-lg md:text-2xl font-bold group-hover:text-emerald-500">
                  {i + 1}
                </span>
                <div className="mt-1 text-[10px] md:text-xs opacity-50">
                  {i + 1} Feb
                </div>
                {/* Penanda khusus misal awal bulan atau puasa */}
                {i === 29 && (
                  <div className="mt-2 text-[8px] md:text-[10px] bg-emerald-500/20 text-emerald-600 p-1 rounded">
                    1 Ramadhan*
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        <footer className="mt-6 text-center">
          <p className="text-xs opacity-50 italic">
            * Estimasi berdasarkan metode Hisab Astronomis. Realita di lapangan
            bisa berubah sesuai hasil Rukyat.
          </p>
        </footer>
      </div>
    </main>
  );
}
