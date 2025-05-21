import { createRouter, createWebHistory } from "vue-router";
import Home from "../pages/Home.vue";
import Calendar from "../pages/Calendar.vue";
import MoonInfo from "../pages/MoonInfo.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      title: "Kalender Hijriyah Hari Ini",
    },
  },
  {
    path: "/calendar",
    name: "Calendar",
    component: Calendar,
    meta: {
      title: "Kalender Hijriyah Bulanan",
    },
  },
  {
    path: "/moon-info",
    name: "MoonInfo",
    component: MoonInfo,
    meta: {
      title: "Informasi Fase Bulan",
    },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.afterEach((to) => {
  const defaultTitle = "Kalender Hijriyah";
  document.title = to.meta.title || defaultTitle;
});

export default router;
