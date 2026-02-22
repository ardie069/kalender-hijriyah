import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Lunar Analytics | Kalender Hijriyah",
  description:
    "Informasi astronomis bulan secara real-time: altitude hilal, elongasi, umur bulan, dan kriteria visibilitas.",
};

export default function MoonInfoLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
