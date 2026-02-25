export default function MoonHeroSection() {
  return (
    <header className="flex flex-col md:flex-row items-start md:items-center gap-4 sm:gap-8 mb-8 sm:mb-12 md:mb-16 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="p-4 sm:p-6 bg-emerald-500/10 rounded-2xl shadow-glow border border-primary/20">
        <span className="text-4xl sm:text-6xl md:text-7xl block animate-pulse">🌙</span>
      </div>
      <div>
        <h1 className="text-3xl sm:text-5xl md:text-7xl font-black tracking-tighter text-gray-900 dark:text-white">
          Lunar <span className="text-primary italic">Analytics</span>
        </h1>
        <p className="mt-2 sm:mt-4 text-xs sm:text-sm md:text-lg font-medium text-gray-500 dark:text-gray-400 max-w-2xl leading-relaxed">
          Sinkronisasi gerak materi angkasa dengan algoritma{" "}
          <span className="font-bold text-primary">Skyfield API</span>. Data
          presisi untuk navigasi waktu Hijriyah.
        </p>
      </div>
    </header>
  );
}
