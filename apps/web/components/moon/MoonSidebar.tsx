"use client";

import { MoonTelemetry } from "@/types/moon";
import DataRow from "./DataRow";

interface MoonSidebarProps {
  telemetry?: MoonTelemetry;
  method: "hisab" | "rukyat";
  onMethodChange: (method: "hisab" | "rukyat") => void;
}

export default function MoonSidebar({
  telemetry,
  method,
  onMethodChange,
}: MoonSidebarProps) {
  return (
    <aside className="lg:col-span-4 space-y-8 animate-in fade-in slide-in-from-right-10 duration-1000">
      {/* --- Calculation Method: Selection Logic --- */}
      <div className="bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-5 sm:p-8 rounded-2xl border border-white/40 dark:border-white/5 shadow-soft hover:shadow-primary/5 transition-all duration-500 group">
        <h3 className="text-[10px] font-black text-primary uppercase tracking-[0.4em] mb-8">
          Metode Kalkulasi
        </h3>

        <div className="space-y-4">
          {/* Active Method: Ephemeris (Hisab) */}
          <button
            onClick={() => onMethodChange("hisab")}
            className={`w-full p-5 flex items-center gap-4 transition-all duration-500 rounded-[2rem] text-left cursor-pointer
              ${
                method === "hisab"
                  ? "bg-primary text-white shadow-lg shadow-primary/20 scale-[1.02]"
                  : "bg-gray-100/50 dark:bg-white/5 opacity-50 grayscale hover:grayscale-0 hover:opacity-100 border border-transparent dark:border-white/5"
              }`}
          >
            <div
              className={`${method === "hisab" ? "bg-white/20" : "bg-gray-200 dark:bg-gray-800"} p-2.5 rounded-xl text-xl`}
            >
              🛰️
            </div>
            <div>
              <p
                className={`text-xs font-black uppercase tracking-tight ${method === "hisab" ? "text-white" : "text-gray-900 dark:text-white"}`}
              >
                Ephemeris Hisab
              </p>
              <p className="text-[9px] opacity-80 font-bold uppercase tracking-tighter">
                Algoritma Jean Meeus
              </p>
            </div>
          </button>

          {/* Inactive Method: Rukyat (Simulasi) */}
          <button
            onClick={() => onMethodChange("rukyat")}
            className={`w-full p-5 flex items-center gap-4 transition-all duration-500 rounded-[2rem] text-left cursor-pointer
              ${
                method === "rukyat"
                  ? "bg-orange-500 text-white shadow-lg shadow-orange-500/20 scale-[1.02]"
                  : "bg-gray-100/50 dark:bg-white/5 opacity-50 grayscale hover:grayscale-0 hover:opacity-100 border border-transparent dark:border-white/5"
              }`}
          >
            <div
              className={`${method === "rukyat" ? "bg-white/20" : "bg-gray-200 dark:bg-gray-800"} p-2.5 rounded-xl text-xl`}
            >
              🔭
            </div>
            <div>
              <p
                className={`text-xs font-black uppercase tracking-tight ${method === "rukyat" ? "text-white" : "text-gray-900 dark:text-white"}`}
              >
                Rukyatul Hilal
              </p>
              <p className="text-[9px] font-bold uppercase tracking-tighter">
                Observasi Visual
              </p>
            </div>
          </button>
        </div>
      </div>

      {/* --- Ephemeris Data List: Telemetry --- */}
      <div className="bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-5 sm:p-8 rounded-2xl border border-white/40 dark:border-white/5 shadow-soft transition-all duration-500 hover:border-primary/20 relative overflow-hidden">
        {/* Subtle Background Decoration */}
        <div className="absolute -bottom-4 -right-4 opacity-[0.03] dark:opacity-[0.07] pointer-events-none">
          <svg
            width="120"
            height="120"
            viewBox="0 0 100 100"
            fill="none"
            stroke="currentColor"
          >
            <path d="M10 90 Q 50 10 90 90" strokeWidth="2" />
          </svg>
        </div>

        <h3 className="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-[0.3em] mb-8 relative z-10">
          Data Ephemeris
        </h3>
        <ul className="space-y-6 relative z-10">
          <DataRow
            label="Azimuth Bulan"
            value={`${telemetry?.azimuth.toFixed(2) ?? "0.00"}°`}
            sub="Sudut dari Utara"
          />
          <DataRow
            label="Altitude (Tinggi)"
            value={`${telemetry?.altitude.toFixed(2) ?? "0.00"}°`}
            sub={(telemetry?.altitude ?? 0) > 0 ? "Di atas ufuk" : "Di bawah ufuk"}
          />
          <DataRow
            label="Jarak Bumi-Bulan"
            value={`${telemetry?.distance_km ?? "0"} km`}
            sub="Toposentris"
          />
        </ul>
      </div>

      {/* --- Wisdom Quote Section: The Spiritual Pulse --- */}
      <div className="p-5 sm:p-8 rounded-2xl bg-emerald-500/5 border border-primary/10 relative overflow-hidden group transition-all duration-700">
        {/* Giant Quote Mark Decoration */}
        <div className="absolute -top-4 -right-2 text-8xl font-black text-primary/10 opacity-0 group-hover:opacity-100 group-hover:translate-y-2 transition-all duration-1000 select-none">
          <q></q>
        </div>

        <div className="flex flex-col gap-4 relative z-10">
          <div className="h-px w-8 bg-primary/40" />
          <p className="text-sm font-bold italic text-emerald-800 dark:text-emerald-300 leading-relaxed">
            <q>
              Dan Dialah yang menjadikan matahari bersinar dan bulan bercahaya
              dan ditetapkan-Nya manzilah-manzilah (tempat-tempat) bagi
              perjalanan bulan itu...
            </q>
          </p>
          <p className="text-[9px] font-black text-primary mt-2 uppercase tracking-[0.3em] flex items-center gap-2">
            <span className="w-1 h-1 rounded-full bg-primary" />
            QS. Yunus: 5
          </p>
        </div>
      </div>
    </aside>
  );
}
