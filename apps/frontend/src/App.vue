<template>
  <div :class="{ dark: darkMode }">
    <Navbar
      :darkMode="darkMode"
      :themeToggleText="themeToggleText"
      :toggleTheme="toggleTheme"
    />
    <router-view />
  </div>
</template>

<script setup>
import { ref, provide, computed, onMounted } from "vue";
import Navbar from "./components/Navbar.vue";

const darkMode = ref(true);

// Provide `darkMode` ke semua komponen anak
provide("darkMode", darkMode);

// Toggle handler
const toggleTheme = () => {
  darkMode.value = !darkMode.value;
  localStorage.setItem("theme", darkMode.value ? "dark" : "light");
};

// Simpan/ambil dari localStorage
onMounted(() => {
  const saved = localStorage.getItem("theme");
  if (saved) darkMode.value = saved === "dark";
  else {
    localStorage.setItem("theme", darkMode.value ? "dark" : "light");
  }
});

// Optional: teks tombol
const themeToggleText = computed(() =>
  darkMode.value ? "🌙 Mode Gelap" : "☀️ Mode Terang"
);
</script>
