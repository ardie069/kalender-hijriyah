<template>
  <div
    :class="themeClass"
    class="min-h-screen p-4 transition-colors duration-300"
  >
    <h1 class="text-2xl font-bold mb-4 mt-4 text-center">🌕 Informasi Bulan</h1>
    <!-- Konten Informasi Bulan -->
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useThemeStore } from "@/stores/themeStore";

const themeStore = useThemeStore();

const themeClass = computed(() =>
  themeStore.darkMode ? "bg-gray-700 text-white" : "bg-white text-black"
);
</script>
