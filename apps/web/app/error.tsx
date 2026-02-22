"use client";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <main className="min-h-[calc(100vh-64px)] flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center space-y-6">
        <div className="text-6xl">⚠️</div>
        <h2 className="text-2xl font-black tracking-tight">
          Terjadi Kesalahan
        </h2>
        <p className="text-sm opacity-60 leading-relaxed">
          {error.message || "Sistem mengalami gangguan. Silakan coba lagi."}
        </p>
        <button
          onClick={() => reset()}
          className="px-6 py-3 bg-emerald-500 text-white font-bold text-sm rounded-2xl hover:bg-emerald-600 transition-colors active:scale-95 cursor-pointer"
        >
          Coba Lagi
        </button>
      </div>
    </main>
  );
}
