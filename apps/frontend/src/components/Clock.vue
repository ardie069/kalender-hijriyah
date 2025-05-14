<template>
  <div>
    <p class="text-lg font-medium mb-2">{{ currentTime }}</p>
    <p :class="[darkMode ? 'text-gray-400' : 'text-gray-600', 'text-sm mb-4']">
      {{ timezone }}
    </p>
  </div>
</template>

<script>
import { isLocationInJava, getWeton } from "../utils/weton.js";

export default {
  props: ["darkMode", "userTimezone", "lat", "lon"],
  data() {
    return {
      currentTime: "🕒 Memuat Waktu...",
      timezone: "🌍 Zona Waktu: -",
      intervalId: null,
    };
  },
  mounted() {
    this.updateRealTime();
    this.checkAndSetWeton(); // Pastikan untuk memeriksa lokasi dan weton saat pertama kali komponen dimuat
  },
  watch: {
    lat: "checkAndSetWeton", // Menambahkan watcher jika lat berubah
    lon: "checkAndSetWeton", // Menambahkan watcher jika lon berubah
  },
  beforeUnmount() {
    clearInterval(this.intervalId);
  },
  methods: {
    updateRealTime() {
      this.intervalId = setInterval(() => {
        const now = new Date();
        const offset = now.getTimezoneOffset() / -60;

        const weton = getWeton(now);

        const formattedDate = now.toLocaleString("id-ID", {
          weekday: "long",
          year: "numeric",
          month: "long",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });

        // Sisipkan weton setelah hari
        const parts = formattedDate.split(", ");
        if (parts.length >= 2) {
          // Menambahkan weton setelah hari dalam formattedDate
          parts[0] += ` ${weton}`;
        }

        this.currentTime = `🕒 ${parts.join(", ")}`; // Memperbarui currentTime dengan weton di dalamnya
        this.timezone = `🌍 Zona Waktu: ${this.userTimezone} (UTC${
          offset >= 0 ? "+" : ""
        }${offset})`;
      }, 1000);
    },
    checkAndSetWeton() {
      // Pastikan lat dan lon ada, dan lokasi berada di Jawa
      if (this.lat && this.lon && isLocationInJava(this.lat, this.lon)) {
        this.updateRealTime(); // Memperbarui waktu dan menambahkan weton saat lokasi berada di Jawa
      }
    },
  },
};
</script>
