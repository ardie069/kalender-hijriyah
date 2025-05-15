<template>
  <div
    v-if="hijriEndPrediction"
    class="mt-4 p-4 rounded-box shadow-md prose max-w-none"
    :class="
      darkMode
        ? 'bg-base-200 text-base-content'
        : 'bg-zinc-100 text-base-content'
    "
    v-html="hijriEndPrediction"
  />
</template>

<script>
export default {
  props: {
    hijriEndPrediction: String,
    darkMode: Boolean,
  },
};
</script>
