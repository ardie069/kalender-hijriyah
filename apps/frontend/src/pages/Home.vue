<script setup lang="ts">
import { computed, ref } from "vue";
import { useThemeStore } from "@/stores/themeStore";

import Clock from "@/components/Clock.vue";
import Method from "@/components/Method.vue";
import HijriDate from "@/components/HijriDate.vue";

import { useHijri } from "@/composables/useHijri";
import type { Method as HijriMethod } from "@/types/hijri";

// Theme
const themeStore = useThemeStore();

const themeClass = computed(() =>
  themeStore.darkMode ? "bg-gray-700 text-white" : "bg-white text-black"
);

// States
const userTimezone = ref<string>(
  Intl.DateTimeFormat().resolvedOptions().timeZone
);

const selectedMethod = ref<HijriMethod>("global");

// Composables
const {
  hijriDate,
  endMonthInfo,
  weton,
  loading,
  error,
} = useHijri(selectedMethod, userTimezone);
</script>

<template>
  <div :class="themeClass" class="flex justify-center items-center min-h-screen transition-colors duration-300">
    <div class="w-full max-w-md rounded-2xl shadow-lg p-6 text-center" :class="themeClass">
      <h1 class="text-2xl font-bold mb-4">🕌 Kalender Hijriyah Hari Ini</h1>
      <p class="mb-6 text-sm opacity-80">
        Menampilkan tanggal Hijriyah berdasarkan lokasi dan metode perhitungan.
      </p>

      <Clock :userTimezone="userTimezone" />

      <Method v-model:selectedMethod="selectedMethod" />

      <HijriDate :hijriDate="hijriDate" :endMonthInfo="endMonthInfo" :weton="weton" :loading="loading"
        :error="error" :method="selectedMethod" />
    </div>
  </div>
</template>