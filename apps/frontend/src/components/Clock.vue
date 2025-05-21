<template>
  <div
    class="p-4 rounded-lg border shadow-sm mb-3 transition-colors duration-300"
    :class="containerClass"
  >
    <p class="text-xl font-semibold mb-1">{{ currentTime }}</p>
    <p class="text-sm" :class="timezoneClass">
      {{ timezone }}
    </p>
  </div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from "vue";
import { useThemeStore } from "@/stores/themeStore";

const props = defineProps(["userTimezone"]);
const theme = useThemeStore();

const currentTime = ref("🕒 Memuat Waktu...");
const timezone = ref("🌍 Zona Waktu: -");
let intervalId = null;

const containerClass = computed(() =>
  theme.darkMode
    ? "bg-gray-800 border-gray-700 text-white"
    : "bg-white border-gray-200 text-gray-800"
);

const timezoneClass = computed(() =>
  theme.darkMode ? "text-gray-400" : "text-gray-600"
);

const updateRealTime = () => {
  clearInterval(intervalId);
  intervalId = setInterval(() => {
    const now = new Date();
    const offset = now.getTimezoneOffset() / -60;

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
    timezone.value = `🌍 Zona Waktu: ${props.userTimezone} (UTC${
      offset >= 0 ? "+" : ""
    }${offset})`;
  }, 1000);
};

onMounted(updateRealTime);
onBeforeUnmount(() => clearInterval(intervalId));
</script>
