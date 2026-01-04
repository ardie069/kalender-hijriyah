import { defineStore } from "pinia";
import { ref, computed, type Ref, type ComputedRef } from "vue";

export const useThemeStore = defineStore("theme", () => {
  const darkMode: Ref<boolean> = ref(false);

  const toggleTheme = (): void => {
    darkMode.value = !darkMode.value;
  };

  const themeToggleText: ComputedRef<string> = computed(() =>
    darkMode.value ? "🌙 Mode Gelap" : "☀️ Mode Terang"
  );

  return {
    darkMode,
    toggleTheme,
    themeToggleText,
  };
});
