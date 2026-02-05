"use client";

import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useMemo,
} from "react";

type ThemeContextType = {
  darkMode: boolean;
  toggleTheme: () => void;
  themeToggleText: string;
};

const ThemeContext = createContext<ThemeContextType | null>(null);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [darkMode, setDarkMode] = useState(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("theme");
      return saved === "dark";
    }
    return false;
  });

  useEffect(() => {
    localStorage.setItem("theme", darkMode ? "dark" : "light");
  }, [darkMode]);

  const toggleTheme = () => setDarkMode((v) => !v);

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
