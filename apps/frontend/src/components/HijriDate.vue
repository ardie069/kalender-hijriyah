<script setup lang="ts">
import HijriInfoCard from "./HijriInfoCard.vue";
import HijriPrediction from "./HijriPrediction.vue";

import type { HijriDate, Method } from "@/types/hijri";

defineProps<{
  hijriDate: HijriDate | null;
  endMonthInfo: any;
  weton: string | null;
  loading: boolean;
  error: string | null;
  method: Method;
}>();
</script>

<template>
  <div>
    <!-- LOADING -->
    <div v-if="loading" class="flex items-center justify-center mt-4">
      <span class="loading loading-spinner mr-2"></span>
      <span class="font-semibold">📍 Menghitung tanggal Hijriyah…</span>
    </div>

    <!-- ERROR -->
    <div v-else-if="error" class="alert alert-error mt-4">
      {{ error }}
    </div>

    <!-- RESULT -->
    <div v-else-if="hijriDate">
      <HijriInfoCard :hijriDate="hijriDate" :weton="weton" :method="method" />

      <HijriPrediction v-if="endMonthInfo" :prediction="endMonthInfo" />
    </div>
  </div>
</template>
