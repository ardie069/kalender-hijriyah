"use client";

export default function MoonSkeleton() {
  return (
    <main className="min-h-screen bg-background-light dark:bg-background-dark transition-colors duration-500 pb-24">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-10 animate-pulse">
        {/* 1. Hero Skeleton */}
        <div className="flex flex-col md:flex-row items-start md:items-center gap-8 mb-16">
          <div className="w-24 h-24 sm:w-32 sm:h-32 bg-gray-200 dark:bg-white/5 rounded-2xl" />
          <div className="space-y-4">
            <div className="h-12 w-64 bg-gray-200 dark:bg-white/5 rounded-2xl" />
            <div className="h-4 w-96 bg-gray-100 dark:bg-white/5 rounded-lg" />
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* 2. Main Telemetry Skeleton (8 Cols) */}
          <div className="lg:col-span-8 space-y-8">
            <div className="bg-gray-100 dark:bg-white/5 h-125 rounded-3xl" />
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="bg-gray-100 dark:bg-white/5 h-64 rounded-2xl" />
              <div className="bg-gray-100 dark:bg-white/5 h-64 rounded-2xl" />
            </div>
          </div>

          {/* 3. Sidebar Skeleton (4 Cols) */}
          <div className="lg:col-span-4 space-y-8">
            <div className="bg-gray-100 dark:bg-white/5 h-64 rounded-2xl" />
            <div className="bg-gray-100 dark:bg-white/5 h-80 rounded-2xl" />
            <div className="bg-gray-100 dark:bg-white/5 h-48 rounded-2xl" />
          </div>
        </div>

        {/* 4. Banner Skeleton */}
        <div className="mt-12 bg-gray-100 dark:bg-white/5 h-24 rounded-3xl" />
      </div>
    </main>
  );
}
