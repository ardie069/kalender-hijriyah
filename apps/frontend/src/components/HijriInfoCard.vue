<template>
  <div
    class="mt-4 p-4 rounded-box shadow-md"
    :class="
      darkMode
        ? 'bg-base-200 text-base-content'
        : 'bg-zinc-100 text-base'
    "
  >
    <p class="text-lg font-medium">
      🕌 Tanggal Hijriyah:
      <span class="font-semibold">{{ hijriDateText }}</span>
    </p>
  </div>
</template>

<script>
export default {
  props: {
    hijriDateText: String,
    darkMode: Boolean,
  },
};
</script>
