"use client";

import type { Method } from "@/types/hijri";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

interface MethodProps {
  value: Method;
  onChange: (method: Method) => void;
}

const METHODS: { id: Method; label: string; icon: string; desc: string }[] = [
  { id: "global", label: "Global", icon: "🌍", desc: "Umm al-Qura" },
  { id: "hisab", label: "Hisab", icon: "🔢", desc: "Wujudul Hilal" },
  { id: "rukyat", label: "Rukyat", icon: "🔭", desc: "MABIMS" },
];

export default function Method({ value, onChange }: MethodProps) {
  const { darkMode } = useTheme();
  const mounted = useMounted();

  const containerBg = mounted && darkMode ? "bg-gray-950/50" : "bg-gray-100";

  return (
    <div className="mb-6">
      <label className="text-[10px] font-bold uppercase tracking-widest opacity-40 mb-3 block ml-2">
        Metode Perhitungan
      </label>

      <div
        className={`grid grid-cols-3 gap-1.5 p-1.5 rounded-2xl transition-colors ${containerBg}`}
      >
        {METHODS.map((m) => {
          const isActive = value === m.id;

          return (
            <button
              key={m.id}
              onClick={() => onChange(m.id)}
              className={`
                flex flex-col items-center justify-center py-3 px-1 rounded-xl transition-all duration-300 cursor-pointer
                ${
                  isActive
                    ? "bg-emerald-500 text-white shadow-md scale-[1.02]"
                    : "hover:bg-current/5 opacity-60 hover:opacity-100"
                }
              `}
            >
              <span className="text-sm mb-0.5">{m.icon}</span>
              <span className="text-[11px] font-bold leading-none">
                {m.label}
              </span>
              <span
                className={`text-[8px] mt-1 opacity-70 ${isActive ? "text-white" : ""}`}
              >
                {m.desc}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
