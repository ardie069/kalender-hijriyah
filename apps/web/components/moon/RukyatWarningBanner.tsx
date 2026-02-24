export default function RukyatWarningBanner() {
  return (
    <div className="mt-12 bg-amber-500/10 border border-amber-500/20 rounded-3xl p-6 flex flex-col sm:flex-row items-center gap-6 animate-in slide-in-from-bottom-8 duration-1000">
      <div className="p-4 bg-amber-500 text-white rounded-2xl shadow-lg shadow-amber-500/20 text-2xl">
        ⚠️
      </div>
      <div>
        <h4 className="text-sm font-black text-amber-700 dark:text-amber-400 uppercase tracking-widest">
          Konfirmasi Rukyat Diperlukan
        </h4>
        <p className="text-sm text-amber-800 dark:text-amber-500/70 mt-1 font-medium italic">
          Meskipun data hisab menunjukkan posisi hilal di atas ufuk,
          penetapan awal bulan resmi menunggu sidang itsbat Kementerian
          Agama.
        </p>
      </div>
    </div>
  );
}
