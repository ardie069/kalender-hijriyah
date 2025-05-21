<template>
  <div class="mb-4 text-left">
    <label
      for="method"
      :class="[
        'text-sm mb-1 block',
        darkMode ? 'text-gray-400' : 'text-gray-600',
      ]"
    >
      Metode Perhitungan:
    </label>
    <select
      id="method"
      :value="selectedMethod"
      @change="onMethodChange"
      :class="[
        'select w-full',
        darkMode
          ? 'bg-gray-800 text-white border-gray-600'
          : 'bg-white text-black border-gray-300',
      ]"
    >
      <option value="global">Global</option>
      <option value="hisab">Hisab</option>
      <option value="rukyat">Rukyat</option>
    </select>
  </div>
</template>

<script setup>
import { storeToRefs } from "pinia";
import { useThemeStore } from "@/stores/themeStore";

const themeStore = useThemeStore();
const { darkMode } = storeToRefs(themeStore);

defineProps({
  selectedMethod: String,
});
const emit = defineEmits(["update:selectedMethod"]);

function onMethodChange(event) {
  emit("update:selectedMethod", event.target.value);
}
</script>
