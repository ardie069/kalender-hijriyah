"use client";

import { useState } from "react";
import Link from "next/link";
import { useTheme } from "@/context/theme-context";
import { useMounted } from "@/hooks/use-mounted";

export default function Navbar() {
  const { darkMode, toggleTheme, themeToggleText } = useTheme();
  const mounted = useMounted();
  const [isOpen, setIsOpen] = useState(false);

  const navBg =
    mounted && darkMode
      ? "bg-gray-900/90 text-white"
      : "bg-white/90 text-black";

  return (
    <nav
      className={`sticky top-0 z-50 backdrop-blur-md shadow-sm transition-colors duration-300 ${navBg}`}
    >
      <div className="max-w-5xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Brand */}
          <Link
            href="/"
            className="font-bold text-lg flex items-center gap-2 active:scale-95 transition-transform"
          >
            <span className="text-2xl">🕌</span>
            <span className="inline">Kalender Hijriyah</span>
          </Link>

          {/* Desktop & Mobile Actions */}
          <div className="flex items-center gap-2">
            {/* Theme Toggle - Selalu muncul tapi kecil di mobile */}
            <button
              onClick={toggleTheme}
              className="p-2 rounded-full hover:bg-black/5 dark:hover:bg-white/5 transition-colors"
              aria-label="Toggle Theme"
            >
              {mounted ? (darkMode ? "🌙" : "☀️") : "..."}
            </button>

            {/* Hamburger Button - Hanya muncul di mobile */}
            <button
              onClick={() => setIsOpen(!isOpen)}
              className="md:hidden p-2 w-10 h-10 relative focus:outline-none"
              aria-label="Toggle Menu"
            >
              <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-6 h-5">
                <span
                  className={`absolute block h-0.5 w-6 bg-current transform transition-all duration-300 ease-in-out ${
                    isOpen ? "rotate-45 top-2.5" : "top-0"
                  }`}
                />
                <span
                  className={`absolute block h-0.5 w-6 bg-current transition-all duration-300 ease-in-out ${
                    isOpen ? "opacity-0" : "top-2.5"
                  }`}
                />
                <span
                  className={`absolute block h-0.5 w-6 bg-current transform transition-all duration-300 ease-in-out ${
                    isOpen ? "-rotate-45 top-2.5" : "top-5"
                  }`}
                />
              </div>
            </button>

            {/* Desktop Links */}
            <div className="hidden md:flex items-center gap-4 ml-4">
              <Link
                href="/calendar"
                className="text-sm font-medium hover:text-emerald-500"
              >
                Kalender
              </Link>
              <Link
                href="/moon-info"
                className="text-sm font-medium hover:text-emerald-500"
              >
                Info Bulan
              </Link>
            </div>
          </div>
        </div>

        {/* Mobile Dropdown - Dialektika antara Ada dan Tiada */}
        <div
          className={`md:hidden overflow-hidden transition-all duration-300 ease-in-out ${
            isOpen
              ? "max-h-60 opacity-100 border-t border-current/10"
              : "max-h-0 opacity-0"
          }`}
        >
          <div className="flex flex-col space-y-1 py-4">
            <Link
              href="/calendar"
              onClick={() => setIsOpen(false)}
              className="px-4 py-3 rounded-xl hover:bg-emerald-500/10 active:bg-emerald-500/20 transition-colors flex items-center gap-3"
            >
              <span>📅</span> Kalender
            </Link>
            <Link
              href="/moon-info"
              onClick={() => setIsOpen(false)}
              className="px-4 py-3 rounded-xl hover:bg-emerald-500/10 active:bg-emerald-500/20 transition-colors flex items-center gap-3"
            >
              <span>🌙</span> Informasi Bulan
            </Link>
            <div className="px-4 pt-2">
              <div className="p-3 rounded-xl bg-black/5 dark:bg-white/5 text-xs opacity-70">
                Mode: {mounted ? themeToggleText : "..."}
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
