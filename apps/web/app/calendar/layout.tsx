import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Kalender | Kalender Hijriyah",
  description:
    "Tampilan kalender bulanan Hijriyah dengan informasi hari besar dan konversi Masehi.",
};

export default function CalendarLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
