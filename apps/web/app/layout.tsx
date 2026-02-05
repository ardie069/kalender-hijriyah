import "./globals.css";
import { ThemeProvider } from "@/context/theme-context";
import Navbar from "@/components/Navbar";

export const metadata = {
  title: "Kalender Hijriyah",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="id">
      <body>
        <ThemeProvider>
          <Navbar />
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
