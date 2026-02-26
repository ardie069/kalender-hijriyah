import { Moon, Sparkles, Sun } from "lucide-react";

interface LocalSliderProps {
  value: number;
  onChange: (val: number) => void;
  alt: number;
  elong: number;
}

export default function LocalSlider({
  value,
  onChange,
  alt,
}: LocalSliderProps) {
  const isAbove = alt > 0;

  const lagTime = Math.ceil(alt * 4);
  const maxRange = isAbove ? Math.max(15, lagTime) : 30;

  return (
    <div className="space-y-3">
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-2">
          <div
            className={`p-1.5 rounded-lg ${isAbove ? "bg-primary/20 text-primary" : "bg-red-500/20 text-red-500"}`}
          >
            {isAbove ? (
              <Sun className="w-3 h-3" />
            ) : (
              <Moon className="w-3 h-3" />
            )}
          </div>
          <div className="flex flex-col">
            <span className="text-[10px] font-black uppercase text-gray-400">
              T + {value} Menit
            </span>
            {isAbove && value === lagTime && (
              <span className="text-[8px] font-bold text-orange-500 uppercase leading-none">
                Moonset Point ⚓
              </span>
            )}
          </div>
        </div>

        {/* Best view biasanya 15-25 menit setelah sunset, tapi hanya jika belum moonset */}
        {value >= 15 && value < lagTime && (
          <span className="text-[9px] font-black text-primary flex items-center gap-1 uppercase animate-pulse">
            <Sparkles className="w-3 h-3" /> Best View
          </span>
        )}
      </div>

      <input
        type="range"
        min="0"
        max={maxRange}
        value={value}
        onChange={(e) => onChange(parseInt(e.target.value))}
        className="w-full h-2 rounded-lg appearance-none cursor-pointer accent-primary bg-gray-200 dark:bg-gray-800 transition-all"
      />

      <div className="flex justify-between text-[8px] font-bold text-gray-500 uppercase px-1">
        <span>Sunset</span>
        {isAbove && <span>Moonset ({lagTime}m)</span>}
      </div>
    </div>
  );
}
