"use client";

import { Sun } from "lucide-react";
import { useMemo } from "react";

interface SkyBoxProps {
  alt: number;
  azimuthDiff?: number;
  offset: number;
  isOrange: boolean;
  label: string;
  icon: React.ReactNode;
}

export default function SkyBox({
  alt,
  azimuthDiff = 0,
  offset,
  isOrange,
  label,
  icon,
}: SkyBoxProps) {
  // --- 1. KONFIGURASI ULTRA-WIDE ---
  const DEG_TO_PCT_X = 0.65;
  const DEG_TO_PCT_Y = 2.0;
  const HORIZON_Y = 15;

  // --- 2. LOGIKA POSISI ---
  const sunAlt = 0 - offset * 0.25;
  const sunX = 50;
  const sunY = HORIZON_Y + sunAlt * DEG_TO_PCT_Y;

  const moonAlt = alt - offset * 0.25;
  const moonX = 50 + azimuthDiff * DEG_TO_PCT_X;
  const moonY = HORIZON_Y + moonAlt * DEG_TO_PCT_Y;

  const rotationAngle = useMemo(() => {
    const dx = moonX - sunX;
    const dy = moonY - sunY;
    return (Math.atan2(dy, dx) * 180) / Math.PI + 90;
  }, [moonX, sunX, moonY, sunY]);

  const skyOpacity = Math.min(0.95, offset / 45);

  return (
    <div className="relative z-0 w-full h-full overflow-hidden bg-[#020617] select-none rounded-4xl border border-white/5 shadow-inner">
      {/* LAYER 1: ATMOSFER */}
      <div
        className={`absolute inset-0 bg-linear-to-b transition-colors duration-1000 ${
          isOrange
            ? "from-[#1e293b] via-[#334155] to-[#f97316]"
            : "from-[#0f172a] via-[#1e293b] to-[#c2410c]"
        }`}
      />
      <div
        className="absolute inset-0 bg-[#020617] transition-opacity duration-1000"
        style={{ opacity: skyOpacity }}
      />

      {/* LAYER 2: WEST MARKER (W) */}
      <div
        className="absolute z-20 flex flex-col items-center"
        style={{ bottom: "15%", left: "50%", transform: "translateX(-50%)" }}
      >
        <div className="bg-blue-600/30 backdrop-blur-md border border-blue-400/30 p-1 rounded-sm rotate-45">
          <span className="block -rotate-45 text-[7px] font-black text-white/90">
            W
          </span>
        </div>
        <div className="w-px h-12 bg-linear-to-b from-blue-400/20 to-transparent" />
      </div>

      {/* LAYER 3: UFUK (GROUND) */}
      <div className="absolute bottom-0 w-full h-[15%] bg-[#08120a] z-40 border-t border-emerald-900/20" />

      {/* LAYER 4: MATAHARI (Ukuran diperkecil agar proporsional) */}
      <div
        className="absolute z-10 transition-all duration-300 ease-out"
        style={{
          bottom: `${sunY}%`,
          left: `${sunX}%`,
          transform: "translate(-50%, 50%)",
        }}
      >
        <div className="relative flex items-center justify-center">
          {/* Pendaran disesuaikan dengan skala wide */}
          <div className="absolute w-64 h-64 bg-orange-600 rounded-full blur-[60px] opacity-20" />
          <Sun className="w-5 h-5 text-orange-100 relative z-10 drop-shadow-[0_0_8px_rgba(255,165,0,0.6)]" />
        </div>
      </div>

      {/* LAYER 5: HILAL (Ukuran ciut mengikuti skala wide) */}
      <div
        className="absolute z-50 transition-all duration-500 ease-out"
        style={{
          bottom: `${moonY}%`,
          left: `${moonX}%`,
          transform: "translate(-50%, 50%)",
          opacity: moonAlt < -5 ? 0 : 1,
        }}
      >
        <div className="relative flex flex-col items-center">
          {/* Moon size dikurangi: w-12/16 -> w-6/8 */}
          <div
            className="w-6 h-6 sm:w-8 sm:h-8 rounded-full border-r-[2.5px] border-b-[0.5px] border-white shadow-[0_0_15px_rgba(255,255,255,0.6)]"
            style={{ transform: `rotate(${rotationAngle}deg)` }}
          />

          {/* Tag Telemetry lebih ramping */}
          <div className="absolute -top-10 flex flex-col items-center">
            <div className="bg-black/60 backdrop-blur-lg px-2 py-0.5 rounded-lg border border-white/5 shadow-2xl">
              <p className="text-[7px] font-bold text-white/80 tabular-nums uppercase tracking-tighter">
                Alt: {moonAlt.toFixed(1)}°
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* LAYER 6: UI INFO */}
      <div className="absolute top-4 left-4 z-50 flex items-center gap-2 bg-black/40 backdrop-blur-xl px-3 py-1.5 rounded-xl border border-white/10">
        <div className="text-white scale-75 opacity-70">{icon}</div>
        <span className="text-[9px] font-black text-white/90 uppercase tracking-widest">
          {label}
        </span>
      </div>
    </div>
  );
}
