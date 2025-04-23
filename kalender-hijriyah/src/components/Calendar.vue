<template>
    <div id="box" class="w-full max-w-md bg-gray-800 rounded-2xl shadow-lg p-6 text-center transition-all duration-300">
        <p id="current-time" class="text-lg text-white font-medium mb-2">{{ currentTime }}</p>
        <p id="timezone" class="text-sm text-gray-400 mb-4">🌍 Zona Waktu: {{ timezone }}</p>

        <div class="mb-4 text-left">
            <label for="method" class="block text-sm font-medium text-gray-300">Metode Perhitungan:</label>
            <select v-model="selectedMethod"
                class="mt-1 block w-full bg-gray-700 border border-gray-600 rounded-lg p-2 text-white">
                <option value="global">Global</option>
                <option value="hisab">Hisab</option>
                <option value="rukyat">Rukyat</option>
            </select>
        </div>

        <p id="loading" v-if="loading" class="text-yellow-400 font-semibold">📍 Menunggu lokasi...</p>

        <div class="mt-4 p-4 bg-gray-700 rounded-lg" id="date-box">
            <p><strong>🕌 Tanggal Hijriyah:</strong> <span>{{ hijriDate }}</span></p>
        </div>

        <div class="mt-4 p-4 bg-gray-700 rounded-lg" id="hijri-end-prediction">
            <p>{{ hijriEndPrediction }}</p>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            currentTime: '🕒 Loading time...',
            timezone: '',
            hijriDate: '-',
            hijriEndPrediction: '📅 Menunggu prediksi akhir bulan...',
            loading: true,
            lat: 0,
            lon: 0,
            selectedMethod: 'global',
            intervalId: null, // Untuk menyimpan ID interval
        };
    },
    methods: {
        updateRealTime() {
            this.intervalId = setInterval(() => {
                const now = new Date();
                const offset = now.getTimezoneOffset() / -60;
                this.currentTime = `🕒 ${now.toLocaleString('id-ID', {
                    weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
                    hour: '2-digit', minute: '2-digit', second: '2-digit',
                })}`;
                this.timezone = `${Intl.DateTimeFormat().resolvedOptions().timeZone} (UTC${offset >= 0 ? '+' : ''}${offset})`;
            }, 1000);
        },
        stopRealTime() {
            clearInterval(this.intervalId);
        },
        loadTheme() {
            const savedTheme = localStorage.getItem('theme');
            if (savedTheme === 'dark') {
                document.body.classList.add('bg-gray-900');
                document.body.classList.add('text-white');
            } else {
                document.body.classList.add('bg-white');
                document.body.classList.add('text-black');
            }
        },
        async fetchLocationAndHijriData() {
            this.loading = true;
            this.hijriDate = '📅 Menghitung hasil tanggal Hijriyah... 🔍';
            this.hijriEndPrediction = '📅 Menunggu prediksi akhir bulan... ⏳';

            if (navigator.geolocation) {
                navigator.geolocation.getCurrentPosition(
                    async (position) => {
                        this.lat = position.coords.latitude;
                        this.lon = position.coords.longitude;
                        await this.fetchHijriData();
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
                const response = await fetch('https://ip-api.com/json');
                const data = await response.json();
                if (data.lat && data.lon) {
                    this.lat = data.lat;
                    this.lon = data.lon;
                    await this.fetchHijriData();
                } else {
                    this.hijriDate = '⚠️ Gagal mendapatkan lokasi.';
                    this.loading = false;
                }
            } catch (error) {
                this.hijriDate = '❌ Gagal mendapatkan lokasi.';
                this.loading = false;
            }
        },
        async fetchHijriData() {
            const API_URL = `/api/hijri-date`;
            const END_MONTH_URL = `/api/hijri-end-month`;

            try {
                const [dateResponse, endMonthResponse] = await Promise.all([ 
                    fetch(`${API_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}`),
                    fetch(`${END_MONTH_URL}?lat=${this.lat}&lon=${this.lon}&method=${this.selectedMethod}`),
                ]);

                if (!dateResponse.ok || !endMonthResponse.ok) throw new Error('HTTP Error saat mengambil data.');

                const dateData = await dateResponse.json();
                const endMonthData = await endMonthResponse.json();

                this.loading = false;

                if (dateData?.hijriDate) {
                    this.hijriDate = `${dateData.hijriDate.day} ${this.getHijriMonthName(dateData.hijriDate.month)} ${dateData.hijriDate.year} H`;
                }

                if (endMonthData?.estimatedEndOfMonth) {
                    this.hijriEndPrediction = `📅 Prediksi akhir bulan: ${endMonthData.estimatedEndOfMonth.hijri.day}`;
                }
            } catch (error) {
                this.hijriDate = '❌ Terjadi kesalahan saat mengambil data.';
                this.hijriEndPrediction = '❌ Gagal mengambil data.';
            }
        },
        getHijriMonthName(month) {
            const hijriMonths = [
                'Muharram', 'Safar', 'Rabiul Awal', 'Rabiul Akhir', 'Jumadil Awal', 'Jumadil Akhir',
                'Rajab', 'Syakban', 'Ramadan', 'Syawal', 'Zulkaidah', 'Zulhijah'
            ];
            return hijriMonths[month - 1] || 'Tidak diketahui';
        }
    },
    mounted() {
        this.updateRealTime();
        this.fetchLocationAndHijriData();
        this.loadTheme();
    },
    beforeDestroy() {
        this.stopRealTime(); // Membersihkan interval ketika komponen dihancurkan
    }
};
</script>

<style scoped>
/* Tambahkan gaya CSS jika diperlukan */
</style>
