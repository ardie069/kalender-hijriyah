<template>
  <p v-if="loading" :class="loadingClass" class="font-semibold">
    📍 Menunggu lokasi...
  </p>

  <div
    :class="[
      hijriDateClass,
      'mt-4 p-4 rounded-lg',
      darkMode ? 'bg-gray-700' : 'bg-gray-200',
    ]"
  >
    <div v-if="!loading">
      <p>
        <strong>🕌 Tanggal Hijriyah:</strong> <span>{{ hijriDateText }}</span>
      </p>
    </div>
  </div>

  <div
    v-if="hijriEndPrediction"
    :class="[
      hijriEndPredictionClass,
      'mt-4 p-4 rounded-lg',
      darkMode ? 'bg-gray-700' : 'bg-gray-200',
    ]"
    v-html="hijriEndPrediction"
  />
</template>

<script>
import { isLocationInJava, getWeton } from "../utils/weton.js";

export default {
  props: {
    hijriDateText: String,
    hijriEndPrediction: String,
    showWeton: Boolean,
    wetonText: String,
    loading: Boolean,
    darkMode: Boolean,
    hijriDateClass: String,
    hijriEndPredictionClass: String,
    selectedMethod: String,
    API_BASE_URL: String,
    userTimezone: String,
  },
  data() {
    return {
      lat: null,
      lon: null,
      showWeton: false,
      wetonText: "",
    };
  },
  computed: {
    loadingClass() {
      return this.darkMode ? "text-yellow-400" : "text-yellow-800";
    },
  },
  mounted() {
    this.fetchLocationAndHijriDate();
  },
  methods: {
    async fetchLocationAndHijriDate() {
      this.$emit("update:loading", true);
      this.$emit(
        "update:hijriDateText",
        "📅 Menghitung hasil tanggal Hijriyah... 🔍"
      );
      this.$emit(
        "update:hijriEndPrediction",
        "📅 Menunggu prediksi akhir bulan... ⏳"
      );

      if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition(
          async ({ coords }) => {
            this.lat = coords.latitude;
            this.lon = coords.longitude;
            await this.fetchHijriData();
            this.checkAndSetWeton();
          },
          async () => {
            await this.fetchLocationByIP();
          }
        );
      } else {
        await this.fetchLocationByIP();
      }
    },

    async fetchLocationByIP() {
      try {
        const res = await fetch("https://ip-api.com/json");
        const { lat, lon } = await res.json();
        if (lat && lon) {
          this.lat = lat;
          this.lon = lon;
          await this.fetchHijriData();
          this.checkAndSetWeton();
        } else {
          this.showError("⚠️ Gagal mendapatkan lokasi.");
        }
      } catch {
        this.showError("❌ Gagal mendapatkan lokasi.");
      }
    },

    checkAndSetWeton() {
      if (isLocationInJava(this.lat, this.lon)) {
        this.showWeton = true;
        this.wetonText = getWeton(new Date());
      } else {
        this.showWeton = false;
        this.wetonText = "";
      }

      this.$emit("update:showWeton", this.showWeton);
      this.$emit("update:wetonText", this.wetonText);
    },

    showError(message) {
      this.$emit("update:hijriDateText", message);
      this.$emit(
        "update:hijriEndPrediction",
        "<p class='text-red-500'>Data tidak tersedia.</p>"
      );
      this.$emit("update:loading", false);
    },

    getHijriMonthName(month) {
      const months = [
        "Muharam",
        "Safar",
        "Rabiulawal",
        "Rabiulakhir",
        "Jumadilawal",
        "Jumadilakhir",
        "Rajab",
        "Syakban",
        "Ramadan",
        "Syawal",
        "Zulkaidah",
        "Zulhijah",
      ];
      return months[month - 1] || "Tidak diketahui";
    },

    formatGregorianDate({ day, month, year }) {
      const months = [
        "Januari",
        "Februari",
        "Maret",
        "April",
        "Mei",
        "Juni",
        "Juli",
        "Agustus",
        "September",
        "Oktober",
        "November",
        "Desember",
      ];
      return `${day} ${months[month - 1]} ${year}`;
    },

    checkImkanurRukyat({ moonAge, moonAltitude, elongation }) {
      const criteria = [
        { name: "Usia Bulan ≥ 8 jam", value: parseFloat(moonAge), required: 8 },
        {
          name: "Ketinggian Bulan ≥ 3°",
          value: parseFloat(moonAltitude),
          required: 3,
        },
        {
          name: "Elongasi ≥ 6,4°",
          value: parseFloat(elongation),
          required: 6.4,
        },
      ];

      const isValid = criteria.every(
        (c) => !isNaN(c.value) && c.value >= c.required
      );
      const summary = isValid
        ? "✅ Memenuhi syarat Imkanur Rukyat"
        : "❌ Tidak memenuhi syarat → Bulan ini 30 hari";

      const details = criteria
        .map((c) => {
          const satuan = c.name.includes("jam") ? " jam" : "°";
          const status = !isNaN(c.value) && c.value >= c.required ? "✅" : "❌";
          return `${status} ${c.name}: ${
            isNaN(c.value) ? "Tidak tersedia" : c.value.toFixed(2) + satuan
          }`;
        })
        .join("<li>");

      return { summary, details };
    },

    async fetchHijriData() {
      if (!this.lat || !this.lon || !this.userTimezone) {
        this.showError("❌ Lokasi atau zona waktu belum diatur.");
        return;
      }

      try {
        const params = `lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}&timezone=${this.userTimezone}`;
        const [dateRes, endRes] = await Promise.all([
          // fetch(`http://localhost:3000/api/hijri-date?${params}`), // development
          // fetch(`http://localhost:3000/api/hijri-end-month?${params}`), // development
          fetch(`${this.API_BASE_URL}/hijri-date?${params}`), // production
          fetch(`${this.API_BASE_URL}/hijri-end-month?${params}`), // production
        ]);

        const dateData = await dateRes.json();
        const endMonthData = await endRes.json();

        this.$emit("update:loading", false);

        if (!dateData?.hijriDate) {
          return this.showError("⚠️ Gagal mendapatkan data.");
        }

        const hijriToday = dateData.hijriDate;
        const hijriText = `${hijriToday.day} ${this.getHijriMonthName(
          hijriToday.month
        )} ${hijriToday.year} H`;
        this.$emit("update:hijriDateText", hijriText);

        if (!endMonthData) return;

        const { message, estimatedEndOfMonth, estimatedStartOfMonth } =
          endMonthData;

        if (message) {
          this.$emit(
            "update:hijriEndPrediction",
            `<p class="text-gray-500">🔍 ${message}</p>`
          );
          return;
        }

        const showPrediction = hijriToday.day >= 29;
        if (!showPrediction) {
          this.$emit(
            "update:hijriEndPrediction",
            `<p class="text-gray-500">⚠️ Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`
          );
          return;
        }

        const hijri29 = estimatedEndOfMonth?.hijri
          ? `${estimatedEndOfMonth.hijri.day} ${this.getHijriMonthName(
              estimatedEndOfMonth.hijri.month
            )} ${estimatedEndOfMonth.hijri.year} H`
          : "Tidak tersedia";

        const formattedStartGregorian = estimatedStartOfMonth?.gregorian
          ? this.formatGregorianDate(estimatedStartOfMonth.gregorian)
          : "Tidak tersedia";

        const { summary, details } =
          this.checkImkanurRukyat(estimatedEndOfMonth);

        const predictionContent = `
            <div class="bg-gray-100 text-gray-800 p-4 rounded-lg">
              <h3 class="text-lg font-semibold mb-2">📅 Informasi Tanggal Hijriyah</h3>
              <p><strong>🗓️ Hari Ini:</strong> ${hijriText}</p>
              <p><strong>🔮 Perkiraan Akhir Bulan:</strong> ${hijri29}</p>
              <p><strong>🌙 Awal Bulan Baru:</strong> ${formattedStartGregorian}</p>
              <br>
              <p><strong>✅ Validasi Imkanur Rukyat:</strong> ${summary}</p>
              <ul><li>${details}</li></ul>
            </div>
          `;
        this.$emit("update:hijriEndPrediction", predictionContent);
      } catch (err) {
        this.showError("❌ Terjadi kesalahan saat mengambil data.");
      }
    },
  },
  watch: {
    selectedMethod(newVal, oldVal) {
      if (newVal !== oldVal) {
        this.fetchHijriData();
      }
    },
    userTimezone() {
      this.checkAndSetWeton();
      this.fetchHijriData();
    },
  },
};
</script>
