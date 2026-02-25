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
