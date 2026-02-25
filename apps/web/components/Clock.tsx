"use client";

import { useEffect, useState } from "react";
import { useMounted } from "@/hooks/use-mounted";

interface ClockProps {
  userTimezone: string;
}

export default function Clock({ userTimezone }: ClockProps) {
  const mounted = useMounted();
  const [time, setTime] = useState({
    clock: "00.00.00",
    date: "Memuat waktu...",
    day: "---",
  });

  useEffect(() => {
    const tick = () => {
      const now = new Date();
      const dayName = now.toLocaleDateString("id-ID", { weekday: "long" });
      const day = dayName === "Minggu" ? "Ahad" : dayName;

      const dateStr = now.toLocaleDateString("id-ID", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });

      const timeStr = now
        .toLocaleTimeString("id-ID", {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hour12: false,
        })
        .replace(/:/g, ".");

      setTime({ clock: timeStr, date: dateStr, day: day });
    };

    tick();
    const id = setInterval(tick, 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="group relative animate-in fade-in slide-in-from-left-10 duration-1000">
      {/* Glow Aura - Dialektika Cahaya & Waktu */}
      <div className="absolute -inset-1 bg-linear-to-br from-primary/20 via-transparent to-emerald-500/10 rounded-3xl blur-2xl opacity-30 group-hover:opacity-60 transition-opacity duration-1000"></div>

      <div
        className={`
          relative overflow-hidden p-10 sm:p-14 rounded-[2.8rem] border transition-all duration-700
          flex flex-col items-center justify-center text-center backdrop-blur-3xl
          ${
            mounted
              ? "bg-white/60 dark:bg-card-dark/40 border-white/40 dark:border-white/5 shadow-card hover:border-primary/20"
              : "bg-gray-100 animate-pulse border-transparent"
          }
        `}
      >
        {/* Dekorasi Ornamen: Grid Geometris Tipis */}
        <div className="absolute inset-0 opacity-[0.03] dark:opacity-[0.07] pointer-events-none">
          <svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <pattern
                id="grid"
                width="40"
                height="40"
                patternUnits="userSpaceOnUse"
              >
                <path
                  d="M 40 0 L 0 0 0 40"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="0.5"
                />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#grid)" />
          </svg>
        </div>

        <div className="relative z-10 flex flex-col items-center">
          {/* Header: Label Teknis */}
          <div className="mb-6 inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/5 border border-primary/10">
            <span className="w-1 h-1 rounded-full bg-primary animate-pulse"></span>
            <p className="text-[9px] font-black uppercase tracking-[0.4em] text-primary">
              Real-time Chronos
            </p>
          </div>

          {/* Jam: Angka Dominan */}
          <div className="mb-4 relative">
            <h2 className="text-7xl sm:text-8xl font-black tabular-nums tracking-tighter transition-all duration-500 text-gray-900 dark:text-white dark:drop-shadow-[0_0_25px_rgba(255,255,255,0.1)]">
              {mounted ? time.clock : "00.00.00"}
            </h2>
          </div>

          {/* Badge Hari & Tanggal */}
          <div className="flex flex-col sm:flex-row items-center gap-3 mb-10">
            <span className="bg-primary text-white dark:text-black px-4 py-1 rounded-xl text-[10px] font-black uppercase tracking-[0.2em] shadow-lg shadow-primary/20">
              {mounted ? time.day : "---"}
            </span>
            <span className="text-xl font-bold opacity-50 text-gray-700 dark:text-gray-300">
              {mounted ? time.date : "Memuat..."}
            </span>
          </div>

          {/* Separator Dialektis */}
          <div className="w-24 border-t-2 border-primary/20 rounded-full mb-8" />

          {/* Footer: Sistem Info */}
          <div className="flex flex-col items-center gap-3">
            <div className="flex items-center gap-3 px-5 py-2 rounded-2xl bg-gray-50/50 dark:bg-white/5 border border-gray-100 dark:border-white/5">
              <div className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
              </div>
              <p className="text-[10px] font-black text-gray-500 dark:text-gray-400 uppercase tracking-[0.3em]">
                {userTimezone}
              </p>
            </div>
            <p className="text-[8px] font-bold opacity-30 uppercase tracking-widest leading-relaxed">
              Syncing with Global Atomic Standards
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
