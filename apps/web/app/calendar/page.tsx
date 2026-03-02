"use client";

import { useState, useEffect } from "react";
import { useMounted } from "@/hooks/use-mounted";
import type { Method as HijriMethod } from "@/types/hijri";
import {
  HIJRI_MONTHS_INDONESIA_GRAMMAR as HIJRI_MONTHS,
  GREGORIAN_MONTHS,
} from "@/lib/constants";
import type { DateSystem, ViewMode } from "@/types/calendar";
import { CalendarHeader } from "@/components/calendar/CalendarHeader";
import { GlobalSettings } from "@/components/calendar/GlobalSettings";
import { CalendarSidebar } from "@/components/calendar/CalendarSidebar";
import { MonthGrid } from "@/components/calendar/MonthGrid";
import { YearlyMonthCard } from "@/components/calendar/YearlyMonthCard";
import { useCalendar } from "@/hooks/use-calendar";

export default function CalendarPage() {
  const mounted = useMounted();

  const [viewMode, setViewMode] = useState<ViewMode>("monthly");
  const [selectedMonth, setSelectedMonth] = useState(9);
  const [year, setYear] = useState(1447);
  const [method, setMethod] = useState<HijriMethod>("umm_al_qura");
  const [dateSystem, setDateSystem] = useState<DateSystem>("hijri");

  // Lokasi default Jakarta — bisa di-extend jadi geolocation nanti
  const lat = -6.2;
  const lon = 106.845;

  const { months, loading, error, today } = useCalendar(
    year,
    dateSystem,
    method,
    lat,
    lon,
  );

  // --- Auto-adjust saat pindah mode ---
  useEffect(() => {
    if (dateSystem === "gregorian") {
      setYear(new Date().getFullYear());
      setSelectedMonth(new Date().getMonth() + 1);
    } else {
      setYear(1447);
      setSelectedMonth(9);
    }
  }, [dateSystem]);

  // --- Auto-generate: navigasi bulan melewati batas tahun ---
  const handleMonthNav = (direction: "prev" | "next") => {
    const maxMonth = dateSystem === "hijri" ? 12 : 12;
    if (direction === "next") {
      if (selectedMonth >= maxMonth) {
        setYear((y) => y + 1);
        setSelectedMonth(1);
      } else {
        setSelectedMonth((m) => m + 1);
      }
    } else {
      if (selectedMonth <= 1) {
        setYear((y) => y - 1);
        setSelectedMonth(maxMonth);
      } else {
        setSelectedMonth((m) => m - 1);
      }
    }
  };

  const currentMonthData = months.find((m) => m.month_id === selectedMonth);

  const monthList =
    dateSystem === "hijri"
      ? HIJRI_MONTHS.map((m) => ({ id: m.id, name: m.name, desc: m.desc, isSpecial: "isSpecial" in m ? m.isSpecial : false }))
      : GREGORIAN_MONTHS.map((m) => ({ id: m.id, name: m.name, desc: m.desc, isSpecial: false }));

  if (!mounted)
    return (
      <div className="min-h-screen bg-background-light dark:bg-background-dark" />
    );

  return (
    <main className="min-h-screen bg-background-light dark:bg-background-dark transition-colors duration-500 pb-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-10">
        <CalendarHeader
          year={year}
          dateSystem={dateSystem}
          method={method}
          viewMode={viewMode}
          onYearChange={setYear}
          onViewModeChange={setViewMode}
        />

        <GlobalSettings
          method={method}
          onMethodChange={setMethod}
          dateSystem={dateSystem}
          onDateSystemChange={setDateSystem}
        />

        {/* --- MAIN LAYOUT --- */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10">
          <CalendarSidebar
            viewMode={viewMode}
            selectedMonth={selectedMonth}
            onMonthSelect={setSelectedMonth}
            dateSystem={dateSystem}
            onMonthNav={handleMonthNav}
          />

          {/* MAIN CONTENT AREA */}
          <div className="lg:col-span-9">
            {error && (
              <div className="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 p-4 rounded-2xl text-sm mb-6">
                {error}
              </div>
            )}

            {viewMode === "monthly" ? (
              <MonthGrid
                monthData={currentMonthData}
                year={year}
                system={dateSystem}
                loading={loading}
                todayDay={
                  today && today.year === year && today.month === currentMonthData?.month_id
                    ? today.day
                    : undefined
                }
              />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-8 duration-1000">
                {monthList.map((m) => (
                  <YearlyMonthCard
                    key={m.id}
                    month={m}
                    dateSystem={dateSystem}
                    onClick={() => {
                      setSelectedMonth(m.id);
                      setViewMode("monthly");
                    }}
                  />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
