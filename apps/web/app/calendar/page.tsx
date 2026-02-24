"use client";

import { useState } from "react";
import { useMounted } from "@/hooks/use-mounted";
import type { Method as HijriMethod } from "@/types/hijri";
import { HIJRI_MONTHS_INDONESIA_GRAMMAR as HIJRI_MONTHS } from "@/lib/constants";
import type { DateSystem, ViewMode } from "@/types/calendar";
import { CalendarHeader } from "@/components/calendar/CalendarHeader";
import { GlobalSettings } from "@/components/calendar/GlobalSettings";
import { CalendarSidebar } from "@/components/calendar/CalendarSidebar";
import { MonthGrid } from "@/components/calendar/MonthGrid";
import { YearlyMonthCard } from "@/components/calendar/YearlyMonthCard";

export default function CalendarPage() {
  const mounted = useMounted();

  const [viewMode, setViewMode] = useState<ViewMode>("monthly");
  const [selectedMonth, setSelectedMonth] = useState(9);
  const [year, setYear] = useState(1447);
  const [method, setMethod] = useState<HijriMethod>("umm_al_qura");
  const [dateSystem, setDateSystem] = useState<DateSystem>("hijri");

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
          />

          {/* MAIN CONTENT AREA */}
          <div className="lg:col-span-9">
            {viewMode === "monthly" ? (
              <MonthGrid
                monthId={selectedMonth}
                year={year}
                system={dateSystem}
              />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-8 duration-1000">
                {HIJRI_MONTHS.map((m) => (
                  <YearlyMonthCard
                    key={m.id}
                    month={m}
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
