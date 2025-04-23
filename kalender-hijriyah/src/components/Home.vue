<template>
    <div class="bg-gray-900 text-white flex justify-center items-center min-h-screen transition-all duration-300">
        <div class="w-full max-w-md bg-gray-800 rounded-2xl shadow-lg p-6 text-center transition-all duration-300">
            <!-- Header dengan Toggle Mode -->
            <div class="flex justify-between items-center mb-4">
                <h1 class="text-2xl font-bold">🕌 Kalender Hijriyah</h1>
                <button @click="toggleTheme"
                    class="bg-gray-700 text-white px-3 py-1 rounded-lg transition-all hover:bg-gray-600 focus:ring-2 focus:ring-gray-500"
                    aria-label="Ganti Tema">
                    {{ isDark ? '🌞 Light Mode' : '🌙 Dark Mode' }}
                </button>
            </div>

            <!-- Indikator Loading -->
            <p v-if="loading" class="text-yellow-400 font-semibold">📍 Menunggu lokasi...</p>

            <!-- Kalender Hijriyah -->
            <Calendar :method="method" :timezone="timezone" v-if="!loading" /> <!-- Hanya tampilkan setelah loading selesai -->
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Calendar from './Calendar.vue'; // Impor Calendar.vue

const isDark = ref(true);
const currentTime = ref('🕒 Loading waktu...');
const timezone = ref('-');
const method = ref('global');
const hijriDate = ref('-');
const predictionText = ref('📅 Menunggu prediksi akhir bulan...');
const loading = ref(true);

// Fungsi untuk Toggle Mode Tema
const toggleTheme = () => {
    isDark.value = !isDark.value;
    document.documentElement.classList.toggle('dark', isDark.value);
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light');
};

// Fungsi untuk memperbarui waktu setiap detik
const updateTime = () => {
    currentTime.value = new Date().toLocaleTimeString();
};

// Memperbarui waktu setiap detik menggunakan setInterval
onMounted(() => {
    timezone.value = Intl.DateTimeFormat().resolvedOptions().timeZone;
    setInterval(updateTime, 1000); // Update setiap detik
    loading.value = false;

    // Simulasi data
    hijriDate.value = '29 Syawal 1446 H';
    predictionText.value = '🌙 Menunggu visibilitas hilal...';
});
</script>

<style scoped>
html.dark {
    background-color: #111827;
}
</style>
