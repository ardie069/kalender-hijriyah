import "./globals.css";
import { ThemeProvider } from "@/context/theme-context";
import Navbar from "@/components/Navbar";
import { Geist, Geist_Mono } from "next/font/google";

// Font modern buat vibe 2026
const geistSans = Geist({
  subsets: ["latin"],
  variable: "--font-sans",
});

const geistMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata = {
  title: "Kalender Hijriyah | Digital Lunar Analytics",
  description:
    "Aplikasi kalender Hijriyah modern berbasis astronomi presisi dan kearifan lokal.",
  icons: {
    icon: "/icon.png",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased font-sans transition-colors duration-500 bg-gray-50 dark:bg-gray-950 text-gray-900 dark:text-gray-100`}
      >
        <ThemeProvider>
          {/* Wrapper biar content gak mentok Navbar */}
          <div className="min-h-screen flex flex-col">
            <Navbar />
            <div className="grow">{children}</div>
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}
