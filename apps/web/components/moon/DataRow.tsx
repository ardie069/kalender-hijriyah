interface DataRowProps {
  label: string;
  value: string;
  sub: string;
}

export default function DataRow({ label, value, sub }: DataRowProps) {
  return (
    <li className="flex justify-between items-center group">
      <div>
        <p className="text-xs font-black text-gray-900 dark:text-white uppercase tracking-tight">
          {label}
        </p>
        <p className="text-[8px] font-bold opacity-30 uppercase tracking-widest">
          {sub}
        </p>
      </div>
      <div className="text-right">
        <p className="text-sm font-black text-primary tabular-nums tracking-tighter">
          {value}
        </p>
      </div>
    </li>
  );
}
