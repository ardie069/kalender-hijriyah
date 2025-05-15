<template>
  <div
    v-if="showWeton"
    class="mt-4 p-4 badge"
    :class="darkMode ? 'bg-base-200 text-base-content' : 'bg-zinc-100 text-black'"
  >
    <p><strong>📆 Weton:</strong> {{ wetonText }}</p>
  </div>
</template>

<script>
export default {
  props: {
    showWeton: Boolean,
    wetonText: String,
    darkMode: Boolean,
  },
};
</script>
