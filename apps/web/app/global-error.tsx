"use client";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="id">
      <body className="bg-gray-950 text-white">
        <main className="min-h-screen flex items-center justify-center px-4">
          <div className="max-w-md w-full text-center space-y-6">
            <div className="text-6xl">💥</div>
            <h2 className="text-2xl font-black tracking-tight">
              Kesalahan Sistem
            </h2>
            <p className="text-sm opacity-60 leading-relaxed">
              {error.message ||
                "Terjadi kesalahan kritis. Silakan muat ulang halaman."}
            </p>
            <button
              onClick={() => reset()}
              className="px-6 py-3 bg-emerald-500 text-white font-bold text-sm rounded-2xl hover:bg-emerald-600 transition-colors active:scale-95 cursor-pointer"
            >
              Muat Ulang
            </button>
          </div>
        </main>
      </body>
    </html>
  );
}
