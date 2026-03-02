import type { DateSystem, ViewMode } from "@/types/calendar";
import type { Method as HijriMethod } from "@/types/hijri";

interface CalendarHeaderProps {
  year: number;
  dateSystem: DateSystem;
  method: HijriMethod;
  viewMode: ViewMode;
  onYearChange: (newYear: number) => void;
  onViewModeChange: (newMode: ViewMode) => void;
}

export function CalendarHeader({
  year,
  dateSystem,
  method,
  viewMode,
  onYearChange,
  onViewModeChange,
}: CalendarHeaderProps) {
  return (
    <header className="flex flex-col lg:flex-row justify-between items-start lg:items-end gap-8 mb-10">
      <div className="animate-in fade-in slide-in-from-left-4 duration-700">
        <h1 className="text-5xl md:text-7xl font-black tracking-tighter text-gray-900 dark:text-white">
          {year}{" "}
          <span className="text-primary italic">
            {dateSystem === "hijri" ? "Hijriyah" : "Masehi"}
          </span>
        </h1>
        <p className="mt-3 text-[10px] md:text-xs font-black uppercase tracking-[0.4em] opacity-40">
          {dateSystem === "hijri"
            ? `${method.toUpperCase()} Engine • Kalender Hijriyah`
            : "Kalender Masehi • Standard Calendar"}
        </p>
      </div>

      {/* View & Year Switcher */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="join bg-card-light dark:bg-card-dark border border-gray-100 dark:border-white/5 rounded-2xl p-1 shadow-soft">
          <button
            onClick={() => onYearChange(year - 1)}
            className="btn btn-ghost btn-sm join-item"
          >
            ←
          </button>
          <span className="px-4 py-1 text-sm font-black">{year}</span>
          <button
            onClick={() => onYearChange(year + 1)}
            className="btn btn-ghost btn-sm join-item"
          >
            →
          </button>
        </div>
        <div className="flex bg-gray-200/50 dark:bg-white/5 p-1 rounded-2xl border border-gray-100 dark:border-white/5">
          <button
            onClick={() => onViewModeChange("monthly")}
            className={`px-5 py-2 rounded-xl text-[9px] font-black uppercase tracking-widest transition-all cursor-pointer ${viewMode === "monthly" ? "bg-white dark:bg-primary text-primary dark:text-black shadow-md" : "opacity-40"}`}
          >
            Month
          </button>
          <button
            onClick={() => onViewModeChange("yearly")}
            className={`px-5 py-2 rounded-xl text-[9px] font-black uppercase tracking-widest transition-all cursor-pointer ${viewMode === "yearly" ? "bg-white dark:bg-primary text-primary dark:text-black shadow-md" : "opacity-40"}`}
          >
            Year
          </button>
        </div>
      </div>
    </header>
  );
}
