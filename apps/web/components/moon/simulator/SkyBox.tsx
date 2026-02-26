import { Sun } from "lucide-react";

interface SkyBoxProps {
  alt: number;
  offset: number;
  isOrange: boolean;
  label: string;
  icon: React.ReactNode;
}

export default function SkyBox({ alt, offset, isOrange, label, icon }: SkyBoxProps) {
  const skyOpacity = Math.min(0.9, offset / 50);
  const ySun = 15 + (0 - offset * 0.25) * 3;
  const yMoon = 15 + (alt - offset * 0.25) * 3;

  return (
    <div className="relative z-0 w-full h-full overflow-hidden bg-[#0f172a]">
      <div
        className={`absolute inset-0 bg-linear-to-b ${isOrange ? "from-[#0f172a] via-[#1e293b] to-[#c2410c]" : "from-[#1e293b] via-[#334155] to-[#f97316]"}`}
      />
      <div
        className="absolute inset-0 bg-[#020617] transition-opacity duration-700"
        style={{ opacity: skyOpacity }}
      />
      <div className="absolute bottom-0 w-full h-[15%] bg-[#112116] z-10 border-t border-emerald-900/50" />

      <div className="absolute top-3 left-3 sm:top-6 sm:left-6 z-30 bg-black/60 px-2 py-1 sm:px-4 sm:py-2 rounded-full text-[8px] sm:text-[10px] font-black text-white uppercase tracking-[0.15em] sm:tracking-[0.2em] flex items-center gap-1.5 sm:gap-2 border border-white/10">
        {icon} {label}
      </div>

      <div
        className="absolute left-1/2 -translate-x-1/2 z-0 transition-all duration-300"
        style={{ bottom: `${ySun}%` }}
      >
        <div className="w-20 h-20 sm:w-32 sm:h-32 bg-orange-600 rounded-full blur-2xl sm:blur-[60px] opacity-40" />
        <Sun className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 sm:w-6 sm:h-6 text-orange-200" />
      </div>

      <div
        className="absolute left-[45%] z-20 transition-all duration-300"
        style={{ bottom: `${yMoon}%` }}
      >
        <div
          className={`w-8 h-8 sm:w-12 sm:h-12 rounded-full border-r-[3px] border-b-[0.5px] border-white rotate-35 shadow-[0_0_20px_white] transition-opacity ${alt - offset * 0.25 < -1 ? "opacity-0" : "opacity-100"}`}
        />
      </div>
    </div>
  );
}