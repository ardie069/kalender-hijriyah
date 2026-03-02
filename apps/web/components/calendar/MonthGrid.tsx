"use client";

import { useMemo } from "react";
import { addDays, format, parseISO } from "date-fns";
import { id } from "date-fns/locale";
import type { UnifiedMonthData, DateSystem } from "@/types/calendar";

const WEEK_DAYS = [
  "Ahad",
  "Senin",
  "Selasa",
  "Rabu",
  "Kamis",
  "Jum'at",
  "Sabtu",
];

interface MonthGridProps {
  monthData: UnifiedMonthData | undefined;
  year: number;
  system: DateSystem;
  loading?: boolean;
  todayDay?: number;
}

export function MonthGrid({
  monthData,
  year,
  system,
  loading,
  todayDay,
}: MonthGridProps) {
  const { days, offset, startDate } = useMemo(() => {
    if (!monthData) return { days: [], offset: 0, startDate: null };
    const gridOffset = (monthData.day_1_weekday + 1) % 7;

    const start = monthData.start_gregorian
      ? parseISO(monthData.start_gregorian)
      : null;

    return {
      days: Array.from({ length: monthData.total_days }),
      offset: gridOffset,
      startDate: start,
    };
  }, [monthData]);

  if (loading || !monthData) {
    return (
      <div className="bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl rounded-3xl p-20 text-center border border-white/10 shadow-soft">
        <div className="flex flex-col items-center gap-4">
          <div className="w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin" />
          <p className="text-[10px] font-black uppercase tracking-[0.4em] opacity-40">
            {system === "hijri"
              ? "Scanning Lunar Orbit..."
              : "Loading Calendar..."}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl rounded-3xl border border-white/40 dark:border-white/5 shadow-soft p-6 md:p-12 animate-in zoom-in-95 duration-700">
      <header className="flex justify-between items-end mb-12">
        <div>
          <h2 className="text-4xl sm:text-5xl font-black tracking-tighter text-gray-900 dark:text-white uppercase">
            {monthData.month_name}
          </h2>
          <p className="text-[10px] font-black text-primary uppercase tracking-[0.5em] mt-2">
            {system === "hijri"
              ? `Projection: Year ${year} H • ${monthData.total_days} Days`
              : `${monthData.month_name} ${year} • ${monthData.total_days} Hari`}
          </p>
        </div>
        <div className="hidden sm:flex gap-1.5 mb-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className={`h-1.5 w-1.5 rounded-full ${i === 0 ? "bg-primary animate-pulse" : "bg-primary/20"}`}
            />
          ))}
        </div>
      </header>

      <div className="grid grid-cols-7 gap-2 md:gap-4">
        {/* Render Header Hari */}
        {WEEK_DAYS.map((d, i) => (
          <div
            key={d}
            className={`text-center py-3 text-[10px] font-black uppercase tracking-[0.3em] opacity-30 ${i === 5 ? "text-primary opacity-100" : ""}`}
          >
            {d}
          </div>
        ))}

        {/* 1. KOTAK KOSONG (OFFSET) */}
        {Array.from({ length: offset }).map((_, i) => (
          <div
            key={`empty-${i}`}
            className="aspect-square sm:aspect-auto sm:min-h-30 opacity-0"
          />
        ))}

        {/* 2. KOTAK HARI (ACTIVE) */}
        {days.map((_, i) => {
          const day = i + 1;
          const isJumat = (i + offset) % 7 === 5;
          const isToday = todayDay === day;
          const gregorianDate = startDate ? addDays(startDate, i) : null;

          return (
            <div
              key={day}
              className={`aspect-square sm:aspect-auto sm:min-h-30 p-5 rounded-[2.2rem] border-2 transition-all duration-300 flex flex-col items-start group cursor-pointer
              ${
                isToday
                  ? "border-primary bg-primary/10 shadow-[0_0_16px_rgba(16,185,129,0.25)] ring-1 ring-primary/30 scale-[1.04]"
                  : isJumat
                    ? "border-primary/20 bg-primary/5"
                    : "border-white/20 dark:border-white/5 bg-white/20 dark:bg-white/5 hover:border-primary/40"
              }
            `}
            >
              <span
                className={`text-2xl sm:text-4xl font-black tabular-nums transition-all ${isToday ? "text-primary" : isJumat ? "text-primary" : "text-gray-900 dark:text-white opacity-80"}`}
              >
                {day}
              </span>

              {/* Mapping Gregorian — hanya di mode Hijriyah */}
              {system === "hijri" && gregorianDate && (
                <span className="text-[8px] sm:text-[10px] font-bold uppercase tracking-tighter opacity-20 group-hover:opacity-60 mt-1">
                  {format(gregorianDate, "d MMM", { locale: id })}
                </span>
              )}

              {/* Indikator Fase & Event — hanya di mode Hijriyah */}
              {system === "hijri" && (
                <div className="mt-auto flex items-center gap-1.5">
                  {day === 1 && (
                    <div
                      className="h-1 w-6 bg-indigo-500 rounded-full shadow-[0_0_8px_rgba(99,102,241,0.5)]"
                      title="New Moon"
                    />
                  )}
                  {day === 15 && (
                    <div
                      className="h-3 w-3 rounded-full border-2 border-gray-400 dark:border-gray-500 bg-white shadow-sm"
                      title="Full Moon"
                    />
                  )}
                  {day === monthData.total_days && (
                    <div
                      className={`h-3 w-3 rounded-full ${monthData.total_days === 30 ? "bg-gray-800 dark:bg-gray-200" : "bg-gray-800"}`}
                      title={
                        monthData.total_days === 30
                          ? "Istikmal (30 Days)"
                          : "Rukyat (29 Days)"
                      }
                    />
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
