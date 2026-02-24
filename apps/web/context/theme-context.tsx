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

// --- 1. EXTERNAL STORE LOGIC ---
let listeners: Array<() => void> = [];

function subscribe(callback: () => void) {
  listeners = [...listeners, callback];
  return () => {
    listeners = listeners.filter((l) => l !== callback);
  };
}

function getSnapshot(): boolean {
  if (typeof window === "undefined") return false;
  return localStorage.getItem("theme") === "dark";
}

function getServerSnapshot(): boolean {
  return false;
}

function setTheme(isDark: boolean) {
  if (typeof window === "undefined") return;

  const root = window.document.documentElement;
  if (isDark) {
    root.classList.add("dark");
    localStorage.setItem("theme", "dark");
  } else {
    root.classList.remove("dark");
    localStorage.setItem("theme", "light");
  }

  for (const listener of listeners) {
    listener();
  }
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const darkMode = useSyncExternalStore(
    subscribe,
    getSnapshot,
    getServerSnapshot,
  );

  const toggleTheme = () => setTheme(!darkMode);

  useEffect(() => {
    const savedTheme = localStorage.getItem("theme");
    const prefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)",
    ).matches;

    if (savedTheme === "dark" || (!savedTheme && prefersDark)) {
      setTheme(true);
    } else {
      setTheme(false);
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
