<script setup lang="ts">
import { computed } from "vue";
import { useThemeStore } from "@/stores/themeStore";
import type { HijriDate } from "@/types/hijri";

// Props
interface HijriPredictionProps {
  prediction: {
    today?: HijriDate;
    estimated_end_of_month?: HijriDate | null;
    message?: string;
    visibility?: {
      moon_altitude: number;
      elongation: number;
      moon_age: number;
      is_visible: boolean;
    };
  } | null;
}

const props = defineProps<HijriPredictionProps>();

// Theme
const themeStore = useThemeStore();

const themeClass = computed(() =>
  themeStore.darkMode
    ? "bg-base-200 text-base-content"
    : "bg-zinc-100 text-base"
);

// Format
function formatHijri(date: HijriDate): string {
  const months = [
    "Muharam",
    "Safar",
    "Rabiulawal",
    "Rabiulakhir",
    "Jumadilawal",
    "Jumadilakhir",
    "Rajab",
    "Syakban",
    "Ramadan",
    "Syawal",
    "Zulkaidah",
    "Zulhijah",
  ];

  return `${date.day} ${months[date.month - 1]} ${date.year} H`;
}
</script>

<template>
  <div v-if="prediction" class="mt-4 p-4 rounded-box shadow-md transition-colors duration-300" :class="themeClass">
    <h3 class="text-lg font-semibold mb-2">🌙 Informasi Akhir Bulan</h3>

    <!-- Pesan Umum -->
    <p v-if="prediction.message" class="text-sm opacity-80">
      {{ prediction.message }}
    </p>

    <!-- Prediksi Akhir Bulan -->
    <div v-if="prediction.estimated_end_of_month" class="mt-2">
      <p>
        <strong>Perkiraan akhir bulan:</strong>
        {{ formatHijri(prediction.estimated_end_of_month) }}
      </p>
    </div>

    <!-- Validasi Rukyat -->
    <div v-if="prediction.visibility" class="mt-3 text-sm">
      <p class="font-medium">Validasi Imkanur Rukyat:</p>
      <ul class="list-disc list-inside">
        <li>
          Ketinggian bulan:
          {{ prediction.visibility.moon_altitude.toFixed(2) }}°
        </li>
        <li>
          Elongasi:
          {{ prediction.visibility.elongation.toFixed(2) }}°
        </li>
        <li>
          Usia bulan:
          {{ prediction.visibility.moon_age.toFixed(2) }} jam
        </li>
        <li>
          Status:
          <span :class="prediction.visibility.is_visible
            ? 'text-green-500 font-semibold'
            : 'text-red-500 font-semibold'
            ">
            {{ prediction.visibility.is_visible
              ? 'Memenuhi syarat'
              : 'Tidak memenuhi syarat' }}
          </span>
        </li>
      </ul>
    </div>
  </div>
</template>
