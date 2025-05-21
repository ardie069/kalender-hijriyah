import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useThemeStore = defineStore("theme", () => {
  const darkMode = ref(false);

  const toggleTheme = () => {
    darkMode.value = !darkMode.value;
  };

  const themeToggleText = computed(() =>
    darkMode.value ? "🌙 Mode Gelap" : "☀️ Mode Terang"
  );

  return {
    darkMode,
    toggleTheme,
    themeToggleText,
  };
});
