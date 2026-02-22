"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useSyncExternalStore,
} from "react";

type ThemeContextType = {
  darkMode: boolean;
  toggleTheme: () => void;
  themeToggleText: string;
};

const ThemeContext = createContext<ThemeContextType | null>(null);

// External store untuk localStorage theme — SSR-safe via useSyncExternalStore
let listeners: Array<() => void> = [];

function subscribe(callback: () => void) {
  listeners = [...listeners, callback];
  return () => {
    listeners = listeners.filter((l) => l !== callback);
  };
}

function getSnapshot(): boolean {
  return localStorage.getItem("theme") === "dark";
}

function getServerSnapshot(): boolean {
  return false;
}

function setTheme(isDark: boolean) {
  localStorage.setItem("theme", isDark ? "dark" : "light");
  // Notify semua subscriber bahwa state berubah
  for (const listener of listeners) {
    listener();
  }
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const darkMode = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  const toggleTheme = () => setTheme(!darkMode);

  // Jika localStorage belum pernah diset, inisialisasi dengan "light"
  useEffect(() => {
    if (localStorage.getItem("theme") === null) {
      localStorage.setItem("theme", "light");
    }
  }, []);

  const themeToggleText = useMemo(
    () => (darkMode ? "🌙 Mode Gelap" : "☀️ Mode Terang"),
    [darkMode],
  );

  return (
    <ThemeContext.Provider value={{ darkMode, toggleTheme, themeToggleText }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error("useTheme must be used inside ThemeProvider");
  }
  return ctx;
}
