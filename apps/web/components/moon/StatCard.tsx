interface StatCardProps {
  icon: string;
  label: string;
  value: string;
  desc: string;
  footer?: string;
  isSuccess?: boolean;
  progress?: number;
}

export default function StatCard({
  label,
  value,
  desc,
  footer,
  isSuccess = false,
  progress,
}: StatCardProps) {
  return (
    <div className="bg-white/40 dark:bg-card-dark/40 backdrop-blur-xl p-8 rounded-4xl border border-white/40 dark:border-white/5 flex flex-col h-full shadow-soft transition-all hover:scale-[1.02]">
      <div className="flex items-center gap-4 mb-6">
        <div className="p-3 bg-primary/10 rounded-2xl text-primary">⚡</div>
        <h3 className="font-black text-lg tracking-tight uppercase">{label}</h3>
      </div>
      <div className="grow">
        <p className="text-4xl font-black tabular-nums tracking-tighter">
          {value}
        </p>
        <p className="text-xs font-medium text-gray-500 dark:text-gray-400 mt-2 leading-relaxed italic">
          {desc}
        </p>
      </div>
      <div className="mt-8">
        {footer && (
          <div
            className={`text-[9px] font-black uppercase tracking-widest py-2 px-4 rounded-xl inline-block border ${isSuccess ? "bg-primary/10 text-primary border-primary/20" : "bg-gray-100 dark:bg-white/5 opacity-50"}`}
          >
            {footer}
          </div>
        )}
        {progress && (
          <div className="w-full mt-5 h-1.5 bg-gray-100 dark:bg-white/5 rounded-full overflow-hidden">
            <div
              className="h-full bg-primary"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
        )}
      </div>
    </div>
  );
}
