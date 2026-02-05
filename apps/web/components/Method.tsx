"use client";

import type { Method } from "@/types/hijri";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

interface MethodProps {
  value: Method;
  onChange: (method: Method) => void;
}

export default function Method({ value, onChange }: MethodProps) {
  const { darkMode } = useTheme();
  const mounted = useMounted();

  return (
    <div className="mb-4 text-left">
      <label className="text-sm mb-1 block opacity-70">
        Metode Perhitungan:
      </label>

      <select
        value={value}
        onChange={(e) => onChange(e.target.value as Method)}
        className={`w-full px-3 py-2 rounded border transition-colors ${
          mounted && darkMode // Kunci utamanya di sini
            ? "bg-gray-800 text-white border-gray-600"
            : "bg-white text-black border-gray-300"
        }`}
      >
        <option value="global">Global (Umm al-Qura)</option>
        <option value="hisab">Hisab Astronomis</option>
        <option value="rukyat">Rukyat Hilal</option>
      </select>
    </div>
  );
}
