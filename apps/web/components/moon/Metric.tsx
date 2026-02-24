interface MetricProps {
  label: string;
  value: string;
  sub: string;
}

export default function Metric({ label, value, sub }: MetricProps) {
  return (
    <div className="flex flex-col gap-1">
      <p className="text-[9px] font-black opacity-30 uppercase tracking-[0.3em]">
        {label}
      </p>
      <p className="text-3xl font-black tabular-nums text-gray-900 dark:text-white tracking-tighter">
        {value}
      </p>
      <p className="text-[8px] font-bold text-primary uppercase tracking-widest">
        {sub}
      </p>
    </div>
  );
}
