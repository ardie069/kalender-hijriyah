<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useThemeStore } from "@/stores/themeStore";
import Navbar from "@/components/Navbar.vue";

const themeStore = useThemeStore();

/**
 * expose state
 */
const darkMode = computed(() => themeStore.darkMode);
const themeToggleText = computed(() => themeStore.themeToggleText);
const toggleTheme = themeStore.toggleTheme;

/**
 * sync localStorage → store
 */
onMounted(() => {
  const saved = localStorage.getItem("theme");
  if (saved === "dark") {
    themeStore.darkMode = true;
  } else if (saved === "light") {
    themeStore.darkMode = false;
  } else {
    localStorage.setItem("theme", themeStore.darkMode ? "dark" : "light");
  }
});

/**
 * watch change → persist
 */
import { watch } from "vue";

watch(
  () => themeStore.darkMode,
  (val) => {
    localStorage.setItem("theme", val ? "dark" : "light");
  }
);
</script>

<template>
  <div :class="{ dark: darkMode }">
    <Navbar :darkMode="darkMode" :themeToggleText="themeToggleText" :toggleTheme="toggleTheme" />
    <router-view />
  </div>
</template>
