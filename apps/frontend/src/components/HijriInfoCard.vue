<script setup lang="ts">
import { computed } from "vue";
import { useThemeStore } from "@/stores/themeStore";
import type { HijriDate, Method } from "@/types/hijri";

const props = defineProps<{
  hijriDate: HijriDate | null;
  weton: string | null;
  method: Method;
}>();

// Theme
const themeStore = useThemeStore();

const themeClass = computed(() =>
  themeStore.darkMode
    ? "bg-base-200 text-base-content"
    : "bg-zinc-100 text-base"
);

//Format
const hijriText = computed(() => {
  if (!props.hijriDate) return "";

  const monthNames = [
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

  return `${props.hijriDate.day} ${monthNames[props.hijriDate.month - 1]
    } ${props.hijriDate.year} H`;
});

const methodLabel = computed(() => {
  switch (props.method) {
    case "global":
      return "Global (Standar Umm al-Qura)";
    case "hisab":
      return "Hisab Wujudul Hilal";
    case "rukyat":
      return "Rukyat Hilal";
    default:
      return props.method;
  }
});
</script>

<template>
  <div class="mt-4 p-4 rounded-box shadow-md transition-colors duration-300" :class="themeClass">
    <p class="text-lg font-medium">🗓️ Tanggal Hijriyah</p>
    <hr class="my-2 border-t opacity-50" />

    <p class="text-lg font-semibold">
      <span v-if="weton">{{ weton }}, </span>
      {{ hijriText }}
    </p>

    <p class="mt-1 text-sm opacity-70">
      Metode: <strong>{{ methodLabel }}</strong>
    </p>
  </div>
</template>
