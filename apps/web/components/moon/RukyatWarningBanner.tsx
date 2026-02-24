"use client";

interface RukyatWarningBannerProps {
  isMet?: boolean;
}

export default function RukyatWarningBanner({
  isMet,
}: RukyatWarningBannerProps) {
  if (isMet === undefined) return null;

  return (
    <div
      className={`
        mt-12 border rounded-3xl p-6 flex flex-col sm:flex-row items-center gap-6 transition-all duration-700 animate-in slide-in-from-bottom-8
        ${
          isMet
            ? "bg-emerald-500/10 border-emerald-500/20 text-emerald-800 dark:text-emerald-400"
            : "bg-amber-500/10 border-amber-500/20 text-amber-800 dark:text-amber-400"
        }
      `}
    >
      <div
        className={`
        p-4 rounded-2xl shadow-lg text-2xl shrink-0
        ${isMet ? "bg-emerald-500 text-white shadow-emerald-500/20" : "bg-amber-500 text-white shadow-amber-500/20"}
      `}
      >
        {isMet ? "✅" : "⚠️"}
      </div>

      <div>
        <h4 className="text-sm font-black uppercase tracking-widest">
          {isMet ? "Potensi Hilal Teramati" : "Konfirmasi Rukyat Diperlukan"}
        </h4>
        <p className="text-sm mt-1 font-medium italic opacity-80 leading-relaxed">
          {isMet
            ? "Secara teknis, posisi hilal sudah berada di atas kriteria MABIMS. Kemungkinan besar bulan baru dimulai esok hari."
            : "Posisi hilal masih di bawah ambang batas visibilitas (3°). Kemungkinan besar bulan berjalan akan digenapkan (Istikmal) menjadi 30 hari."}
        </p>
      </div>
    </div>
  );
}
