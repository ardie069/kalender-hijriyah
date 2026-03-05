"use client";

import { useState, useEffect, useMemo } from "react";
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
  const [dateSystem, setDateSystemState] = useState<DateSystem>("hijri");
  const [method, setMethodState] = useState<HijriMethod>("umm_al_qura");

  const [year, setYear] = useState<number | null>(null);
  const [selectedMonth, setSelectedMonth] = useState<number | null>(null);
  const [hasSynced, setHasSynced] = useState(false);

  const setDateSystem = (val: DateSystem) => {
    setHasSynced(false);
    setYear(null);
    setSelectedMonth(null);
    setDateSystemState(val);
  };

  const setMethod = (val: HijriMethod) => {
    setHasSynced(false);
    setYear(null);
    setSelectedMonth(null);
    setMethodState(val);
  };

  const [location, setLocation] = useState<{ lat: number; lon: number } | null>(
    null,
  );

  useEffect(() => {
    const fetchIPLocation = async () => {
      try {
        const response = await fetch("https://ipapi.co/json/");
        const data = await response.json();
        if (data.latitude && data.longitude) {
          setLocation({ lat: data.latitude, lon: data.longitude });
        }
      } catch {
        setLocation({ lat: -6.2, lon: 106.845 });
      }
    };

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (pos) =>
          setLocation({ lat: pos.coords.latitude, lon: pos.coords.longitude }),
        () => fetchIPLocation(),
        { timeout: 5000 },
      );
    } else {
      fetchIPLocation();
    }
  }, []);

  const initialYear = useMemo(() => {
    try {
      const now = new Date();
      const gregYear = now.getFullYear();

      if (dateSystem === "gregorian") return gregYear;

      return Math.floor((gregYear - 622) * 1.0307);
    } catch (err) {
      console.error("Error calculating initial year:", err);
      return new Date().getFullYear();
    }
  }, [dateSystem]);

  const { months, loading, error, today } = useCalendar(
    year,
    dateSystem,
    method,
    location?.lat ?? null,
    location?.lon ?? null,
  );

  useEffect(() => {
    if (today && today.system === dateSystem && !hasSynced) {
      const timer = setTimeout(() => {
        setYear(today.year);
        setSelectedMonth(today.month);
        setHasSynced(true);
      }, 0);
      return () => clearTimeout(timer);
    }
  }, [today, dateSystem, hasSynced]);

  // ─── 4. DISPLAY LOGIC ───
  const activeToday = today?.system === dateSystem ? today : null;
  const displayYear = year ?? activeToday?.year ?? initialYear;
  const initialMonth = useMemo(() => {
    const now = new Date();

    if (dateSystem === "gregorian") {
      return now.getMonth() + 1;
    }

    // ─── DYNAMIC HIJRI MONTH DETECTION ───
    try {
      const formatter = new Intl.DateTimeFormat("id-ID-u-ca-islamic-uma", {
        month: "numeric",
      });

      const hijriMonth = parseInt(formatter.format(now).replace(/\D/g, ""));

      if (hijriMonth >= 1 && hijriMonth <= 12) {
        return hijriMonth;
      }
    } catch (e) {
      console.error("Failed to detect Hijri month using Intl:", e);
    }

    const startOfHijriYear = new Date(now.getFullYear(), 5, 26);
    const diffDays = Math.floor(
      (now.getTime() - startOfHijriYear.getTime()) / (1000 * 60 * 60 * 24),
    );
    const estimatedMonth = Math.floor(diffDays / 29.5) + 1;

    return Math.min(Math.max(estimatedMonth, 1), 12);
  }, [dateSystem]);

  const currentMonthId = selectedMonth ?? activeToday?.month ?? initialMonth;

  const currentMonthData = months.find((m) => m.month_id === currentMonthId);

  const monthList = useMemo(() => {
    return dateSystem === "hijri"
      ? HIJRI_MONTHS.map((m) => ({
        id: m.id,
        name: m.name,
        desc: m.desc,
        isSpecial: "isSpecial" in m ? m.isSpecial : false,
      }))
      : GREGORIAN_MONTHS.map((m) => ({
        id: m.id,
        name: m.name,
        desc: m.desc,
        isSpecial: false,
      }));
  }, [dateSystem]);

  if (!mounted)
    return (
      <div className="min-h-screen bg-background-light dark:bg-background-dark" />
    );

  return (
    <main className="min-h-screen bg-background-light dark:bg-background-dark transition-colors duration-500 pb-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-10">
        <CalendarHeader
          year={displayYear}
          dateSystem={dateSystem}
          method={method}
          viewMode={viewMode}
          onViewModeChange={setViewMode}
        />

        <GlobalSettings
          method={method}
          onMethodChange={setMethod}
          dateSystem={dateSystem}
          onDateSystemChange={setDateSystem}
        />

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10">
          <CalendarSidebar
            viewMode={viewMode}
            selectedMonth={currentMonthId}
            onMonthSelect={setSelectedMonth}
            dateSystem={dateSystem}
          />

          <div className="lg:col-span-9">
            {error && (
              <div className="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 p-4 rounded-2xl text-sm mb-6">
                {error}
              </div>
            )}

            {viewMode === "monthly" ? (
              <MonthGrid
                monthData={currentMonthData}
                year={displayYear}
                system={dateSystem}
                loading={loading || !hasSynced}
                todayDay={today?.day}
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
