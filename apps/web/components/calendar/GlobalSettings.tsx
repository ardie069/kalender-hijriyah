import { METHODS } from "@/lib/constants";
import type { DateSystem } from "@/types/calendar";
import type { Method as HijriMethod } from "@/types/hijri";

interface GlobalSettingsProps {
  method: HijriMethod;
  onMethodChange: (newMethod: HijriMethod) => void;
  dateSystem: DateSystem;
  onDateSystemChange: (newDateSystem: DateSystem) => void;
}

export function GlobalSettings({
  method,
  onMethodChange,
  dateSystem,
  onDateSystemChange,
}: GlobalSettingsProps) {
  return (
    <section className="grid grid-cols-1 md:grid-cols-12 gap-4 mb-12">
      <div className="md:col-span-8 bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-2 rounded-3xl border border-white/40 dark:border-white/5 flex flex-wrap gap-2">
        {METHODS.map((m) => (
          <button
            key={m.id}
            onClick={() => onMethodChange(m.id)}
            className={`flex-1 min-w-25 flex items-center justify-center gap-2 py-3 rounded-2xl text-[9px] font-black uppercase tracking-widest transition-all cursor-pointer ${method === m.id ? "bg-primary text-white shadow-lg shadow-primary/20" : "opacity-40 hover:opacity-100"}`}
          >
            <span>{m.icon}</span> {m.label}
          </button>
        ))}
      </div>
      <div className="md:col-span-4 flex bg-gray-200/50 dark:bg-white/5 p-2 rounded-3xl border border-gray-100 dark:border-white/5">
        <button
          onClick={() => onDateSystemChange("hijri")}
          className={`flex-1 py-3 rounded-2xl text-[9px] font-black uppercase tracking-widest transition-all cursor-pointer ${dateSystem === "hijri" ? "bg-white dark:bg-primary text-primary dark:text-black shadow-md" : "opacity-40"}`}
        >
          Hijriyah
        </button>
        <button
          onClick={() => onDateSystemChange("gregorian")}
          className={`flex-1 py-3 rounded-2xl text-[9px] font-black uppercase tracking-widest transition-all cursor-pointer ${dateSystem === "gregorian" ? "bg-white dark:bg-primary text-primary dark:text-black shadow-md" : "opacity-40"}`}
        >
          Masehi
        </button>
      </div>
    </section>
  );
}
