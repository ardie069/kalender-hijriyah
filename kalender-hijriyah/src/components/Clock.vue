<template>
    <div>
        <p class="text-lg font-medium mb-2">{{ currentTime }}</p>
        <p class="text-sm text-gray-400 mb-4">{{ timezone }}</p>
    </div>
</template>

<script>
export default {
    data() {
        return {
            currentTime: "🕒 Memuat Waktu...",
            timezone: "🌍 Zona Waktu: -",
            userTimezone: Intl.DateTimeFormat().resolvedOptions().timeZone
        };
    },
    mounted() {
        this.updateRealTime();
    },
    methods: {
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
        }
    }
};
</script>
