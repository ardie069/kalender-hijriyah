<template>
  <div
    class="mt-4 p-4 rounded-box shadow-md"
    :class="[themeClass, 'transition-colors duration-300']"
  >
    <p class="text-lg font-medium">🗓️ Tanggal Hijriyah:</p>
    <hr class="my-2 border-t border-gray-300 dark:border-gray-700 opacity-50" />

    <p v-if="!loading" class="text-lg font-semibold">
      <span v-if="showWeton && wetonText"> {{ wetonText }}, </span>
      <span v-else> {{ weekdayText }}, </span>
      {{ hijriDateText }}
    </p>

    <p v-else class="text-gray-400 italic">Memuat tanggal Hijriyah...</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useThemeStore } from "@/stores/themeStore"; // Menggunakan store

const themeStore = useThemeStore();
const darkMode = computed(() => themeStore.darkMode); // Ambil darkMode dari store

const themeClass = computed(() =>
  darkMode.value ? "bg-base-200 text-base-content" : "bg-zinc-100 text-base"
);

defineProps({
  hijriDateText: String,
  showWeton: Boolean,
  wetonText: String,
  loading: Boolean,
});

const currentDate = ref(new Date());
const weekdayText = ref("");

// Update waktu real-time
const updateTime = () => {
  currentDate.value = new Date();
  let day = currentDate.value.toLocaleDateString("id-ID", { weekday: "long" });

  // Ubah "Minggu" menjadi "Ahad"
  if (day === "Minggu") {
    day = "Ahad";
  }

  weekdayText.value = day;
};

// Real-time update setiap detik
let intervalId;
onMounted(() => {
  updateTime();
  intervalId = setInterval(updateTime, 1000);
});

onUnmounted(() => {
  clearInterval(intervalId);
});
</script>
