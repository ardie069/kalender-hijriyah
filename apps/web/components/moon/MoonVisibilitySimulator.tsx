"use client";

import { useState, useMemo } from "react";
import {
  ArrowUpToLine,
  Ruler,
  Sun,
  Moon,
  Eye,
  EyeOff,
  Navigation,
  Sparkles,
} from "lucide-react";
import { MoonTelemetry, MoonStatus } from "@/types/moon";
import MiniStat from "./simulator/MiniStat";
import SimulatorHeader from "./simulator/SimulatorHeader";
import SimulationCard from "./simulator/SimulationCard";
import DetailedDataTable from "./simulator/DetailedDataTable";

interface MoonVisibilitySimulatorProps {
  telemetry: MoonTelemetry;
  status: MoonStatus;
}

export default function MoonVisibilitySimulator({
  telemetry,
  status,
}: MoonVisibilitySimulatorProps) {
  const [offsets, setOffsets] = useState({ day29: 15, day30: 15, day1: 15 });

  const updateOffset = (key: keyof typeof offsets, val: number) =>
    setOffsets((prev) => ({ ...prev, [key]: val }));

  const simulationDays = useMemo(() => {
    const days = [
      {
        id: "day29" as const,
        title: "Day 29: Pre-Conjunction",
        desc: "Kondisi akhir bulan saat konjungsi belum sempurna atau hilal terlalu rendah.",
        alt: telemetry.altitude - 2.5,
        elong: 2.11,
        isOrange: false,
        icon: <EyeOff className="w-4 h-4" />,
        badge: null,
        show: !status.is_mabims_met,
      },
      {
        id: "day30" as const,
        title: status.is_mabims_met
          ? "Observation Day"
          : "Day 30: Critical Observation",
        desc: "Masa kritis rukyatul hilal. Geser slider untuk melihat simulasi terbenam.",
        alt: telemetry.altitude - 0.5,
        elong: telemetry.elongation - 0.5,
        isOrange: true,
        icon: <Navigation className="w-4 h-4 rotate-45" />,
        badge: null,
        show: true,
      },
      {
        id: "day1" as const,
        title: "Day 1: New Moon Cycle",
        desc: "Kondisi bulan baru yang mulai terlihat jelas. Simulasi kegelapan langit.",
        alt: telemetry.altitude,
        elong: telemetry.elongation,
        isOrange: false,
        icon: <Eye className="w-4 h-4" />,
        badge: (
          <span
            className={`text-[10px] px-3 py-1 rounded-full font-black uppercase ${status.is_mabims_met ? "bg-primary text-white" : "bg-white/10 text-gray-400"}`}
          >
            {status.is_mabims_met ? "Visible" : "Wait Phase"}
          </span>
        ),
        show: true,
      },
    ];
    return days.filter((day) => day.show);
  }, [telemetry, status.is_mabims_met]);

  return (
    <div className="space-y-8 animate-in fade-in duration-700">
      <SimulatorHeader
        sunsetTime="17:54"
        moonAltitudeAtSunset={telemetry.altitude}
        elongation={telemetry.elongation}
      />

      <div className="flex flex-col gap-12">
        {simulationDays.map((day) => (
          <SimulationCard
            key={day.id}
            title={day.title}
            description={day.desc}
            badge={day.badge}
            className={
              day.id === "day1" && status.is_mabims_met
                ? "border-primary/40 bg-white dark:bg-card-dark"
                : "bg-white/40 border-white/10"
            }
            visual={
              <SkyBox
                alt={day.alt}
                offset={offsets[day.id]}
                isOrange={day.isOrange}
                label={`${day.title} View`}
                icon={day.icon}
              />
            }
            stats={
              <>
                <MiniStat
                  icon={
                    <ArrowUpToLine
                      className={`w-4 h-4 ${day.id === "day1" ? "text-primary" : ""}`}
                    />
                  }
                  label="Altitude"
                  value={`${(day.alt - offsets[day.id] * 0.25).toFixed(2)}°`}
                />
                <MiniStat
                  icon={
                    <Ruler
                      className={`w-4 h-4 ${day.id === "day1" ? "text-primary" : ""}`}
                    />
                  }
                  label="Elongasi"
                  value={`${day.elong.toFixed(2)}°`}
                />
              </>
            }
            controls={
              <LocalSlider
                value={offsets[day.id]}
                onChange={(val) => updateOffset(day.id, val)}
                alt={day.alt}
                elong={day.elong}
              />
            }
          />
        ))}
      </div>

      <DetailedDataTable telemetry={telemetry} />
    </div>
  );
}

interface LocalSliderProps {
  value: number;
  onChange: (val: number) => void;
  alt: number;
  elong: number;
}

function LocalSlider({ value, onChange, alt, elong }: LocalSliderProps) {
  const isAbove = alt > 0;
  const maxRange = Math.min(120, Math.max(45, Math.ceil(elong * 5)));

  return (
    <div className="space-y-3">
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-2">
          <div
            className={`p-1.5 rounded-lg ${isAbove ? "bg-primary/20 text-primary" : "bg-red-500/20 text-red-500"}`}
          >
            {isAbove ? (
              <Sun className="w-3 h-3" />
            ) : (
              <Moon className="w-3 h-3" />
            )}
          </div>
          <span className="text-[10px] font-black uppercase text-gray-400">
            T + {value} Menit
          </span>
        </div>
        {value >= 15 && value <= 25 && isAbove && (
          <span className="text-[9px] font-black text-primary flex items-center gap-1 uppercase animate-pulse">
            <Sparkles className="w-3 h-3" /> Best View
          </span>
        )}
      </div>
      <input
        type="range"
        min="0"
        max={maxRange}
        value={value}
        onChange={(e) => onChange(parseInt(e.target.value))}
        className={`w-full h-2 rounded-lg appearance-none cursor-pointer accent-primary bg-gray-200 dark:bg-gray-800`}
      />
    </div>
  );
}

interface SkyBoxProps {
  alt: number;
  offset: number;
  isOrange: boolean;
  label: string;
  icon: React.ReactNode;
}

function SkyBox({ alt, offset, isOrange, label, icon }: SkyBoxProps) {
  const skyOpacity = Math.min(0.9, offset / 50);
  const ySun = 15 + (0 - offset * 0.25) * 3;
  const yMoon = 15 + (alt - offset * 0.25) * 3;

  return (
    <div className="relative z-0 w-full h-full overflow-hidden bg-[#0f172a]">
      <div
        className={`absolute inset-0 bg-linear-to-b ${isOrange ? "from-[#0f172a] via-[#1e293b] to-[#c2410c]" : "from-[#1e293b] via-[#334155] to-[#f97316]"}`}
      />
      <div
        className="absolute inset-0 bg-[#020617] transition-opacity duration-700"
        style={{ opacity: skyOpacity }}
      />
      <div className="absolute bottom-0 w-full h-[15%] bg-[#112116] z-10 border-t border-emerald-900/50" />

      <div className="absolute top-3 left-3 sm:top-6 sm:left-6 z-30 bg-black/60 px-2 py-1 sm:px-4 sm:py-2 rounded-full text-[8px] sm:text-[10px] font-black text-white uppercase tracking-[0.15em] sm:tracking-[0.2em] flex items-center gap-1.5 sm:gap-2 border border-white/10">
        {icon} {label}
      </div>

      <div
        className="absolute left-1/2 -translate-x-1/2 z-0 transition-all duration-300"
        style={{ bottom: `${ySun}%` }}
      >
        <div className="w-20 h-20 sm:w-32 sm:h-32 bg-orange-600 rounded-full blur-[40px] sm:blur-[60px] opacity-40" />
        <Sun className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 sm:w-6 sm:h-6 text-orange-200" />
      </div>

      <div
        className="absolute left-[45%] z-20 transition-all duration-300"
        style={{ bottom: `${yMoon}%` }}
      >
        <div
          className={`w-8 h-8 sm:w-12 sm:h-12 rounded-full border-r-[3px] border-b-[0.5px] border-white rotate-35 shadow-[0_0_20px_white] transition-opacity ${alt - offset * 0.25 < -1 ? "opacity-0" : "opacity-100"}`}
        />
      </div>
    </div>
  );
}
