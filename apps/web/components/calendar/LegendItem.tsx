export function LegendItem({
  color,
  label,
  sub,
}: {
  color: string;
  label: string;
  sub: string;
}) {
  return (
    <li className="flex items-center gap-4 group">
      <div
        className={`w-3.5 h-3.5 rounded-full shrink-0 shadow-sm transition-transform group-hover:scale-125 ${color}`}
      ></div>
      <div>
        <p className="text-xs font-black text-gray-900 dark:text-white leading-none">
          {label}
        </p>
        <p className="text-[8px] font-bold uppercase tracking-widest opacity-30 mt-1">
          {sub}
        </p>
      </div>
    </li>
  );
}
