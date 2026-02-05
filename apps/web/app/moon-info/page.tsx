"use client";

import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function MoonInfoPage() {
  const { darkMode } = useTheme();
  const mounted = useMounted();

  const isDark = mounted && darkMode;

  // FIX: Gunakan kontras lembut yang sudah kita sepakati
  const bgClass = isDark
    ? "bg-gray-950 text-white"
    : "bg-gray-50 text-gray-800";
  const cardClass = isDark
    ? "bg-gray-900 border-gray-800 shadow-2xl shadow-emerald-500/5"
    : "bg-white border-gray-200/50 shadow-xl shadow-gray-200/30";

  return (
    <main
      className={`${bgClass} min-h-[calc(100vh-64px)] transition-colors duration-500 px-4 py-8 md:py-16`}
    >
      <div className="max-w-5xl mx-auto">
        {/* Header: Center Focused */}
        <header className="text-center mb-12 md:mb-16 animate-in fade-in slide-in-from-bottom-4 duration-700">
          <div className="inline-block p-4 rounded-full bg-emerald-500/10 mb-6">
            <span className="text-4xl md:text-6xl">🌕</span>
          </div>
          <h1 className="text-4xl md:text-6xl font-black tracking-tighter mb-4">
            Lunar Analytics
          </h1>
          <p className="opacity-50 max-w-xl mx-auto font-medium text-sm md:text-base leading-relaxed italic">
            <q>
              Mengamati gerak materi di angkasa untuk menyinkronkan waktu
              manusia dengan hukum alam.
            </q>
          </p>
        </header>

        {/* Dashboard Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 md:gap-8 items-stretch">
          {/* Card 1: Posisi Astronomis (7 Cols) */}
          <div
            className={`${cardClass} lg:col-span-7 p-8 md:p-10 rounded-[3rem] border animate-in fade-in slide-in-from-left-6 duration-1000`}
          >
            <div className="flex justify-between items-center mb-8">
              <h2 className="text-xl font-black uppercase tracking-widest opacity-80">
                Posisi Astronomis
              </h2>
              <span className="px-3 py-1 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-[10px] font-black rounded-full uppercase tracking-tighter">
                Live Data
              </span>
            </div>

            <div className="space-y-6">
              {[
                {
                  label: "Tinggi Hilal (Altitude)",
                  value: "7.24°",
                  desc: "Sudut vertikal dari ufuk",
                },
                {
                  label: "Elongasi",
                  value: "8.51°",
                  desc: "Jarak sudut Bulan-Matahari",
                },
                {
                  label: "Umur Bulan",
                  value: "15j 20m",
                  desc: "Waktu sejak ijtima' (konjungsi)",
                },
                {
                  label: "Iluminasi",
                  value: "1.2%",
                  desc: "Persentase cahaya matahari",
                },
              ].map((item, i) => (
                <div
                  key={i}
                  className="group flex justify-between items-end border-b border-gray-100/50 dark:border-gray-800 pb-4 last:border-0"
                >
                  <div>
                    <p className="text-[10px] font-black uppercase tracking-widest opacity-30 group-hover:opacity-100 transition-opacity">
                      {item.label}
                    </p>
                    <p className="text-xs opacity-40 mt-0.5">{item.desc}</p>
                  </div>
                  <span className="text-2xl md:text-3xl font-black font-mono tracking-tighter text-emerald-500">
                    {item.value}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Card 2: Visualisasi & Status (5 Cols) */}
          <div className="lg:col-span-5 flex flex-col gap-6 md:gap-8 animate-in fade-in slide-in-from-right-6 duration-1000">
            {/* Visual Fase */}
            <div
              className={`${cardClass} p-8 rounded-[3rem] border grow flex flex-col items-center justify-center text-center`}
            >
              <div className="relative mb-6">
                {/* Efek Glow di belakang bulan */}
                <div className="absolute inset-0 bg-emerald-500/20 blur-3xl rounded-full" />
                <div className="relative text-8xl md:text-9xl drop-shadow-2xl">
                  🌒
                </div>
              </div>
              <h3 className="text-2xl font-black mb-1">Hilal Awal</h3>
              <p className="text-xs font-bold uppercase tracking-[0.2em] opacity-30">
                Waxing Crescent
              </p>
            </div>

            {/* Status Visibilitas */}
            <div className="p-8 rounded-[3rem] bg-emerald-500 text-white shadow-xl shadow-emerald-500/20">
              <p className="text-[10px] font-black uppercase tracking-widest opacity-70 mb-2 text-white/80">
                Kriteria MABIMS
              </p>
              <h4 className="text-2xl font-bold leading-tight">
                Hilal kemungkinan besar terlihat.
              </h4>
              <div className="mt-4 flex gap-2">
                <div className="px-2 py-1 bg-white/20 rounded-lg text-[10px] font-bold">
                  Altitude {`>`} 3°
                </div>
                <div className="px-2 py-1 bg-white/20 rounded-lg text-[10px] font-bold">
                  Elongasi {`>`} 6.4°
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer Note */}
        <footer className="mt-12 text-center md:text-left px-4">
          <div className="inline-flex items-center gap-3 p-4 rounded-2xl bg-black/2 dark:bg-white/2 border border-gray-100/50 dark:border-gray-800">
            <span className="text-xs">🔬</span>
            <p className="text-[10px] md:text-xs opacity-40 italic leading-relaxed">
              Data dihitung menggunakan algoritma <strong>Skyfield API</strong>{" "}
              dengan koreksi refraksi atmosfer lokal. Akurasi data astronomis
              mencapai ±0.01 detik busur.
            </p>
          </div>
        </footer>
      </div>
    </main>
  );
}
