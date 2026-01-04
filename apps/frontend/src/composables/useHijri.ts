import { ref, watch, type Ref } from "vue";
import type { Method, HijriDate } from "@/types/hijri";
import { fetchHijriDate, fetchHijriEndMonth } from "@/api/hijri";
import { isLocationInJava, getWeton } from "@/utils/weton";

export function useHijri(method: Ref<Method>, timezone: Ref<string>) {
  const lat = ref<number | null>(null);
  const lon = ref<number | null>(null);

  const hijriDate = ref<HijriDate | null>(null);
  const endMonthInfo = ref<Record<string, any> | null>(null);

  const weton = ref<string | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function resolveLocation() {
    return new Promise<{ lat: number; lon: number }>((resolve, reject) => {
      if (!navigator.geolocation) reject();

      navigator.geolocation.getCurrentPosition(
        (pos) =>
          resolve({
            lat: pos.coords.latitude,
            lon: pos.coords.longitude,
          }),
        reject
      );
    });
  }

  async function loadHijri() {
    if (!timezone.value) return;

    loading.value = true;
    error.value = null;

    try {
      if (!lat.value || !lon.value) {
        const loc = await resolveLocation();
        lat.value = loc.lat;
        lon.value = loc.lon;
      }

      const [dateRes, endRes] = await Promise.all([
        fetchHijriDate(lat.value!, lon.value!, method.value, timezone.value),
        fetchHijriEndMonth(
          lat.value!,
          lon.value!,
          method.value,
          timezone.value
        ),
      ]);

      hijriDate.value = dateRes.hijri_date;
      endMonthInfo.value = endRes;

      if (isLocationInJava(lat.value!, lon.value!)) {
        const day = new Date().toLocaleDateString("id-ID", { weekday: "long" });
        weton.value = `${day === "Minggu" ? "Ahad" : day} ${getWeton(
          new Date()
        )}`;
      } else {
        weton.value = null;
      }
    } catch {
      error.value = "Gagal mengambil data Hijriyah";
    } finally {
      loading.value = false;
    }
  }

  watch([method, timezone], loadHijri, { immediate: true });

  return {
    hijriDate,
    endMonthInfo,
    weton,
    loading,
    error,
    reload: loadHijri,
  };
}
