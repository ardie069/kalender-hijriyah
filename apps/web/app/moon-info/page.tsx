"use client";

import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function MoonInfoPage() {
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
      <div className="max-w-4xl mx-auto">
        <header className="text-center mb-10">
          <h1 className="text-3xl md:text-4xl font-extrabold mb-4 flex items-center justify-center gap-3">
            <span>🌕</span> Informasi Bulan
          </h1>
          <p className="opacity-70 max-w-2xl mx-auto">
            Data spesifik mengenai posisi bulan, elongasi, dan visibilitas hilal
            untuk menentukan awal bulan Hijriyah secara ilmiah.
          </p>
        </header>

        {/* Grid Konten - Responsif: 1 kolom di HP, 2 di Desktop */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className={`${cardClass} p-6 rounded-3xl border shadow-sm`}>
            <h2 className="text-xl font-bold mb-4">Posisi Astronomis</h2>
            <div className="space-y-4 opacity-80">
              <p className="flex justify-between border-b border-current/5 pb-2">
                <span>Tinggi Hilal:</span>
                <span className="font-mono">7.2°</span>
              </p>
              <p className="flex justify-between border-b border-current/5 pb-2">
                <span>Elongasi:</span>
                <span className="font-mono">8.5°</span>
              </p>
              <p className="flex justify-between">
                <span>Umur Bulan:</span>
                <span className="font-mono">15j 20m</span>
              </p>
            </div>
          </div>

          <div
            className={`${cardClass} p-6 rounded-3xl border shadow-sm flex items-center justify-center`}
          >
            {/* Tempat buat visualisasi fase bulan nantinya */}
            <div className="text-center">
              <div className="text-6xl mb-2">🌒</div>
              <p className="text-sm font-medium">Fase: Hilal Awal</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
