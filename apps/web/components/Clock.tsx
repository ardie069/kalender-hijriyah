"use client";

import { useEffect, useState } from "react";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

interface ClockProps {
  userTimezone: string;
}

export default function Clock({ userTimezone }: ClockProps) {
  const { darkMode } = useTheme();
  const mounted = useMounted();
  const [time, setTime] = useState({
    clock: "--:--:--",
    date: "Memuat waktu...",
    day: "",
  });

  useEffect(() => {
    const tick = () => {
      const now = new Date();

      // Dialektika Hari: Ahad lebih estetik buat kalender Hijriyah
      const dayName = now.toLocaleDateString("id-ID", { weekday: "long" });
      const day = dayName === "Minggu" ? "Ahad" : dayName;

      const dateStr = now.toLocaleDateString("id-ID", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });

      const timeStr = now.toLocaleTimeString("id-ID", {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: false,
      });

      setTime({ clock: timeStr, date: dateStr, day: day });
    };

    tick();
    const id = setInterval(tick, 1000);
    return () => clearInterval(id);
  }, []);

  const containerClass =
    mounted && darkMode
      ? "bg-gray-900/50 border-gray-800 text-white shadow-inner"
      : "bg-gray-50 border-gray-200 text-gray-900 shadow-sm";

  return (
    <div
      className={`p-6 rounded-4xl border transition-all duration-500 ${containerClass}`}
    >
      <div className="flex flex-col items-center justify-center space-y-1">
        {/* Visual Pulse: Jam yang dominan */}
        <p className="text-5xl font-black tracking-tighter tabular-nums font-mono mb-1">
          {mounted ? time.clock : "00:00:00"}
        </p>

        {/* Meta info yang rapi */}
        <div className="flex items-center gap-2">
          <span className="px-2 py-0.5 rounded-lg bg-emerald-500/10 text-emerald-500 text-[10px] font-bold uppercase tracking-widest">
            {mounted ? time.day : "---"}
          </span>
          <p className="text-sm font-medium opacity-60">
            {mounted ? time.date : "Memuat..."}
          </p>
        </div>

        {/* Zona Waktu dengan gaya 'pill' */}
        <div className="mt-4 pt-3 border-t border-current/5 w-full flex justify-center">
          <p className="text-[10px] font-bold opacity-40 uppercase tracking-widest flex items-center gap-2">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            {userTimezone}
          </p>
        </div>
      </div>
    </div>
  );
}
