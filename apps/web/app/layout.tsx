import type { Metadata } from "next";
import "./globals.css";
import { ThemeProvider } from "@/context/theme-context";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Geist, Geist_Mono } from "next/font/google";

const geistSans = Geist({
  subsets: ["latin"],
  variable: "--font-sans",
});

const geistMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "Kalender Hijriyah | Digital Lunar Analytics",
  description:
    "Integrasi sains astronomi dan kearifan lokal dalam satu genggaman digital.",
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
    <html
      lang="id"
      className={`${geistSans.variable} ${geistMono.variable}`}
      suppressHydrationWarning
    >
      <head>
        {/* Anti-Flicker Script: Dialektika Materi Sebelum Paint */}
        <script
          dangerouslySetInnerHTML={{
            __html: `
              (function() {
                try {
                  var theme = localStorage.getItem('theme');
                  var supportDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches === true;
                  if (theme === 'dark' || (!theme && supportDarkMode)) {
                    document.documentElement.classList.add('dark');
                  } else {
                    document.documentElement.classList.remove('dark');
                  }
                } catch (e) {}
              })();
            `,
          }}
        />
      </head>
      <body
        className={`
          antialiased font-sans transition-colors duration-300
          bg-background-light dark:bg-background-dark 
          text-gray-900 dark:text-gray-100
          min-h-screen flex flex-col
        `}
      >
        <ThemeProvider>
          <header className="bg-white/80 dark:bg-card-dark/80 backdrop-blur-md border-b border-gray-100 dark:border-gray-800 sticky top-0 z-50">
            <Navbar />
          </header>

          <main className="grow w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            {children}
          </main>

          <Footer />
        </ThemeProvider>
      </body>
    </html>
  );
}
