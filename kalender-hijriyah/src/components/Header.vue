<template>
    <div class="flex justify-between items-center mb-4">
        <h1 class="text-2xl font-bold">🕌 Kalender Hijriyah</h1>
        <button @click="toggleTheme" :class="[
            darkMode ? 'bg-gray-700 hover:bg-gray-600 text-white' : 'bg-gray-200 hover:bg-gray-300 text-black',
            'px-3 py-1 rounded-lg transition-all focus:ring-2 focus:ring-gray-500'
        ]" aria-label="Ganti Tema">
            {{ themeToggleText }}
        </button>
    </div>
</template>

<script>
export default {
    props: {
        darkMode: {
            type: Boolean,
            required: true
        },
        themeToggleText: {
            type: String,
            required: true
        },
        toggleTheme: {
            type: Function,
            required: true
        }
    }
};
</script>
