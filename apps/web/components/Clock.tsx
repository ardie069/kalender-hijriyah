"use client";

import { useEffect, useState } from "react";
import { useTheme } from "@/context/theme-context";

interface ClockProps {
  userTimezone: string;
}

export default function Clock({ userTimezone }: ClockProps) {
  const { darkMode } = useTheme();
  const [mounted, setMounted] = useState(false);
  const [timeText, setTimeText] = useState("🕒 Memuat waktu…");

  useEffect(() => {
    requestAnimationFrame(() => setMounted(true));

    const tick = () => {
      const now = new Date();
      let day = now.toLocaleDateString("id-ID", { weekday: "long" });
      if (day === "Minggu") day = "Ahad";

      const date = now.toLocaleDateString("id-ID", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });

      const time = now.toLocaleTimeString("id-ID", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });

      setTimeText(`🕒 ${day}, ${date}, ${time}`);
    };

    tick();
    const id = setInterval(tick, 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <div
      className={`mb-3 p-4 rounded-lg border transition-colors ${
        mounted && darkMode
          ? "bg-gray-800 border-gray-700 text-white"
          : "bg-white border-gray-200 text-gray-800"
      }`}
    >
      <p className="text-xl font-semibold">{timeText}</p>
      <p className="text-sm opacity-70">🌍 Zona Waktu: {userTimezone}</p>
    </div>
  );
}
