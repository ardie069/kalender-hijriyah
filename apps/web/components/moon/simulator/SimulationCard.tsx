import { ReactNode } from "react";

interface SimulationCardProps {
  visual: ReactNode;
  title: string;
  badge?: ReactNode;
  description: ReactNode;
  stats: ReactNode;
  controls?: ReactNode;
  className?: string;
}

export default function SimulationCard({
  visual,
  title,
  badge,
  description,
  stats,
  controls,
  className,
}: SimulationCardProps) {
  return (
    <div
      className={`group flex flex-col lg:flex-row gap-6 p-6 rounded-3xl transition-all ${className}`}
    >
      <div className="relative w-full lg:w-2/3 aspect-video rounded-3xl overflow-hidden">
        {visual}
      </div>
      <div className="flex-1 space-y-4 py-2">
        <div className="flex justify-between items-center">
          <h3 className="text-xl font-black dark:text-white uppercase tracking-tight">
            {title}
          </h3>
          {badge}
        </div>
        <div className="text-sm text-gray-500 leading-relaxed italic">
          {description}
        </div>
        <div className="grid grid-cols-1 gap-3 pt-2">{stats}</div>
        <div className="mt-8 pt-6 border-t border-white/10">{controls}</div>
      </div>
    </div>
  );
}
