<template>
  <div v-if="loading" class="flex items-center justify-center mt-4">
    <span :class="['loading loading-spinner mr-2', loadingClass]"></span>
    <span :class="['font-semibold', loadingClass]">📍 Menunggu lokasi...</span>
  </div>

  <HijriInfoCard
    :hijriDateText="hijriDateText"
    :darkMode="darkMode"
    :showWeton="showWeton"
    :wetonText="wetonText"
    :weekdayText="weekdayText"
    :loading="loading"
  />
  <HijriPrediction
    :hijriEndPrediction="hijriEndPrediction"
    :darkMode="darkMode"
  />
</template>

<script setup>
import { ref, computed, watch } from "vue";
import HijriInfoCard from "./HijriInfoCard.vue";
import HijriPrediction from "./HijriPrediction.vue";
import { useHijriDate } from "../composables/useHijriDate";

const dayMap = {
  Minggu: "Ahad",
  Senin: "Senin",
  Selasa: "Selasa",
  Rabu: "Rabu",
  Kamis: "Kamis",
  Jumat: "Jumat",
  Sabtu: "Sabtu",
};

const now = new Date();
const day = now.toLocaleDateString("id-ID", { weekday: "long" });
const weekdayText = dayMap[day] || day;

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
