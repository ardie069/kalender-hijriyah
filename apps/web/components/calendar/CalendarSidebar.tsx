import {
  HIJRI_MONTHS_INDONESIA_GRAMMAR as HIJRI_MONTHS,
  GREGORIAN_MONTHS,
} from "@/lib/constants";
import type { ViewMode, DateSystem } from "@/types/calendar";
import { LegendItem } from "./LegendItem";

interface CalendarSidebarProps {
  viewMode: ViewMode;
  selectedMonth: number;
  onMonthSelect: (monthId: number) => void;
  dateSystem: DateSystem;
  onMonthNav: (direction: "prev" | "next") => void;
}

export function CalendarSidebar({
  viewMode,
  selectedMonth,
  onMonthSelect,
  dateSystem,
  onMonthNav,
}: CalendarSidebarProps) {
  const monthList =
    dateSystem === "hijri"
      ? HIJRI_MONTHS.map((m) => ({
          id: m.id,
          name: m.name,
          isSpecial: "isSpecial" in m ? Boolean(m.isSpecial) : false,
        }))
      : GREGORIAN_MONTHS.map((m) => ({
          id: m.id,
          name: m.name,
          isSpecial: false,
        }));

  return (
    <aside className="lg:col-span-3 flex flex-col gap-8">
      {viewMode === "monthly" && (
        <>
          {/* Navigasi Prev / Next */}
          <div className="flex items-center justify-between px-2">
            <button
              onClick={() => onMonthNav("prev")}
              className="btn btn-ghost btn-sm rounded-xl text-xs font-black opacity-60 hover:opacity-100"
            >
              ← Prev
            </button>
            <button
              onClick={() => onMonthNav("next")}
              className="btn btn-ghost btn-sm rounded-xl text-xs font-black opacity-60 hover:opacity-100"
            >
              Next →
            </button>
          </div>

          <div className="space-y-2 max-h-[50vh] overflow-y-auto pr-2 scrollbar-hide">
            <p className="text-[10px] font-black uppercase tracking-[0.3em] opacity-30 mb-4 px-4">
              {dateSystem === "hijri" ? "Bulan Hijriyah" : "Bulan Masehi"}
            </p>
            {monthList.map((m) => (
              <button
                key={m.id}
                onClick={() => onMonthSelect(m.id)}
                className={`w-full flex items-center justify-between p-4 rounded-3xl transition-all duration-300 border cursor-pointer ${selectedMonth === m.id ? "bg-primary border-primary text-white shadow-xl shadow-primary/20 scale-[1.02]" : "bg-white/40 dark:bg-card-dark/40 border-gray-100 dark:border-white/5 opacity-60 hover:opacity-100"}`}
              >
                <span className="text-xs font-black tracking-tight">
                  {m.name}
                </span>
                <span className="text-xs">{m.isSpecial ? "✨" : ""}</span>
              </button>
            ))}
          </div>
        </>
      )}

      {/* BLOCK KETERANGAN (The Legend) */}
      <div className="bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-8 rounded-2xl border border-white/40 dark:border-white/5 shadow-soft">
        <h2 className="text-[10px] font-black tracking-[0.3em] text-primary uppercase mb-6">
          Keterangan
        </h2>
        <ul className="space-y-5">
          {dateSystem === "hijri" ? (
            <>
              <LegendItem
                color="bg-indigo-500"
                label="Awal Bulan"
                sub="New Crescent"
              />
              <LegendItem
                color="bg-white border-2 border-gray-400 dark:border-gray-500"
                label="Bulan Purnama"
                sub="Full Moon"
              />
              <LegendItem
                color="bg-gray-800 dark:bg-gray-200"
                label="Akhir Bulan"
                sub="Old Crescent"
              />
              <LegendItem
                color="bg-primary/10 border-2 border-primary"
                label="Hari Ini"
                sub="Today"
              />
            </>
          ) : (
            <LegendItem
              color="bg-primary/30 border border-primary"
              label="Hari Ini"
              sub="Today"
            />
          )}
          <LegendItem
            color="bg-primary/20 border border-primary"
            label="Hari Jumat"
            sub="Special Day"
          />
        </ul>
      </div>
    </aside>
  );
}
