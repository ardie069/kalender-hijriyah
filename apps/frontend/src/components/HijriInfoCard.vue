<template>
  <div
    class="mt-4 p-4 rounded-box shadow-md"
    :class="
      darkMode ? 'bg-base-200 text-base-content' : 'bg-zinc-100 text-base'
    "
  >
    <p class="text-lg font-medium">🕌 Tanggal Hijriyah:</p>
    <hr class="my-2 border-t border-gray-300 dark:border-gray-700 opacity-50" />
    <p class="text-lg font-semibold">
      <span v-if="showWeton && wetonText">{{ wetonText }}, </span>
      <span v-else-if="weekdayText">{{ weekdayText }}, </span>
      {{ hijriDateText }}
    </p>
  </div>
</template>

<script>
export default {
  props: {
    hijriDateText: String,
    darkMode: Boolean,
    showWeton: Boolean,
    wetonText: String,
    weekdayText: String,
  },
};
</script>
