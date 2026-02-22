import Link from "next/link";

export default function NotFound() {
  return (
    <main className="min-h-[calc(100vh-64px)] flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center space-y-6">
        <div className="text-6xl">🔍</div>
        <h2 className="text-2xl font-black tracking-tight">
          Halaman Tidak Ditemukan
        </h2>
        <p className="text-sm opacity-60 leading-relaxed">
          Halaman yang Anda cari tidak ada atau sudah dipindahkan.
        </p>
        <Link
          href="/"
          className="inline-block px-6 py-3 bg-emerald-500 text-white font-bold text-sm rounded-2xl hover:bg-emerald-600 transition-colors active:scale-95"
        >
          Kembali ke Beranda
        </Link>
      </div>
    </main>
  );
}
