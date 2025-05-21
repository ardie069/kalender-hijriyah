<template>
  <div
    :class="themeClass"
    class="flex justify-center items-center min-h-screen transition-colors duration-300"
  >
    <div
      id="box"
      :class="themeClass"
      class="w-full max-w-md bg-gray-800 rounded-2xl shadow-lg p-6 text-center"
    >
      <!-- Judul Halaman -->
      <h1 class="text-2xl font-bold mb-4">🕌 Kalender Hijriyah Hari Ini</h1>
      <p class="mb-6 text-sm opacity-80">
        Menampilkan tanggal Hijriyah saat ini berdasarkan lokasi dan metode
        perhitungan.
      </p>

      <!-- Komponen Clock -->
      <Clock :userTimezone="userTimezone" />

      <!-- Komponen Metode -->
      <Method v-model:selectedMethod="selectedMethod" />

      <!-- Komponen Hijri Date -->
      <HijriDate
        :hijriDateText="hijriDateText"
        :hijriEndPrediction="hijriEndPrediction"
        :loading="loading"
        :hijriDateClass="hijriDateClass"
        :hijriEndPredictionClass="hijriEndPredictionClass"
        :selectedMethod="selectedMethod"
        :API_BASE_URL="API_BASE_URL"
        :userTimezone="userTimezone"
      />
    </div>
  </div>
</template>

<script>
import { useThemeStore } from "@/stores/themeStore";
import { computed, ref, watch } from "vue";
import Clock from "@/components/Clock.vue";
import Method from "@/components/Method.vue";
import HijriDate from "@/components/HijriDate.vue";

export default {
  components: {
    Clock,
    Method,
    HijriDate,
  },
  setup() {
    const themeStore = useThemeStore();
    const userTimezone = ref(Intl.DateTimeFormat().resolvedOptions().timeZone);
    const selectedMethod = ref("global");
    const hijriDateText = ref("");
    const hijriEndPrediction = ref("");
    const loading = ref(true);
    const API_BASE_URL = import.meta.env.VITE_APP_API_BASE_URL;

    // Menggunakan store untuk tema
    const themeClass = computed(() =>
      themeStore.darkMode ? "bg-gray-700 text-white" : "bg-white text-black"
    );
    const hijriDateClass = computed(() =>
      themeStore.darkMode ? "text-white" : "text-black"
    );
    const hijriEndPredictionClass = computed(() =>
      themeStore.darkMode ? "text-white" : "text-black"
    );

    // Watchers untuk memperbarui data jika props berubah
    watch([selectedMethod, userTimezone], () => {
      loading.value = true;
      // Logika untuk fetch data hijri berdasarkan selectedMethod dan userTimezone
      // Set hijriDateText, hijriEndPrediction dan loading menjadi false setelah data diterima
    });

    return {
      userTimezone,
      selectedMethod,
      hijriDateText,
      hijriEndPrediction,
      loading,
      API_BASE_URL,
      themeClass,
      hijriDateClass,
      hijriEndPredictionClass,
    };
  },
};
</script>
