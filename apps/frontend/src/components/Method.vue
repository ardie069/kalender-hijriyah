<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useThemeStore } from "@/stores/themeStore";
import type { Method } from "@/types/hijri";

// Theme
const themeStore = useThemeStore();
const { darkMode } = storeToRefs(themeStore);

// Props & Emits
const props = defineProps<{
  selectedMethod: Method;
}>();

const emit = defineEmits<{
  (e: "update:selectedMethod", value: Method): void;
}>();

// Handler
function onMethodChange(event: Event) {
  const value = (event.target as HTMLSelectElement).value as Method;
  emit("update:selectedMethod", value);
}
</script>

<template>
  <div class="mb-4 text-left">
    <label for="method" :class="[
      'text-sm mb-1 block',
      darkMode ? 'text-gray-400' : 'text-gray-600',
    ]">
      Metode Perhitungan:
    </label>

    <select id="method" :value="selectedMethod" @change="onMethodChange" :class="[
      'select w-full',
      darkMode
        ? 'bg-gray-800 text-white border-gray-600'
        : 'bg-white text-black border-gray-300',
    ]">
      <option value="global">Global</option>
      <option value="hisab">Hisab</option>
      <option value="rukyat">Rukyat</option>
    </select>
  </div>
</template>