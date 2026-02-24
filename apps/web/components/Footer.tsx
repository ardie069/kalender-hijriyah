"use client";

export default function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="mt-auto py-12 border-t border-gray-100 dark:border-gray-800 bg-white/30 dark:bg-background-dark/30 backdrop-blur-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row justify-between items-center gap-8">
          {/* Section 1: Branding & Philosophy */}
          <div className="flex flex-col items-center md:items-start gap-2">
            <div className="flex items-center gap-2">
              <span className="text-xl">🕌</span>
              <span className="text-xs font-black uppercase tracking-[0.3em] text-gray-900 dark:text-white">
                Kalender Hijriyah Digital
              </span>
            </div>
            <p className="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest">
              © {currentYear} • Build with Logic
            </p>
          </div>

          {/* Section 2: Technical Specs (The Engine) */}
          <div className="flex flex-wrap justify-center gap-4 sm:gap-6">
            <TechBadge label="Framework" value="Next.js 15" />
            <TechBadge label="Backend" value="FastAPI" />
            <TechBadge label="Astronomy" value="Skyfield" />
          </div>

          {/* Section 3: System Status Indicator */}
          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2 px-4 py-2 rounded-2xl bg-primary/5 border border-primary/10">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
              </span>
              <span className="text-[10px] font-black text-primary uppercase tracking-tighter">
                System Online
              </span>
            </div>
          </div>
        </div>

        {/* Bottom Line: Quote Minimalis */}
        <div className="mt-12 pt-8 border-t border-gray-100 dark:border-gray-800 text-center">
          <p className="text-[9px] font-medium text-gray-400 dark:text-gray-600 uppercase tracking-[0.4em] max-w-2xl mx-auto leading-relaxed">
            <q>
              Sains adalah obor yang membakar kegelapan ghaib dengan cahaya
              logika materialistik.
            </q>
          </p>
        </div>
      </div>
    </footer>
  );
}

// Sub-komponen biar rapi
function TechBadge({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex flex-col items-center md:items-start">
      <span className="text-[8px] font-black text-gray-400 dark:text-gray-600 uppercase tracking-tighter mb-0.5">
        {label}
      </span>
      <span className="text-[10px] font-bold text-gray-700 dark:text-gray-300 hover:text-primary transition-colors cursor-default">
        {value}
      </span>
    </div>
  );
}
