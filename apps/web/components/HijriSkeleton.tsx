"use client";

export default function HijriSkeleton() {
  return (
    <div className="space-y-6 animate-pulse px-1">
      {/* 1. Main Info Card Skeleton */}
      <div className="bg-card-light dark:bg-card-dark rounded-4xl p-8 border border-gray-100 dark:border-gray-800 shadow-card">
        <div className="flex flex-col items-center gap-4">
          <div className="h-4 w-24 bg-gray-200 dark:bg-gray-700 rounded-full" />
          <div className="h-12 w-64 bg-gray-300 dark:bg-gray-600 rounded-2xl" />
          <div className="h-6 w-32 bg-gray-200 dark:bg-gray-700 rounded-full" />
        </div>
      </div>

      {/* 2. Reasoning Box Skeleton */}
      <div className="p-6 rounded-4xl bg-black/5 dark:bg-white/5 border border-gray-100 dark:border-gray-800 space-y-4">
        <div className="flex justify-between items-center">
          <div className="h-3 w-32 bg-gray-200 dark:bg-gray-700 rounded-full" />
          <div className="h-4 w-16 bg-gray-200 dark:bg-gray-700 rounded-full" />
        </div>

        <div className="space-y-3">
          <div className="h-3 w-full bg-gray-200 dark:bg-gray-700 rounded-full opacity-60" />
          <div className="h-3 w-[90%] bg-gray-200 dark:bg-gray-700 rounded-full opacity-40" />
          <div className="h-3 w-[80%] bg-gray-200 dark:bg-gray-700 rounded-full opacity-20" />
        </div>

        <div className="pt-4 border-t border-gray-100 dark:border-gray-800 grid grid-cols-2 gap-6">
          <div className="space-y-2">
            <div className="h-2 w-16 bg-gray-200 dark:bg-gray-700 rounded-full" />
            <div className="h-6 w-20 bg-gray-300 dark:bg-gray-600 rounded-lg" />
          </div>
          <div className="space-y-2">
            <div className="h-2 w-16 bg-gray-200 dark:bg-gray-700 rounded-full" />
            <div className="h-6 w-20 bg-gray-300 dark:bg-gray-600 rounded-lg" />
          </div>
        </div>
      </div>

      {/* 3. Prediction Card Skeleton */}
      <div className="p-6 rounded-4xl bg-emerald-500/5 dark:bg-emerald-500/10 border border-emerald-500/10 space-y-4">
        <div className="flex justify-between">
          <div className="h-3 w-28 bg-emerald-200/50 dark:bg-emerald-800/50 rounded-full" />
          <div className="h-4 w-20 bg-emerald-200/50 dark:bg-emerald-800/50 rounded-full" />
        </div>
        <div className="h-10 w-48 bg-emerald-500/20 rounded-xl" />
        <div className="h-4 w-32 bg-emerald-200/50 dark:bg-emerald-800/50 rounded-full" />
      </div>
    </div>
  );
}
