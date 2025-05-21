<template>
  <div
    :data-theme="theme.darkMode ? 'dark' : 'light'"
    class="w-full sticky top-0 z-50 backdrop-blur bg-base-200 bg-opacity-90 shadow"
  >
    <div class="navbar bg-base-200 shadow-md px-4">
      <!-- Logo -->
      <div class="flex-1">
        <RouterLink to="/" class="btn btn-ghost normal-case text-xl"
          >🕌 Kalender Hijriyah</RouterLink
        >
      </div>

      <!-- Menu Desktop -->
      <div class="hidden lg:flex">
        <ul class="menu menu-horizontal px-1">
          <li>
            <RouterLink
              to="/"
              :class="{ 'text-primary font-bold': $route.path === '/' }"
              >Beranda</RouterLink
            >
          </li>
          <li>
            <RouterLink
              to="/calendar"
              :class="{
                'text-primary font-bold': $route.path.startsWith('/calendar'),
              }"
              >Kalender</RouterLink
            >
          </li>
          <li>
            <RouterLink
              to="/moon-info"
              :class="{
                'text-primary font-bold': $route.path.startsWith('/moon-info'),
              }"
              >Informasi Bulan</RouterLink
            >
          </li>
        </ul>
      </div>

      <!-- Menu Mobile -->
      <div class="dropdown dropdown-end lg:hidden">
        <div tabindex="0" role="button" class="btn btn-ghost btn-circle">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </div>
        <ul
          tabindex="0"
          class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52"
        >
          <li>
            <RouterLink
              to="/"
              :class="{ 'text-primary font-bold': $route.path === '/' }"
              >Beranda</RouterLink
            >
          </li>
          <li>
            <RouterLink
              to="/calendar"
              :class="{
                'text-primary font-bold': $route.path.startsWith('/calendar'),
              }"
              >Kalender</RouterLink
            >
          </li>
          <li>
            <RouterLink
              to="/moon-info"
              :class="{
                'text-primary font-bold': $route.path.startsWith('/moon-info'),
              }"
              >Informasi Bulan</RouterLink
            >
          </li>
        </ul>
      </div>

      <!-- Toggle Tema -->
      <div class="ml-2">
        <button @click="theme.toggleTheme" class="btn btn-sm btn-outline">
          {{ theme.themeToggleText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink, useRoute } from "vue-router";
import { useThemeStore } from "@/stores/themeStore";

const route = useRoute();
const theme = useThemeStore();
</script>
