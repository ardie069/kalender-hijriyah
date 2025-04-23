<template>
    <!-- Indikator Loading -->
    <p v-if="loading" :class="darkMode ? 'text-yellow-400' : 'text-yellow-800'" class="font-semibold">
        📍 Menunggu lokasi...
    </p>

    <div :class="[hijriDateClass, 'mt-4 p-4 rounded-lg', darkMode ? 'bg-gray-700' : 'bg-gray-200']">
        <div v-if="!loading">
            <p><strong>🕌 Tanggal Hijriyah:</strong> <span>{{ hijriDateText }}</span></p>
        </div>
    </div>

    <div v-if="hijriEndPrediction"
        :class="[hijriEndPredictionClass, 'mt-4 p-4 rounded-lg', darkMode ? 'bg-gray-700' : 'bg-gray-200']">
        <div v-html="hijriEndPrediction"></div>
    </div>
</template>

<script>
export default {
    props: {
        hijriDateText: String,
        hijriEndPrediction: String,
        loading: Boolean,
        darkMode: Boolean,
        hijriDateClass: String,
        hijriEndPredictionClass: String,
        selectedMethod: String,
        API_BASE_URL: String,
        userTimezone: String
    },
    mounted() {
        this.fetchLocationAndHijriDate();
    },
    methods: {
        async fetchLocationAndHijriDate() {
            this.$emit('update:loading', true);
            this.$emit('update:hijriDateText', "📅 Menghitung hasil tanggal Hijriyah... 🔍");
            this.$emit('update:hijriEndPrediction', "📅 Menunggu prediksi akhir bulan... ⏳");

            if ("geolocation" in navigator) {
                navigator.geolocation.getCurrentPosition(async (position) => {
                    this.lat = position.coords.latitude;
                    this.lon = position.coords.longitude;
                    await this.fetchHijriData();
                }, async () => {
                    await this.fetchLocationByIP();
                });
            } else {
                await this.fetchLocationByIP();
            }
        },
        async fetchLocationByIP() {
            try {
                const response = await fetch("https://ip-api.com/json");
                const data = await response.json();
                if (data.lat && data.lon) {
                    this.lat = data.lat;
                    this.lon = data.lon;
                    await this.fetchHijriData();
                } else {
                    this.$emit('update:hijriDateText', "⚠️ Gagal mendapatkan lokasi.");
                    this.$emit('update:loading', false);
                }
            } catch {
                this.$emit('update:hijriDateText', "❌ Gagal mendapatkan lokasi.");
                this.$emit('update:loading', false);
            }
        },
        async fetchHijriData() {
            if (!this.lat || !this.lon || !this.userTimezone) {
                this.$emit('update:hijriDateText', "❌ Lokasi atau zona waktu belum diatur.");
                this.$emit('update:hijriEndPrediction', "<p class='text-red-500'>Data tidak tersedia.</p>");
                this.$emit('update:loading', false);
                return;
            }

            try {
                const API_URL = `${this.API_BASE_URL}/hijri-date`;
                const END_MONTH_URL = `${this.API_BASE_URL}/hijri-end-month`;

                const [dateResponse, endMonthResponse] = await Promise.all([
                    fetch(`${API_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}&timezone=${this.userTimezone}`, { mode: 'cors' }),
                    fetch(`${END_MONTH_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}&timezone=${this.userTimezone}`, { mode: 'cors' })
                ]);

                const dateData = await dateResponse.json();
                const endMonthData = await endMonthResponse.json();

                this.$emit('update:loading', false);

                if (dateData?.hijriDate) {
                    const hijriToday = dateData.hijriDate;
                    const hijriText = `${hijriToday.day} ${this.getHijriMonthName(hijriToday.month)} ${hijriToday.year} H`;
                    this.$emit('update:hijriDateText', hijriText);

                    if (endMonthData?.estimatedEndOfMonth) {
                        const { hijri: estimatedHijri, moonAltitude, elongation, sunAltitude, moonAge, conjunctionBeforeSunset, explanation } = endMonthData.estimatedEndOfMonth;

                        const hijri29 = estimatedHijri ? `${estimatedHijri.day} ${this.getHijriMonthName(estimatedHijri.month)} ${estimatedHijri.year} H` : "Tidak tersedia";
                        const formattedStartGregorian = endMonthData.estimatedStartOfMonth?.gregorian ? this.formatGregorianDate(endMonthData.estimatedStartOfMonth.gregorian) : "Tidak tersedia";

                        const parsedMoonAge = parseFloat(moonAge);
                        const syaratImkanurRukyat = [
                            { name: "Usia Bulan ≥ 8 jam", value: parsedMoonAge, required: 8 },
                            { name: "Ketinggian Bulan ≥ 3°", value: parseFloat(moonAltitude), required: 3 },
                            { name: "Elongasi ≥ 6,4°", value: parseFloat(elongation), required: 6.4 }
                        ];

                        const imkanurRukyat = syaratImkanurRukyat.every(s => !isNaN(s.value) && s.value >= s.required)
                            ? `✅ Memenuhi syarat Imkanur Rukyat → ${explanation}`
                            : "❌ Tidak memenuhi syarat → Bulan ini 30 hari";

                        const imkanurRukyatDetails = syaratImkanurRukyat.map(s => {
                            const satuan = s.name.includes("jam") ? " jam" : "°";
                            const status = !isNaN(s.value) && s.value >= s.required ? "✅" : "❌";
                            return `${status} ${s.name}: ${isNaN(s.value) ? "Tidak tersedia" : s.value.toFixed(2) + satuan}`;
                        }).join("<li>");

                        const now = new Date();
                        const showEndPrediction = hijriToday.day >= 29 && now.getHours() >= 0;

                        const predictionContent = showEndPrediction
                            ? `
                          <div class="bg-gray-100 text-gray-800 p-4 rounded-lg">
                              <h3 class="text-lg font-semibold mb-2">📅 Informasi Tanggal Hijriyah</h3>
                              <p><strong>🗓️ Hari Ini:</strong> ${hijriText}</p>
                              <p><strong>🔮 Perkiraan Akhir Bulan:</strong> ${hijri29}</p>
                              <p><strong>🌙 Awal Bulan Baru:</strong> ${formattedStartGregorian}</p>
                              <br>
                              <p><strong>✅ Validasi Imkanur Rukyat:</strong> ${imkanurRukyat}</p>
                              <ul><li>${imkanurRukyatDetails}</li></ul>
                          </div>
                          `
                            : `<p class="text-gray-500">🔍 Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`;

                        this.$emit('update:hijriEndPrediction', predictionContent);
                    }
                } else {
                    this.$emit('update:hijriDateText', "⚠️ Gagal mendapatkan data.");
                }
            } catch (error) {
                this.$emit('update:hijriDateText', "❌ Terjadi kesalahan saat mengambil data.");
                this.$emit('update:hijriEndPrediction', "<p class='text-red-500'>❌ Gagal mengambil data.</p>");
                this.$emit('update:loading', false);
            }
        },
        getHijriMonthName(month) {
            const months = [
                "Muharram", "Safar", "Rabiul Awal", "Rabiul Akhir", "Jumadil Awal", "Jumadil Akhir",
                "Rajab", "Syaban", "Ramadhan", "Syawal", "Zulkaidah", "Zulhijah"
            ];
            return months[month - 1] || "Tidak diketahui";
        },
        formatGregorianDate({ day, month, year }) {
            const months = [
                "Januari", "Februari", "Maret", "April", "Mei", "Juni",
                "Juli", "Agustus", "September", "Oktober", "November", "Desember"
            ];
            return `${day} ${months[month - 1]} ${year}`;
        }
    },
    watch: {
        selectedMethod(newVal, oldVal) {
            if (newVal !== oldVal) {
                this.fetchHijriData();
            }
        }
    }
};
</script>
