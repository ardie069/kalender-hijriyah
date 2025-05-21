<template>
  <div :class="{ dark: darkMode }">
    <!-- Menambahkan class dark pada elemen root -->
    <div v-if="loading" class="flex items-center justify-center mt-4">
      <span :class="['loading loading-spinner mr-2', loadingClass]"></span>
      <span :class="['font-semibold', loadingClass]">
        📍 Menunggu lokasi...
      </span>
    </div>

    <HijriInfoCard
      :hijriDateText="hijriDateText"
      :showWeton="showWeton"
      :wetonText="wetonText"
      :loading="loading"
    />
    <HijriPrediction
      :hijriEndPrediction="hijriEndPrediction"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { useThemeStore } from "@/stores/themeStore";
import HijriInfoCard from "./HijriInfoCard.vue";
import HijriPrediction from "./HijriPrediction.vue";
import { useHijriDate } from "../composables/useHijriDate";

const themeStore = useThemeStore();
const darkMode = computed(() => themeStore.darkMode);

// props dari parent
const props = defineProps({
  selectedMethod: String,
  userTimezone: String,
  API_BASE_URL: String,
});

// reactive props
const selectedMethod = ref(props.selectedMethod);
const userTimezone = ref(props.userTimezone);

// composable
const {
  hijriDateText,
  hijriEndPrediction,
  showWeton,
  wetonText,
  loading,
  fetchLocationAndHijriDate,
} = useHijriDate(selectedMethod, userTimezone, props.API_BASE_URL);

// styling loading
const loadingClass = computed(() =>
  darkMode.value ? "text-yellow-400" : "text-yellow-800"
);

// fetch saat mounted
fetchLocationAndHijriDate();

// watch props berubah
watch(
  () => props.selectedMethod,
  (newVal) => {
    selectedMethod.value = newVal;
    fetchLocationAndHijriDate();
  }
);

watch(
  () => props.userTimezone,
  (newVal) => {
    userTimezone.value = newVal;
    fetchLocationAndHijriDate();
  }
);
</script>
