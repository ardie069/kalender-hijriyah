<template>
  <header
    class="flex justify-between items-center px-4 py-3 rounded-lg shadow-sm border mb-4"
    :class="
      darkMode
        ? 'bg-gray-800 border-gray-700 text-white'
        : 'bg-white border-gray-200 text-gray-800'
    "
  >
    <h1 class="text-2xl font-bold flex items-center gap-2">
      🕌 Kalender Hijriyah
    </h1>
    <button
      @click="toggleTheme"
      :class="['btn btn-sm', darkMode ? 'btn-neutral' : 'btn-outline']"
      aria-label="Ganti Tema"
    >
      {{ themeToggleText }}
    </button>
  </header>
</template>

<script>
export default {
  props: {
    darkMode: {
      type: Boolean,
      required: true,
    },
    themeToggleText: {
      type: String,
      required: true,
    },
    toggleTheme: {
      type: Function,
      required: true,
    },
  },
};
</script>
