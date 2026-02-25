import { ReactNode } from "react";

interface MiniStatProps {
  icon: ReactNode;
  label: string;
  value: string;
}

export default function MiniStat({ icon, label, value }: MiniStatProps) {
  return (
    <div className="flex items-center gap-3 bg-white/5 dark:bg-white/5 p-3 rounded-2xl border border-white/5">
      <div className="p-2 bg-gray-100 dark:bg-gray-800 rounded-lg text-gray-500">
        {icon}
      </div>
      <div>
        <p className="text-[9px] font-black uppercase tracking-widest text-gray-400">
          {label}
        </p>
        <p className="text-sm font-bold text-gray-900 dark:text-white">
          {value}
        </p>
      </div>
    </div>
  );
}
