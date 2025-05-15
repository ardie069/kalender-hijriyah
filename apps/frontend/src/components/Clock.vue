<template>
  <div
    class="p-4 rounded-lg border shadow-sm mb-3"
    :class="
      darkMode
        ? 'bg-gray-800 border-gray-700 text-white'
        : 'bg-white border-gray-200 text-gray-800'
    "
  >
    <p class="text-xl font-semibold mb-1">{{ currentTime }}</p>
    <p class="text-sm" :class="darkMode ? 'text-gray-400' : 'text-gray-600'">
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
    this.checkAndSetWeton();
  },
  watch: {
    lat: "checkAndSetWeton",
    lon: "checkAndSetWeton",
  },
  beforeUnmount() {
    clearInterval(this.intervalId);
  },
  methods: {
    updateRealTime() {
      clearInterval(this.intervalId); // Bersihkan sebelumnya untuk mencegah dobel timer
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

        const parts = formattedDate.split(", ");
        if (parts.length >= 2) {
          parts[0] += ` ${weton}`;
        }

        this.currentTime = `🕒 ${parts.join(", ")}`;
        this.timezone = `🌍 Zona Waktu: ${this.userTimezone} (UTC${
          offset >= 0 ? "+" : ""
        }${offset})`;
      }, 1000);
    },
    checkAndSetWeton() {
      if (this.lat && this.lon && isLocationInJava(this.lat, this.lon)) {
        this.updateRealTime();
      }
    },
  },
};
</script>
