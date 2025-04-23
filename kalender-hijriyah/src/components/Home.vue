<template>
    <div :class="themeClass" class="flex justify-center items-center min-h-screen transition-colors duration-300">
        <div id="box" :class="themeClass" class="w-full max-w-md bg-gray-800 rounded-2xl shadow-lg p-6 text-center">
            <!-- Menggunakan Komponen Header -->
            <Header :darkMode="darkMode" :themeToggleText="themeToggleText" :toggleTheme="toggleTheme" />

            <!-- Memanggil komponen Clock -->
            <Clock />

            <!-- Dropdown untuk Metode -->
            <div class="mb-4 text-left">
                <label for="method" class="block text-sm font-medium text-gray-300">Metode Perhitungan:</label>
                <select id="method" v-model="selectedMethod" @change="fetchHijriData" :class="[
                    methodClass,
                    'mt-1 block w-full rounded-lg p-2 pl-2 pr-8 focus:ring-2 focus:outline-none appearance-none',
                    darkMode ? 'border-gray-600 focus:ring-gray-500' : 'border-gray-400 focus:ring-gray-300'
                ]">
                    <option value="global">Global</option>
                    <option value="hisab">Hisab</option>
                    <option value="rukyat">Rukyat</option>
                </select>
            </div>

            <!-- Indikator Loading -->
            <p v-if="loading" :class="darkMode ? 'text-yellow-400' : 'text-brown-600'" class="font-semibold">
                📍 Menunggu lokasi...
            </p>

            <!-- Tanggal Masehi & Hijriyah -->
            <div :class="[
                hijriDateClass,
                'mt-4 p-4 rounded-lg',
                darkMode ? 'bg-gray-700' : 'bg-gray-200'
            ]">
                <p><strong>🕌 Tanggal Hijriyah:</strong> <span>{{ hijriDateText }}</span></p>
            </div>

            <div :class="[
                hijriEndPredictionClass,
                'mt-4 p-4 rounded-lg',
                darkMode ? 'bg-gray-700' : 'bg-gray-200'
            ]" v-html="hijriEndPrediction"></div>
        </div>
    </div>
</template>

<script>
import Header from './Header.vue';
import Clock from './Clock.vue';

export default {
    components: {
        Header,
        Clock
    },
    data() {
        return {
            currentTime: "🕒 Memuat Waktu...",
            timezone: "🌍 Zona Waktu: -",
            userTimezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            hijriDateText: "",
            hijriEndPrediction: "",
            selectedMethod: "global",
            loading: true,
            lat: 0,
            lon: 0,
            darkMode: true,
            API_BASE_URL: window.location.hostname.includes("localhost") ? "http://localhost:3000" : "/api",
        };
    },
    computed: {
        themeClass() {
            return this.darkMode ? "bg-gray-700 text-white" : "bg-white text-black";
        },
        methodClass() {
            return this.darkMode ? "bg-gray-700 text-white" : "bg-gray-300 text-black";
        },
        hijriDateClass() {
            return this.darkMode ? "text-white" : "text-black";
        },
        hijriEndPredictionClass() {
            return this.darkMode ? "text-white" : "text-black";
        },
        themeToggleText() {
            return this.darkMode ? "🌞 Light Mode" : "🌙 Dark Mode";
        }
    },
    mounted() {
        this.applySavedTheme();
        this.updateRealTime();
        this.fetchLocationAndHijriDate();
    },
    methods: {
        applySavedTheme() {
            const saved = localStorage.getItem("theme");
            this.darkMode = saved !== "light"; // default to dark
        },
        toggleTheme() {
            this.darkMode = !this.darkMode;
            localStorage.setItem("theme", this.darkMode ? "dark" : "light");
        },
        updateRealTime() {
            setInterval(() => {
                const now = new Date();
                const offset = now.getTimezoneOffset() / -60;
                this.currentTime = `🕒 ${now.toLocaleString("id-ID", {
                    weekday: "long", year: "numeric", month: "long", day: "numeric",
                    hour: "2-digit", minute: "2-digit", second: "2-digit"
                })}`;
                this.timezone = `🌍 Zona Waktu: ${this.userTimezone} (UTC${offset >= 0 ? "+" : ""}${offset})`;
            }, 1000);
        },
        async fetchLocationAndHijriDate() {
            this.hijriDateText = "📅 Menghitung hasil tanggal Hijriyah... 🔍";
            this.hijriEndPrediction = "📅 Menunggu prediksi akhir bulan... ⏳";

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
                    this.hijriDateText = "⚠️ Gagal mendapatkan lokasi.";
                    this.loading = false;
                }
            } catch {
                this.hijriDateText = "❌ Gagal mendapatkan lokasi.";
                this.loading = false;
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
        },
        async fetchHijriData() {
            if (!this.lat || !this.lon || !this.userTimezone) {
                this.hijriDateText = "❌ Lokasi atau zona waktu belum diatur.";
                this.hijriEndPrediction = "<p class='text-red-500'>Data tidak tersedia.</p>";
                this.loading = false;
                return;
            }

            try {
                const API_URL = `${this.API_BASE_URL}/hijri-date`;
                const END_MONTH_URL = `${this.API_BASE_URL}/hijri-end-month`;

                const [dateResponse, endMonthResponse] = await Promise.all([
                    fetch(`${API_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}&timezone=${this.userTimezone}`),
                    fetch(`${END_MONTH_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}&timezone=${this.userTimezone}`)
                ]);

                const dateData = await dateResponse.json();
                const endMonthData = await endMonthResponse.json();

                this.loading = false;

                if (dateData?.hijriDate) {
                    const hijriToday = dateData.hijriDate;
                    this.hijriDateText = `${hijriToday.day} ${this.getHijriMonthName(hijriToday.month)} ${hijriToday.year} H`;

                    if (endMonthData?.estimatedEndOfMonth) {
                        const {
                            hijri: estimatedHijri, moonAltitude, elongation, sunAltitude,
                            moonAge, conjunctionBeforeSunset, explanation
                        } = endMonthData.estimatedEndOfMonth;

                        const hijri29 = estimatedHijri
                            ? `${estimatedHijri.day} ${this.getHijriMonthName(estimatedHijri.month)} ${estimatedHijri.year} H`
                            : "Tidak tersedia";

                        const formattedStartGregorian = endMonthData.estimatedStartOfMonth?.gregorian
                            ? this.formatGregorianDate(endMonthData.estimatedStartOfMonth.gregorian)
                            : "Tidak tersedia";

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
                        if (hijriToday.day >= 29 && now.getHours() >= 0) {
                            this.hijriEndPrediction = `
                  <div class="bg-gray-100 text-gray-800 p-4 rounded-lg">
                    <h3 class="text-lg font-semibold mb-2">📅 Informasi Tanggal Hijriyah</h3>
                    <p><strong>🗓️ Hari Ini:</strong> ${this.hijriDateText}</p>
                    <p><strong>🔮 Perkiraan Akhir Bulan:</strong> ${hijri29}</p>
                    <p><strong>🌙 Awal Bulan Baru:</strong> ${formattedStartGregorian}</p>
                    <br>
                    <p><strong>✅ Validasi Imkanur Rukyat:</strong> ${imkanurRukyat}</p>
                    <ul><li>${imkanurRukyatDetails}</li></ul>
                  </div>
                `;
                        } else {
                            this.hijriEndPrediction = `<p class="text-gray-500">🔍 Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`;
                        }
                    }
                } else {
                    this.hijriDateText = "⚠️ Gagal mendapatkan data.";
                }
            } catch (error) {
                this.hijriDateText = "❌ Terjadi kesalahan saat mengambil data.";
                this.hijriEndPrediction = "<p class='text-red-500'>❌ Gagal mengambil data.</p>";
                this.loading = false;
            }
        }
    }
};
</script>