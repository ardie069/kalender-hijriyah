import { Moon, Sparkles, Sun } from "lucide-react";


interface LocalSliderProps {
  value: number;
  onChange: (val: number) => void;
  alt: number;
  elong: number;
}

export default function LocalSlider({ value, onChange, alt, elong }: LocalSliderProps) {
  const isAbove = alt > 0;
  const maxRange = Math.min(120, Math.max(45, Math.ceil(elong * 5)));

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
          <span className="text-[10px] font-black uppercase text-gray-400">
            T + {value} Menit
          </span>
        </div>
        {value >= 15 && value <= 25 && isAbove && (
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
        className={`w-full h-2 rounded-lg appearance-none cursor-pointer accent-primary bg-gray-200 dark:bg-gray-800`}
      />
    </div>
  );
}