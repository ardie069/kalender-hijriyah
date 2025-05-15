<template>
  <div v-if="loading" class="flex items-center mt-4">
    <span class="loading loading-spinner text-accent mr-2"></span>
    <span class="font-semibold text-accent">📍 Menunggu lokasi...</span>
  </div>

  <HijriInfoCard :hijriDateText="hijriDateText" :darkMode="darkMode" />
  <WetonInfo :showWeton="showWeton" :wetonText="wetonText" :darkMode="darkMode" />
  <HijriPrediction
    :hijriEndPrediction="hijriEndPrediction"
    :darkMode="darkMode"
  />
</template>

<script setup>
import { ref, computed, watch } from "vue";
import HijriInfoCard from "./HijriInfoCard.vue";
import HijriPrediction from "./HijriPrediction.vue";
import WetonInfo from "./WetonInfo.vue";
import { useHijriDate } from "../composables/useHijriDate";

// props dari parent
const props = defineProps({
  selectedMethod: String,
  userTimezone: String,
  API_BASE_URL: String,
  darkMode: Boolean,
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
  props.darkMode ? "text-yellow-400" : "text-yellow-800"
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
