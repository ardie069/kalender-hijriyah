"use client";

import { useState } from "react";
import Clock from "@/components/Clock";
import Method from "@/components/Method";
import HijriDate from "@/components/HijriDate";
import { useHijri } from "@/hooks/useHijri";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";
import type { Method as HijriMethod } from "@/types/hijri";

export default function HomeClient() {
  const { darkMode } = useTheme();
  const mounted = useMounted();
  const [selectedMethod, setSelectedMethod] =
    useState<HijriMethod>("umm_al_qura");

  const userTimezone = mounted
    ? Intl.DateTimeFormat().resolvedOptions().timeZone
    : "UTC";

  const {
    hijriDate,
    explanation,
    endMonthInfo,
    weton,
    loading,
    error,
    generatedAt,
  } = useHijri(selectedMethod, userTimezone);

  const methodLabels: Record<HijriMethod, string> = {
    umm_al_qura: "Umm Al-Qura",
    local_hisab: "Hisab Lokal",
    local_rukyat: "Rukyat Lokal",
    ughc: "KHGT",
  };

  if (!mounted) {
    return <div className="min-h-screen bg-background-light" />;
  }

  return (
    <div className="flex flex-col gap-8 lg:gap-12 animate-in fade-in duration-1000">
      {/* 1. Header Section */}
      <header className="space-y-4">
        <div className="flex items-center gap-4">
          <span className="text-5xl sm:text-6xl drop-shadow-md">
            {darkMode ? "🌙" : "🕌"}
          </span>
          <h1 className="text-4xl sm:text-5xl font-black tracking-tighter text-gray-900 dark:text-white transition-colors duration-500">
            Kalender Hijriyah
          </h1>
        </div>
        <p className="text-lg sm:text-xl text-gray-500 dark:text-gray-400 leading-relaxed max-w-2xl font-medium">
          Integrasi sains astronomi dan kearifan lokal dalam satu genggaman
          digital.
        </p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start">
        {/* KOLOM KIRI: Data Utama */}
        <div className="lg:col-span-7 flex flex-col gap-8">
          <section className="animate-in slide-in-from-left-8 duration-700">
            <Clock userTimezone={userTimezone} />
          </section>

          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-800/30 rounded-3xl p-6 flex items-start gap-4">
              {/* Konten Error */}
            </div>
          )}

          <section className="animate-in slide-in-from-bottom-8 duration-1000 delay-200">
            <HijriDate
              hijriDate={hijriDate}
              explanation={explanation}
              endMonthInfo={endMonthInfo}
              weton={weton}
              loading={loading}
              error={error}
              method={selectedMethod}
              generatedAt={generatedAt ?? undefined}
            />
          </section>
        </div>

        {/* KOLOM KANAN: Sidebar Kontrol */}
        <aside className="lg:col-span-5 flex flex-col gap-6 sticky top-24">
          <div className="rounded-4xl p-6 sm:p-8 shadow-card border transition-all duration-500 bg-card-light dark:bg-card-dark dark:text-white border-gray-100 dark:border-gray-800 shadow-gray-200/50 dark:shadow-emerald-500/5">
            <div className="mb-8">
              <h3 className="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-[0.2em] mb-6">
                Metode Perhitungan
              </h3>
              <Method value={selectedMethod} onChange={setSelectedMethod} />
            </div>

            <div className="w-full border-t border-gray-100 dark:border-gray-800 mb-8 border-dashed" />

            <div className="space-y-4">
              <h3 className="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-[0.2em] mb-4">
                Sistem Informasi
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <InfoBlock label="Zona Waktu" value={userTimezone} />
                <InfoBlock
                  label="Standardisasi"
                  value={methodLabels[selectedMethod]}
                />
              </div>
            </div>
          </div>

          {/* Quote Section ala Madilog */}
          <div className="rounded-4xl p-8 border transition-all duration-500 group bg-[#ecfdf5] dark:bg-primary/5 border-primary/10 dark:border-primary/20">
            <p className="italic font-semibold leading-relaxed transition-opacity text-primary-dark dark:text-primary opacity-90 dark:opacity-80 group-hover:opacity-100">
              <q>
                Waktu bukan sekadar angka, melainkan gerak materi di ruang
                angkasa yang kita tangkap melalui logika.
              </q>
            </p>
          </div>
        </aside>
      </div>
    </div>
  );
}

function InfoBlock({ label, value }: { label: string; value: string }) {
  return (
    <div className="border border-gray-100 dark:border-gray-800 rounded-2xl p-4 bg-gray-50/50 dark:bg-white/5 transition-colors duration-500">
      <span className="block text-[9px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-wider mb-1">
        {label}
      </span>
      <span className="font-bold text-xs text-gray-900 dark:text-white block truncate">
        {value}
      </span>
    </div>
  );
}
