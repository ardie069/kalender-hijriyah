import { HIJRI_MONTHS_INDONESIA_GRAMMAR as HIJRI_MONTHS } from "@/lib/constants";
import type { DateSystem } from "@/types/calendar";

const WEEK_DAYS = [
  "Ahad",
  "Senin",
  "Selasa",
  "Rabu",
  "Kamis",
  "Jum'at",
  "Sabtu",
];

interface MonthGridProps {
  monthId: number;
  year: number;
  system: DateSystem;
}

export function MonthGrid({ monthId, year, system }: MonthGridProps) {
  const currentMonth = HIJRI_MONTHS.find((m) => m.id === monthId);
  return (
    <div className="bg-white/60 dark:bg-card-dark/40 backdrop-blur-3xl rounded-3xl border border-white/40 dark:border-white/5 shadow-soft p-6 md:p-12 animate-in zoom-in-95 duration-700">
      <header className="flex justify-between items-end mb-12">
        <div>
          <h2 className="text-4xl sm:text-5xl font-black tracking-tighter text-gray-900 dark:text-white uppercase">
            {currentMonth?.name}
          </h2>
          <p className="text-[10px] font-black text-primary uppercase tracking-[0.5em] mt-2">
            Projection: Year {year} H
          </p>
        </div>
        <div className="hidden sm:flex gap-1.5 mb-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className={`h-1.5 w-1.5 rounded-full ${i === 0 ? "bg-primary animate-pulse" : "bg-primary/20"}`}
            />
          ))}
        </div>
      </header>
      <div className="grid grid-cols-7 gap-2 md:gap-4">
        {WEEK_DAYS.map((d, i) => (
          <div
            key={d}
            className={`text-center py-3 text-[10px] font-black uppercase tracking-[0.3em] opacity-30 ${i === 5 ? "text-primary opacity-100" : ""}`}
          >
            {d}
          </div>
        ))}
        {Array.from({ length: 30 }).map((_, i) => {
          const day = i + 1;
          const isJumat = i % 7 === 5;
          return (
            <div
              key={i}
              className={`aspect-square sm:aspect-auto sm:min-h-30 p-5 rounded-[2.2rem] border transition-all duration-300 flex flex-col items-start group cursor-pointer
              ${isJumat ? "bg-primary/5 border-primary/20" : "bg-white/20 dark:bg-white/5 border-white/20 dark:border-white/5 hover:border-primary/40"}
            `}
            >
              <span
                className={`text-2xl sm:text-4xl font-black tabular-nums transition-all ${isJumat ? "text-primary" : "text-gray-900 dark:text-white opacity-80"}`}
              >
                {day}
              </span>
              <span className="text-[8px] sm:text-[10px] font-bold uppercase tracking-tighter opacity-20 group-hover:opacity-60 mt-1">
                {(i % 28) + 1} Mar
              </span>
              {day === 1 && (
                <div className="mt-auto h-1 w-6 bg-indigo-500 rounded-full shadow-[0_0_8px_rgba(99,102,241,0.5)]"></div>
              )}
              {day === 14 && (
                <div className="mt-auto h-3 w-3 rounded-full border-2 border-gray-400 dark:border-gray-500 bg-white"></div>
              )}
              {day === 29 && (
                <div className="mt-auto h-3 w-3 rounded-full bg-gray-800 dark:bg-gray-200"></div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
