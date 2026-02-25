import { LineChart, ArrowUpToLine, Ruler, History } from "lucide-react";
import { MoonTelemetry } from "@/types/moon";

interface DetailedDataTableProps {
  telemetry: MoonTelemetry;
}

function TableRow({
  icon,
  label,
  value,
  threshold,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
  threshold: string;
}) {
  return (
    <tr className="hover:bg-white/5 transition-colors">
      <td className="px-4 py-3 sm:px-6 md:px-8 sm:py-4 md:py-5 font-bold flex items-center gap-3 sm:gap-4 text-gray-900 dark:text-gray-100">
        <div className="p-2 bg-gray-100 dark:bg-gray-800 rounded-xl">
          {icon}
        </div>
        {label}
      </td>
      <td className="px-4 py-3 sm:px-6 md:px-8 sm:py-4 md:py-5 font-mono text-primary font-black">{value}</td>
      <td className="px-4 py-3 sm:px-6 md:px-8 sm:py-4 md:py-5 text-gray-500 text-[9px] sm:text-[10px] font-black uppercase tracking-widest hidden sm:table-cell">
        {threshold}
      </td>
    </tr>
  );
}

export default function DetailedDataTable({
  telemetry,
}: DetailedDataTableProps) {
  return (
    <div className="bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl rounded-2xl border border-white/40 dark:border-white/5 overflow-hidden shadow-soft">
      <div className="px-4 py-4 sm:px-6 md:px-8 sm:py-5 md:py-6 border-b border-white/10 flex justify-between items-center bg-white/20">
        <h3 className="text-[10px] font-black text-gray-400 uppercase tracking-[0.4em]">
          Detailed Data
        </h3>
        <LineChart className="w-4 h-4 text-gray-400" />
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <tbody className="divide-y divide-white/5 text-sm">
            <TableRow
              icon={<ArrowUpToLine className="w-4 h-4 text-blue-500" />}
              label="Moon Altitude"
              value={`${telemetry.altitude}°`}
              threshold="Min. 3.0°"
            />
            <TableRow
              icon={<Ruler className="w-4 h-4 text-purple-500" />}
              label="Elongation"
              value={`${telemetry.elongation}°`}
              threshold="Min. 6.4°"
            />
            <TableRow
              icon={<History className="w-4 h-4 text-amber-500" />}
              label="Moon Age"
              value={`${telemetry.age_days} Days`}
              threshold="After New Moon"
            />
          </tbody>
        </table>
      </div>
    </div>
  );
}
