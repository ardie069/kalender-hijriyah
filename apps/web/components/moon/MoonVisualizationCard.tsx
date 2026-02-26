"use client";

import { ArrowUpFromLine, Radius } from "lucide-react";
import Metric from "./Metric";
import StatCard from "./StatCard";

interface MoonVisualizationProps {
  illumination: number;
  phaseName: string;
  age: number;
  distance: string;
  elongation: number;
  altitude: number;
  lat: number;
  is_rukyat_time: boolean;
}

export default function MoonVisualizationCard({
  illumination = 0,
  phaseName = "Loading...",
  age = 0,
  distance = "0",
  elongation = 0,
  altitude = 0,
  lat = 0,
  is_rukyat_time = false,
}: MoonVisualizationProps) {
  const isWaning =
    phaseName.toLowerCase().includes("waning") ||
    phaseName.toLowerCase().includes("terakhir") ||
    phaseName.toLowerCase().includes("third");

  const isWaxing = !isWaning;

  const isNorth = lat > 0;

  const lightOnRight = (isNorth && isWaxing) || (!isNorth && isWaning);

  const percentage = illumination / 100;

  return (
    <div className="space-y-8 animate-in fade-in duration-1000">
      <div className="bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl p-5 sm:p-8 md:p-12 rounded-2xl sm:rounded-3xl border border-white/40 dark:border-white/5 shadow-soft overflow-hidden relative group transition-all duration-500">
        <div className="relative z-10">
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 mb-6 sm:mb-12">
            <div>
              <h2 className="text-lg sm:text-xl md:text-2xl font-black tracking-tight text-gray-900 dark:text-white">
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

          <div className="flex flex-col md:flex-row items-center gap-8 sm:gap-12 md:gap-16 justify-center">
            <div className="relative w-40 h-40 sm:w-56 sm:h-56 md:w-64 md:h-64 shrink-0">
              <div className="absolute inset-0 bg-primary/20 blur-[80px] rounded-full animate-pulse"></div>

              <div className="relative w-full h-full rounded-full bg-gray-950 shadow-2xl overflow-hidden border border-white/10 flex items-center justify-center">
                <div className="absolute inset-0 bg-gray-900 shadow-inner"></div>

                {/* Sisi Terang */}
                <div
                  className={`absolute w-1/2 h-full bg-gray-100 shadow-[0_0_30px_rgba(255,255,255,0.3)] transition-all duration-1000
                    ${lightOnRight ? "right-0" : "left-0"}`}
                ></div>

                {/* Elips Terminator: Inti dari visualisasi fase */}
                <div
                  className={`absolute h-[101%] rounded-full transition-all duration-1000 ease-in-out
                    ${illumination < 50 ? "bg-gray-900" : "bg-gray-100"}`}
                  style={{
                    width: "100%",
                    transform: `scaleX(${Math.abs(1 - percentage * 2)})`,
                    zIndex: 10,
                  }}
                ></div>

                <div className="absolute inset-0 opacity-10 pointer-events-none bg-[url('https://www.transparenttextures.com/patterns/pollen.png')] mix-blend-overlay"></div>
                <div className="absolute inset-0 bg-linear-to-tr from-black/40 via-transparent to-white/10 rounded-full"></div>
              </div>

              <div className="absolute bottom-4 right-4 z-20 flex items-center gap-2 bg-white/90 dark:bg-black/50 backdrop-blur px-3 py-1.5 rounded-full shadow-sm border border-gray-100 dark:border-gray-800">
                <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse"></span>
                <span className="text-[10px] font-black uppercase tracking-widest">
                  Live
                </span>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-y-6 gap-x-4 sm:gap-y-10 sm:gap-x-8 w-full max-w-sm">
              <Metric
                label="Iluminasi"
                value={`${illumination.toFixed(1)}%`}
                sub="Light Density"
              />
              <Metric
                label="Umur Bulan"
                value={`${age.toFixed(1)} Hari`}
                sub="Cycle Age"
              />
              <Metric
                label="Jarak"
                value={`${distance} km`}
                sub="Earth-Moon Distance"
              />
              <Metric
                label="Elongasi"
                value={`${elongation.toFixed(1)}°`}
                sub="Sun-Moon Angle"
              />
            </div>
          </div>

          <div className="mt-8 sm:mt-16 pt-6 sm:pt-8 border-t border-gray-100 dark:border-white/5">
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

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <StatCard
          icon={<ArrowUpFromLine />}
          label="Altitude Bulan"
          value={`${altitude.toFixed(2)}°`}
          desc={
            altitude > 0
              ? "Tinggi di atas ufuk saat pengamatan."
              : "Bulan berada di bawah ufuk"
          }
          footer={
            is_rukyat_time
              ? altitude > 3
                ? "Lolos Kriteria MABIMS"
                : "Belum Lolos Kriteria"
              : "Posisi Real-time"
          }
          isSuccess={is_rukyat_time && altitude > 3}
        />

        <StatCard
          icon={<Radius />}
          label={is_rukyat_time ? "Status Visibilitas" : "Elongasi Sudut"}
          value={
            is_rukyat_time
              ? altitude > 3 && elongation > 6.4
                ? "Visible"
                : "Not Visible"
              : `${elongation.toFixed(2)}°`
          }
          desc={
            is_rukyat_time
              ? "Parameter Altitude & Elongasi"
              : "Jarak Sudut Bulan-Matahari"
          }
          footer="Kalkulasi Toposentris"
          progress={
            is_rukyat_time ? (altitude > 3 && elongation > 6.4 ? 100 : 30) : 100
          }
        />
      </div>
    </div>
  );
}
