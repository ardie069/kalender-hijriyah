import MoonHeroSection from "@/components/moon/MoonHeroSection";
import MoonVisualizationCard from "@/components/moon/MoonVisualizationCard";
import MoonSidebar from "@/components/moon/MoonSidebar";
import RukyatWarningBanner from "@/components/moon/RukyatWarningBanner";

export default function MoonInfoPage() {
  return (
    <main className="min-h-screen bg-background-light dark:bg-background-dark transition-colors duration-500 pb-24">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-10">
        <MoonHeroSection />

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* KOLOM KIRI (8 Cols): Main Telemetry */}
          <div className="lg:col-span-8 space-y-8">
            <MoonVisualizationCard
              illumination={14.2}
              phaseName="Waxing Crescent"
              age="3.4"
              distance="384,400"
              elongation="42.5"
            />
          </div>

          {/* KOLOM KANAN (4 Cols): Meta & Religion */}
          <MoonSidebar />
        </div>

        <RukyatWarningBanner />
      </div>
    </main>
  );
}
