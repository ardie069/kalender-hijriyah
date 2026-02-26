"use client";

import { useState, useMemo } from "react";
import { ArrowUpToLine, Ruler, EyeOff, Navigation } from "lucide-react";
import { MoonTelemetry, MoonStatus, RukyatResponse } from "@/types/moon";
import MiniStat from "./simulator/MiniStat";
import SimulatorHeader from "./simulator/SimulatorHeader";
import SimulationCard from "./simulator/SimulationCard";
import DetailedDataTable from "./simulator/DetailedDataTable";
import SkyBox from "./simulator/SkyBox";
import LocalSlider from "./simulator/LocalSlider";

interface MoonVisibilitySimulatorProps {
  telemetry: MoonTelemetry;
  status: MoonStatus;
  rukyatData?: RukyatResponse | null;
}

export default function MoonVisibilitySimulator({
  telemetry,
  status,
  rukyatData,
}: MoonVisibilitySimulatorProps) {
  const [offsets, setOffsets] = useState({ day29: 15, day30: 15 });

  const updateOffset = (key: keyof typeof offsets, val: number) =>
    setOffsets((prev) => ({ ...prev, [key]: val }));

  const simulationDays = useMemo(() => {
    const d29Alt = rukyatData?.altitude_at_sunset ?? telemetry.altitude - 2.5;
    const d29Elon = rukyatData?.elongation_at_sunset ?? 0;

    const days = [
      {
        id: "day29" as const,
        title: "Day 29: Penentuan Hilal",
        desc: `Kondisi kritis penentuan bulan baru berdasarkan kriteria ${rukyatData?.criteria_used || "MABIMS"}.`,
        alt: d29Alt,
        elong: d29Elon,
        isOrange: false,
        icon: <EyeOff className="w-4 h-4" />,
        badge: (
          <span
            className={`text-[10px] px-2 py-0.5 rounded-lg font-black uppercase ${rukyatData?.is_visible ? "bg-primary text-white" : "bg-red-500/10 text-red-500"}`}
          >
            {rukyatData?.is_visible ? "Sighting Possible" : "Impossible"}
          </span>
        ),
        show: true,
      },
      {
        id: "day30" as const,
        title: "Observation Follow-up",
        desc: "Simulasi posisi hilal satu hari setelah hari rukyat (H+1 Isbat).",
        alt: telemetry.altitude,
        elong: telemetry.elongation,
        isOrange: true,
        icon: <Navigation className="w-4 h-4 rotate-45" />,
        badge: null,
        show: true,
      },
    ];
    return days.filter((day) => day.show);
  }, [telemetry, rukyatData]);

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
              day.id === "day29" && status.is_visible
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
                      className={`w-4 h-4 ${day.id === "day29" ? "text-primary" : ""}`}
                    />
                  }
                  label="Altitude"
                  value={`${(day.alt - offsets[day.id] * 0.25).toFixed(2)}°`}
                />
                <MiniStat
                  icon={
                    <Ruler
                      className={`w-4 h-4 ${day.id === "day29" ? "text-primary" : ""}`}
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
