"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function Navbar() {
  const { darkMode, toggleTheme } = useTheme();
  const mounted = useMounted();
  const pathname = usePathname();
  const [isOpen, setIsOpen] = useState(false);
  const isDark = mounted && darkMode;

  const isActive = (path: string) => pathname === path;

  return (
    <nav className="sticky top-0 z-50 bg-background-light/80 dark:bg-background-dark/80 backdrop-blur-xl border-b border-gray-100 dark:border-white/5 transition-all duration-500 shadow-soft">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-row items-center justify-between h-16 sm:h-20 relative">
          {/* --- BAGIAN KIRI: Identity --- */}
          <div className="flex-none">
            <Link
              href="/"
              className="flex items-center gap-3 active:scale-95 transition-all duration-300 group"
            >
              <div className="text-3xl sm:text-4xl group-hover:rotate-12 transition-transform duration-500 ease-out">
                {isDark ? "🌙" : "🕌"}
              </div>
              <div className="flex flex-col">
                <span className="font-black text-lg sm:text-xl tracking-tighter text-gray-900 dark:text-white leading-tight">
                  Kalender <span className="text-primary">Hijriyah</span>
                </span>
                <span className="text-[9px] font-black tracking-[0.4em] opacity-40 uppercase">
                  Digital Analytics
                </span>
              </div>
            </Link>
          </div>

          {/* --- BAGIAN KANAN: Control --- */}
          <div className="flex flex-row items-center gap-1 sm:gap-4">
            {/* Desktop Menu */}
            <div className="hidden md:flex flex-row items-center gap-2 mr-4 border-r border-gray-100 dark:border-white/10 pr-6">
              <NavLink
                href="/calendar"
                active={isActive("/calendar")}
                label="Kalender"
              />
              <NavLink
                href="/moon-info"
                active={isActive("/moon-info")}
                label="Info Bulan"
              />
            </div>

            {/* Theme Toggle */}
            <button
              onClick={toggleTheme}
              className={`btn btn-ghost btn-circle btn-sm sm:btn-md transition-all duration-500 hover:scale-110 active:scale-90 cursor-pointer
                ${isDark ? "hover:bg-yellow-500/10" : "hover:bg-indigo-500/10"}`}
            >
              {mounted ? (
                isDark ? (
                  <span className="text-xl">☀️</span>
                ) : (
                  <span className="text-xl">🌙</span>
                )
              ) : (
                <span className="loading loading-spinner loading-xs opacity-20"></span>
              )}
            </button>

            {/* --- FIX MOBILE DROPDOWN --- */}
            <div className="md:hidden relative">
              <button
                onClick={() => setIsOpen(!isOpen)}
                className="btn btn-ghost btn-circle btn-sm hover:bg-primary/10"
              >
                {/* Pakai SVG biar Ikon Pasti Muncul */}
                {isOpen ? (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={2.5}
                    stroke="currentColor"
                    className="w-6 h-6"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                ) : (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={2.5}
                    stroke="currentColor"
                    className="w-6 h-6"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                    />
                  </svg>
                )}
              </button>

              {/* Menu Content: Diposisikan absolut terhadap pemicu */}
              {isOpen && (
                <>
                  {/* Backdrop buat nutup menu pas klik di luar */}
                  <div
                    className="fixed inset-0 z-[-1]"
                    onClick={() => setIsOpen(false)}
                  />
                  <ul className="absolute right-0 mt-4 p-3 shadow-2xl bg-card-light dark:bg-card-dark rounded-3xl w-60 border border-gray-100 dark:border-white/10 animate-in fade-in slide-in-from-top-4 duration-300">
                    <li className="text-[10px] font-black uppercase tracking-[0.2em] opacity-40 mb-3 px-4 pt-2">
                      Navigasi
                    </li>
                    <li>
                      <Link
                        href="/calendar"
                        onClick={() => setIsOpen(false)}
                        className={`flex items-center gap-3 p-4 font-bold rounded-2xl transition-all ${isActive("/calendar") ? "bg-primary text-white" : "hover:bg-primary/10 hover:text-primary"}`}
                      >
                        📅 Kalender
                      </Link>
                    </li>
                    <li className="mt-1">
                      <Link
                        href="/moon-info"
                        onClick={() => setIsOpen(false)}
                        className={`flex items-center gap-3 p-4 font-bold rounded-2xl transition-all ${isActive("/moon-info") ? "bg-primary text-white" : "hover:bg-primary/10 hover:text-primary"}`}
                      >
                        🌙 Info Bulan
                      </Link>
                    </li>
                  </ul>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}

function NavLink({
  href,
  active,
  label,
}: {
  href: string;
  active: boolean;
  label: string;
}) {
  return (
    <Link
      href={href}
      className={`relative px-4 py-2 rounded-xl text-sm font-black transition-all duration-300
        ${active ? "text-primary bg-primary/5 shadow-soft" : "text-gray-500 hover:text-primary hover:bg-primary/5"}`}
    >
      {label}
      {active && (
        <span className="absolute -bottom-1 left-1/2 -translate-x-1/2 w-1 h-1 bg-primary rounded-full animate-pulse" />
      )}
    </Link>
  );
}
