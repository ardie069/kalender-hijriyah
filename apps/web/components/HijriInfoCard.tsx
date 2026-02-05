"use client";

import { useMemo } from "react";
import { useTheme } from "@/context/theme-context";
import type { HijriDate, Method } from "@/types/hijri";

interface HijriInfoCardProps {
  hijriDate: HijriDate;
  weton: string | null;
  method: Method;
}

const MONTHS = [
  "Muharam",
  "Safar",
  "Rabiulawal",
  "Rabiulakhir",
  "Jumadilawal",
  "Jumadilakhir",
  "Rajab",
  "Syakban",
  "Ramadan",
  "Syawal",
  "Zulkaidah",
  "Zulhijah",
];

export default function HijriInfoCard({
  hijriDate,
  weton,
  method,
}: HijriInfoCardProps) {
  const { darkMode } = useTheme();

  const themeClass = useMemo(
    () => (darkMode ? "bg-gray-800 text-white" : "bg-gray-100 text-black"),
    [darkMode],
  );

  const hijriText = useMemo(() => {
    return `${hijriDate.day} ${
      MONTHS[hijriDate.month - 1]
    } ${hijriDate.year} H`;
  }, [hijriDate]);

  const methodLabel = useMemo(() => {
    switch (method) {
      case "global":
        return "Global (Standar Umm al-Qura)";
      case "hisab":
        return "Hisab Wujudul Hilal";
      case "rukyat":
        return "Rukyat Hilal";
      default:
        return method;
    }
  }, [method]);

  return (
    <div
      className={`mt-4 p-4 rounded-lg shadow transition-colors ${themeClass}`}
    >
      <p className="text-lg font-medium">🗓️ Tanggal Hijriyah</p>
      <hr className="my-2 opacity-50" />

      <p className="text-lg font-semibold">
        {weton && <span>{weton}, </span>}
        {hijriText}
      </p>

      <p className="mt-1 text-sm opacity-70">
        Metode: <strong>{methodLabel}</strong>
      </p>
    </div>
  );
}
