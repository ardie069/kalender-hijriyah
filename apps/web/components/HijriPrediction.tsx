"use client";

import { useMemo } from "react";
import { useTheme } from "@/context/theme-context";
import type { HijriDate } from "@/types/hijri";

interface HijriVisibility {
  moon_altitude: number;
  elongation: number;
  moon_age: number;
  is_visible: boolean;
}

interface HijriPredictionData {
  today?: HijriDate;
  estimated_end_of_month?: HijriDate | null;
  message?: string;
  visibility?: HijriVisibility;
}

interface HijriPredictionProps {
  prediction: HijriPredictionData | null;
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

function formatHijri(date: HijriDate): string {
  return `${date.day} ${MONTHS[date.month - 1]} ${date.year} H`;
}

export default function HijriPrediction({ prediction }: HijriPredictionProps) {
  const { darkMode } = useTheme();

  const themeClass = useMemo(
    () => (darkMode ? "bg-gray-800 text-white" : "bg-gray-100 text-black"),
    [darkMode],
  );

  if (!prediction) return null;

  return (
    <div
      className={`mt-4 p-4 rounded-lg shadow transition-colors ${themeClass}`}
    >
      <h3 className="text-lg font-semibold mb-2">🌙 Informasi Akhir Bulan</h3>

      {prediction.message && (
        <p className="text-sm opacity-80">{prediction.message}</p>
      )}

      {prediction.estimated_end_of_month && (
        <div className="mt-2">
          <p>
            <strong>Perkiraan akhir bulan:</strong>{" "}
            {formatHijri(prediction.estimated_end_of_month)}
          </p>
        </div>
      )}

      {prediction.visibility && (
        <div className="mt-3 text-sm">
          <p className="font-medium">Validasi Imkanur Rukyat:</p>
          <ul className="list-disc list-inside">
            <li>
              Ketinggian bulan: {prediction.visibility.moon_altitude.toFixed(2)}
              °
            </li>
            <li>Elongasi: {prediction.visibility.elongation.toFixed(2)}°</li>
            <li>Usia bulan: {prediction.visibility.moon_age.toFixed(2)} jam</li>
            <li>
              Status:{" "}
              <span
                className={
                  prediction.visibility.is_visible
                    ? "text-green-500 font-semibold"
                    : "text-red-500 font-semibold"
                }
              >
                {prediction.visibility.is_visible
                  ? "Memenuhi syarat"
                  : "Tidak memenuhi syarat"}
              </span>
            </li>
          </ul>
        </div>
      )}
    </div>
  );
}
