<template>
    <div :class="themeClass" class="flex justify-center items-center min-h-screen transition-colors duration-300">
        <div id="box" :class="themeClass" class="w-full max-w-md bg-gray-800 rounded-2xl shadow-lg p-6 text-center">
            <!-- Menggunakan Komponen Header -->
            <Header :darkMode="darkMode" :themeToggleText="themeToggleText" :toggleTheme="toggleTheme" />

            <!-- Memanggil komponen Clock -->
            <Clock :darkMode="darkMode" :userTimezone="userTimezone" />

            <!-- Memanggil Komponen Metode -->
            <Method :darkMode="darkMode" v-model:selectedMethod="selectedMethod" />

            <!-- Memanggil Komponen Hijri Date -->
            <HijriDate :hijriDateText="hijriDateText" @update:hijriDateText="hijriDateText = $event"
                :hijriEndPrediction="hijriEndPrediction" @update:hijriEndPrediction="hijriEndPrediction = $event"
                :loading="loading" @update:loading="loading = $event" :darkMode="darkMode"
                :hijriDateClass="hijriDateClass" :hijriEndPredictionClass="hijriEndPredictionClass"
                :selectedMethod="selectedMethod" :API_BASE_URL="API_BASE_URL" :userTimezone="userTimezone" />
        </div>
    </div>
</template>

<script>
import Header from '../components/Header.vue';
import Clock from '../components/Clock.vue';
import Method from '../components/Method.vue';
import HijriDate from '../components/HijriDate.vue';

export default {
    components: {
        Header,
        Clock,
        Method,
        HijriDate
    },
    data() {
        return {
            userTimezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            hijriDateText: "",
            hijriEndPrediction: "",
            selectedMethod: "global",
            loading: true,
            lat: 0,
            lon: 0,
            darkMode: true,
            API_BASE_URL: import.meta.env.VITE_APP_API_BASE_URL || import.meta.env.VUE_APP_API_BASE_URL,
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
    },
    methods: {
        applySavedTheme() {
            const saved = localStorage.getItem("theme");
            this.darkMode = saved !== "light";
        },
        toggleTheme() {
            this.darkMode = !this.darkMode;
            localStorage.setItem("theme", this.darkMode ? "dark" : "light");
        },
    }
};
</script>