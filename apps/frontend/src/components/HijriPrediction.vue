<template>
  <div
    v-if="hijriEndPrediction"
    class="mt-4 p-4 rounded-box shadow-md prose max-w-none transition-colors duration-300"
    :class="themeClass"
    v-html="hijriEndPrediction"
  />
</template>

<script setup>
import { computed } from "vue";
import { useThemeStore } from "@/stores/themeStore"; // Menggunakan store untuk tema

// Mengambil darkMode dari store
const themeStore = useThemeStore();
const darkMode = computed(() => themeStore.darkMode);

// Menggunakan computed untuk menentukan kelas tema
const themeClass = computed(() =>
  darkMode.value ? "bg-base-200 text-zinc-100" : "bg-zinc-100 text-base-content"
);

defineProps({
  hijriEndPrediction: String,
});
</script>
