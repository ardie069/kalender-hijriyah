"use client";

import Metric from "./Metric";
import StatCard from "./StatCard";

interface MoonVisualizationProps {
  illumination: number; // 0 - 100
  phaseName: string;
  age: string;
  distance: string;
  elongation: string;
  lat?: number; // Ditambahin buat deteksi lokasi
}

export default function MoonVisualizationCard({
  illumination = 14.2,
  phaseName = "Waxing Crescent",
  age = "3.4",
  distance = "384,400",
  elongation = "42.5",
  lat = -6.2, // Default Jakarta (Southern Hemisphere)
}: MoonVisualizationProps) {
  // 1. Deteksi Fase (Membesar vs Mengecil)
  const isWaning =
    phaseName.toLowerCase().includes("waning") ||
    phaseName.toLowerCase().includes("third");
  const isWaxing = !isWaning;

  // 2. Deteksi Hemisfer (Utara vs Selatan)
  const isNorth = lat >= 0;

  // 3. LOGIKA ARAH CAHAYA (Dialektika Lokasi)
  // Di Utara: Waxing di Kanan, Waning di Kiri.
  // Di Selatan (Indonesia): Waxing di Kiri, Waning di Kanan.
  const lightOnRight = (isNorth && isWaxing) || (!isNorth && isWaning);

  // 4. Hitung lebar elips shaper
  const percentage = illumination / 100;

  return (
    <div className="space-y-8">
      <div className="bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl p-8 sm:p-12 rounded-5xl border border-white/40 dark:border-white/5 shadow-soft overflow-hidden relative group transition-all duration-500">
        <div className="relative z-10">
          <div className="flex justify-between items-center mb-12">
            <div>
              <h2 className="text-2xl font-black tracking-tight text-gray-900 dark:text-white">
                Visualisasi Real-Time
              </h2>
              <p className="text-[10px] font-black uppercase tracking-[0.2em] opacity-40">
                Observer: {lat >= 0 ? "Northern" : "Southern"} Hemisphere
              </p>
            </div>
            <span className="badge badge-primary font-black text-[10px] py-3 px-4 rounded-full uppercase tracking-widest">
              {phaseName}
            </span>
          </div>

          <div className="flex flex-col md:flex-row items-center gap-16 justify-center">
            {/* --- MOON PHASE ENGINE (LOCATION AWARE) --- */}
            <div className="relative w-56 h-56 sm:w-64 sm:h-64 shrink-0">
              <div className="absolute inset-0 bg-primary/20 blur-[80px] rounded-full animate-pulse"></div>

              <div className="relative w-full h-full rounded-full bg-gray-950 shadow-2xl overflow-hidden border border-white/10 flex items-center justify-center">
                {/* Layer 1: Sisi Gelap Total */}
                <div className="absolute inset-0 bg-gray-900 shadow-inner"></div>

                {/* Layer 2: Setengah Lingkaran Terang (Posisi tergantung Lokasi + Fase) */}
                <div
                  className={`absolute w-1/2 h-full bg-gray-100 shadow-[0_0_30px_rgba(255,255,255,0.3)] transition-all duration-700
                    ${lightOnRight ? "right-0" : "left-0"}`}
                ></div>

                {/* Layer 3: Elips Shaper (Pembentuk Sabit/Cembung) */}
                <div
                  className={`absolute h-[101%] rounded-full transition-all duration-700 ease-in-out
                    ${illumination < 50 ? "bg-gray-900" : "bg-gray-100"}`}
                  style={{
                    width: "100%",
                    transform: `scaleX(${Math.abs(1 - percentage * 2)})`,
                    zIndex: 10,
                  }}
                ></div>

                {/* Texture & Depth */}
                <div className="absolute inset-0 opacity-10 pointer-events-none bg-[url('https://www.transparenttextures.com/patterns/pollen.png')] mix-blend-overlay"></div>
                <div className="absolute inset-0 bg-linear-to-tr from-black/40 via-transparent to-white/10 rounded-full"></div>
              </div>

              {/* Live Indicator */}
              <div className="absolute bottom-4 right-4 z-20 flex items-center gap-2 bg-white/90 dark:bg-black/50 backdrop-blur px-3 py-1.5 rounded-full shadow-sm border border-gray-100 dark:border-gray-800">
                <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
                <span className="text-[10px] font-black uppercase tracking-widest">
                  Live Analytics
                </span>
              </div>
            </div>

            {/* Quick Stats Grid */}
            <div className="grid grid-cols-2 gap-y-10 gap-x-8 w-full max-w-sm">
              <Metric
                label="Iluminasi"
                value={`${illumination}%`}
                sub="Light Density"
              />
              <Metric
                label="Umur Bulan"
                value={`${age} Hari`}
                sub="Cycle Age"
              />
              <Metric
                label="Jarak"
                value={`${distance} km`}
                sub="Earth-Moon Distance"
              />
              <Metric
                label="Elongasi"
                value={`${elongation}°`}
                sub="Sun-Moon Angle"
              />
            </div>
          </div>

          {/* Phase Progress Bar */}
          <div className="mt-16 pt-8 border-t border-gray-100 dark:border-white/5">
            <div className="flex justify-between text-[9px] font-black uppercase tracking-widest opacity-40 mb-3">
              <span>Bulan Baru</span>
              <span className="text-primary opacity-100 uppercase tracking-widest font-black">
                Progress Siklus
              </span>
              <span>Bulan Purnama</span>
            </div>
            <div className="h-3 w-full bg-gray-100 dark:bg-white/5 rounded-full overflow-hidden">
              <div
                className="h-full bg-primary shadow-[0_0_15px_rgba(16,185,129,0.5)] rounded-full transition-all duration-1000"
                style={{ width: `${illumination}%` }}
              ></div>
            </div>
          </div>
        </div>
      </div>

      {/* Secondary Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <StatCard
          icon="vertical_align_top"
          label="Altitude Bulan"
          value="+12° 45'"
          desc="Tinggi di atas ufuk saat matahari terbenam."
          footer="Memenuhi kriteria MABIMS (> 3°)"
          isSuccess
        />
        <StatCard
          icon="schedule"
          label="Lag Time"
          value="48 Menit"
          desc="Durasi bulan di atas ufuk setelah sunset."
          footer="Optimal untuk observasi"
          progress={75}
        />
      </div>
    </div>
  );
}
