<template>
  <div
    class="mt-4 p-4 rounded-box shadow-md"
    :class="
      darkMode ? 'bg-base-200 text-base-content' : 'bg-zinc-100 text-base'
    "
  >
    <p class="text-lg font-medium">🕌 Tanggal Hijriyah:</p>
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
import { ref, onMounted, onUnmounted } from "vue";

defineProps({
  hijriDateText: String,
  darkMode: Boolean,
  showWeton: Boolean,
  wetonText: String,
  loading: Boolean,
});

// Ambil tanggal saat ini
const currentDate = ref(new Date());
const weekdayText = ref("");

// Fungsi untuk memperbarui waktu dan hari
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
  updateTime(); // Update pertama kali
  intervalId = setInterval(updateTime, 1000); // Update setiap detik
});

onUnmounted(() => {
  clearInterval(intervalId); // Bersihkan interval ketika komponen di-unmount
});
</script>
