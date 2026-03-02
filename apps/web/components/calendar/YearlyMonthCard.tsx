import type { DateSystem } from "@/types/calendar";

interface MonthInfo {
  id: number;
  name: string;
  desc: string;
  isSpecial?: boolean;
}

interface YearlyMonthCardProps {
  month: MonthInfo;
  dateSystem: DateSystem;
  onClick: () => void;
}

export function YearlyMonthCard({
  month,
  dateSystem,
  onClick,
}: YearlyMonthCardProps) {
  return (
    <div
      onClick={onClick}
      className="group cursor-pointer bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-8 rounded-2xl border border-white/40 dark:border-white/5 hover:border-primary/40 transition-all duration-500 hover:scale-[1.03] shadow-soft"
    >
      <div className="flex justify-between items-start mb-6">
        <div>
          <p className="text-[8px] font-black opacity-30 uppercase tracking-[0.3em] mb-1">
            {dateSystem === "hijri"
              ? `Orbit Stage ${month.id}`
              : `Bulan ke-${month.id}`}
          </p>
          <h3 className="text-2xl font-black text-gray-900 dark:text-white group-hover:text-primary transition-colors">
            {month.name}
          </h3>
        </div>
        {month.isSpecial && <span className="text-2xl drop-shadow-sm">✨</span>}
      </div>
      <div className="grid grid-cols-7 gap-1.5 opacity-10 group-hover:opacity-100 transition-opacity duration-1000">
        {Array.from({ length: 30 }).map((_, i) => (
          <div
            key={i}
            className="aspect-square bg-gray-400/30 dark:bg-white/10 rounded-sm"
          ></div>
        ))}
      </div>
      <p className="mt-6 text-[9px] font-black opacity-30 uppercase tracking-[0.4em] text-center">
        {month.desc}
      </p>
    </div>
  );
}
