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
export default {
  props: ["darkMode", "userTimezone"],
  data() {
    return {
      currentTime: "🕒 Memuat Waktu...",
      timezone: "🌍 Zona Waktu: -",
      intervalId: null,
    };
  },
  mounted() {
    this.updateRealTime();
  },
  beforeUnmount() {
    clearInterval(this.intervalId);
  },
  methods: {
    updateRealTime() {
      clearInterval(this.intervalId);
      this.intervalId = setInterval(() => {
        const now = new Date();
        const offset = now.getTimezoneOffset() / -60;

        // Mengambil nama hari dan mengubah "Minggu" menjadi "Ahad" jika perlu
        let day = now.toLocaleDateString("id-ID", { weekday: "long" });

        // Ubah "Minggu" menjadi "Ahad"
        if (day === "Minggu") {
          day = "Ahad";
        }

        const rest = now.toLocaleDateString("id-ID", {
          year: "numeric",
          month: "long",
          day: "numeric",
        });

        const time = now.toLocaleTimeString("id-ID", {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });

        this.currentTime = `🕒 ${day}, ${rest}, ${time}`;
        this.timezone = `🌍 Zona Waktu: ${this.userTimezone} (UTC${
          offset >= 0 ? "+" : ""
        }${offset})`;
      }, 1000);
    },
  },
};
</script>
