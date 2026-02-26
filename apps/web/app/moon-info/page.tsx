"use client";

import { useState } from "react";
import { useMoon } from "@/hooks/use-moon";
import { useRukyat } from "@/hooks/use-rukyat";
import MoonHeroSection from "@/components/moon/MoonHeroSection";
import MoonVisualizationCard from "@/components/moon/MoonVisualizationCard";
import MoonVisibilitySimulator from "@/components/moon/MoonVisibilitySimulator";
import MoonSidebar from "@/components/moon/MoonSidebar";
import RukyatWarningBanner from "@/components/moon/RukyatWarningBanner";
import MoonSkeleton from "@/components/moon/MoonSkeleton";

export default function MoonInfoPage() {
  const { data, loading, error, lat } = useMoon();
  const { data: rukyatData } = useRukyat();
  const [method, setMethod] = useState<"hisab" | "rukyat">("hisab");

  if (loading) return <MoonSkeleton />;

  return (
    <main className="min-h-screen bg-background-light dark:bg-background-dark transition-colors duration-500 pb-24">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-10">
        <MoonHeroSection />

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-xs font-bold uppercase tracking-widest text-center">
            ⚠️ {error}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-8 space-y-8">
            {method === "hisab" && (
              <MoonVisualizationCard
                illumination={data?.telemetry.illumination ?? 0}
                phaseName={data?.status.phase_name ?? "Loading..."}
                age={data?.telemetry.age_days ?? 0}
                distance={data?.telemetry.distance_km ?? "0"}
                elongation={data?.telemetry.elongation ?? 0}
                altitude={data?.telemetry.altitude ?? 0}
                lat={lat ?? 0}
                is_rukyat_time={data?.status.is_rukyat_time ?? false}
              />
            )}

            {method === "rukyat" && data && (
              <MoonVisibilitySimulator
                telemetry={data.telemetry}
                status={data.status}
                rukyatData={rukyatData}
              />
            )}
          </div>
          <MoonSidebar
            telemetry={data?.telemetry}
            method={method}
            onMethodChange={setMethod}
          />
        </div>

        {data?.status.is_rukyat_time && (
          <RukyatWarningBanner isMet={!!data?.status.is_visible} />
        )}
      </div>
    </main>
  );
}
