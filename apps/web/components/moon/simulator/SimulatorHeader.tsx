"use client";

import { Sunset, Timer } from "lucide-react";

interface SimulatorHeaderProps {
  sunsetTime: string;
  moonAltitudeAtSunset: number;
  elongation: number;
}

export default function SimulatorHeader({
  sunsetTime,
  moonAltitudeAtSunset,
  elongation,
}: SimulatorHeaderProps) {
  const isMoonAboveHorizon = moonAltitudeAtSunset > 0;

  const dynamicMaxRange = isMoonAboveHorizon
    ? Math.min(120, Math.max(45, Math.ceil(elongation * 5)))
    : 30;

  return (
    <div className="flex flex-col md:flex-row md:items-end justify-between gap-6">
      <div className="space-y-4">
        <h1 className="text-xl sm:text-2xl md:text-3xl font-black tracking-tighter text-gray-900 dark:text-white">
          Simulator Perbandingan
        </h1>
        <div className="flex flex-wrap items-center gap-3">
          <div className="px-3 py-1.5 bg-orange-500/10 border border-orange-500/20 rounded-full flex items-center gap-2">
            <Sunset className="text-orange-500 w-4 h-4" />
            <p className="text-[10px] font-black text-orange-600 dark:text-orange-400 uppercase tracking-widest">
              Sunset: {sunsetTime} WIB
            </p>
          </div>
          <div className="px-3 py-1.5 bg-blue-500/10 border border-blue-500/20 rounded-full flex items-center gap-2">
            <Timer className="text-blue-500 w-3.5 h-3.5" />
            <p className="text-[9px] font-black text-blue-600 dark:text-blue-400 uppercase tracking-widest">
              Lag Time: ~{dynamicMaxRange}m
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
