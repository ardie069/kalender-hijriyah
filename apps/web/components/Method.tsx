"use client";

import type { Method as MethodType } from "@/types/hijri";

interface MethodProps {
  value: MethodType;
  onChange: (method: MethodType) => void;
}

const METHODS: { id: MethodType; label: string; icon: string; desc: string }[] =
  [
    { id: "ughc", label: "KHGT", icon: "🌍", desc: "Global Standard" },
    {
      id: "umm_al_qura",
      label: "Umm al-Qura",
      icon: "🕋",
      desc: "Saudi / Civil",
    },
    {
      id: "local_hisab",
      label: "Hisab Lokal",
      icon: "🔢",
      desc: "Wujudul Hilal",
    },
    {
      id: "local_rukyat",
      label: "Rukyat Lokal",
      icon: "🔭",
      desc: "MABIMS Kriteria",
    },
  ];

export default function Method({ value, onChange }: MethodProps) {

  return (
    <div className="w-full space-y-4 animate-in fade-in slide-in-from-right-8 duration-1000">
      {/* Header Label: Scientific Tag */}
      <div className="flex items-center gap-3 ml-1">
        <div className="h-px w-4 bg-primary/40"></div>
        <label className="text-[9px] font-black uppercase tracking-[0.4em] text-primary opacity-80">
          Calculation Engine
        </label>
      </div>

      {/* Main Container: The Bento Grid */}
      <div className="grid grid-cols-2 gap-3 p-2 bg-white/40 dark:bg-card-dark/40 backdrop-blur-2xl rounded-[2rem] border border-white/40 dark:border-white/5 shadow-soft">
        {METHODS.map((m) => {
          const isActive = value === m.id;

          return (
            <button
              key={m.id}
              onClick={() => onChange(m.id)}
              className={`
                group relative flex flex-col items-start p-4 rounded-[1.4rem] transition-all duration-500 overflow-hidden cursor-pointer
                ${
                  isActive
                    ? "bg-primary text-white shadow-lg shadow-primary/20 scale-[1.02] z-10"
                    : "hover:bg-primary/5 text-gray-500 dark:text-gray-400"
                }
              `}
            >
              {/* Background Glow Effect untuk yang Aktif */}
              {isActive && (
                <div className="absolute inset-0 bg-linear-to-br from-white/20 to-transparent pointer-events-none" />
              )}

              {/* Status Dot: Logic Indicator */}
              <div className="absolute top-4 right-4">
                <div
                  className={`h-1.5 w-1.5 rounded-full transition-all duration-500 ${
                    isActive
                      ? "bg-white animate-pulse"
                      : "bg-gray-300 dark:bg-gray-700"
                  }`}
                />
              </div>

              <div className="flex items-center gap-3 mb-3">
                <span
                  className={`text-xl transition-transform duration-500 ${isActive ? "scale-110 rotate-12" : "group-hover:scale-110 opacity-60"}`}
                >
                  {m.icon}
                </span>
                <span
                  className={`text-xs font-black tracking-tight ${isActive ? "text-white" : "text-gray-900 dark:text-gray-100"}`}
                >
                  {m.label}
                </span>
              </div>

              <div className="flex flex-col items-start">
                <span
                  className={`text-[8px] font-mono font-bold uppercase tracking-widest leading-none ${isActive ? "text-white/80" : "opacity-40"}`}
                >
                  {m.desc}
                </span>
              </div>

              {/* Hover Scanline Effect */}
              {!isActive && (
                <div className="absolute inset-0 bg-primary/10 translate-y-full group-hover:translate-y-0 transition-transform duration-700 pointer-events-none opacity-20" />
              )}
            </button>
          );
        })}
      </div>

      {/* Footer Hint: System Status */}
      <div className="flex justify-between items-center px-2 mt-2">
        <p className="text-[8px] font-bold opacity-30 uppercase tracking-widest">
          Engine: {value.toUpperCase()}_REV_2026
        </p>
        <div className="flex gap-1">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className={`h-1 w-1 rounded-full ${i === 1 ? "bg-primary" : "bg-gray-200 dark:bg-gray-800"}`}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
