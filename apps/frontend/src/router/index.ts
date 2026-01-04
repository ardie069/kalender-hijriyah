import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";

import Home from "@/pages/Home.vue";
import Calendar from "@/pages/Calendar.vue";
import MoonInfo from "@/pages/MoonInfo.vue";

const routes: RouteRecordRaw[] = [
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
  const title =
    typeof to.meta.title === "string" ? to.meta.title : defaultTitle;

  document.title = title;
});

export default router;
