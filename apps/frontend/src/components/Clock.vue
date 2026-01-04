<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from "vue";
import { useThemeStore } from "@/stores/themeStore";

// Props
const props = defineProps<{
  userTimezone: string;
}>();

// Theme
const themeStore = useThemeStore();

// State
const currentTime = ref<string>("🕒 Memuat waktu...");
const timezoneText = ref<string>("🌍 Zona Waktu: -");

let intervalId: number | undefined;

//Computed Style
const containerClass = computed(() =>
  themeStore.darkMode
    ? "bg-gray-800 border-gray-700 text-white"
    : "bg-white border-gray-200 text-gray-800"
);

const timezoneClass = computed(() =>
  themeStore.darkMode ? "text-gray-400" : "text-gray-600"
);

// Logic
function updateRealTime() {
  if (intervalId) clearInterval(intervalId);

  intervalId = window.setInterval(() => {
    const now = new Date();

    // Nama hari (Ahad override)
    let day = now.toLocaleDateString("id-ID", { weekday: "long" });
    if (day === "Minggu") day = "Ahad";

    const date = now.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });

    const time = now.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });

    currentTime.value = `🕒 ${day}, ${date}, ${time}`;

    // Offset UTC dari browser (bukan timezone target, tapi cukup informatif)
    const offset = now.getTimezoneOffset() / -60;
    timezoneText.value = `🌍 Zona Waktu: ${props.userTimezone} (UTC${offset >= 0 ? "+" : ""
      }${offset})`;
  }, 1000);
}

// Lifecycle
onMounted(updateRealTime);

onBeforeUnmount(() => {
  if (intervalId) clearInterval(intervalId);
});
</script>

<template>
  <div class="p-4 rounded-lg border shadow-sm mb-3 transition-colors duration-300" :class="containerClass">
    <p class="text-xl font-semibold mb-1">{{ currentTime }}</p>
    <p class="text-sm" :class="timezoneClass">
      {{ timezoneText }}
    </p>
  </div>
</template>