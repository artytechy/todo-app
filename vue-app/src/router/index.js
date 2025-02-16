import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/stores/user";

import LoginView from "@/components/auth/LoginView.vue";
import RegisterView from "@/components/auth/RegisterView.vue";
import TodoView from "@/components/pages/todo/TodoView.vue";

const routes = [
  {
    path: "/",
    redirect: "/todo",
  },
  {
    path: "/login",
    name: "Login",
    component: LoginView,
    meta: {
      hideForAuth: true,
    },
  },
  {
    path: "/register",
    name: "Register",
    component: RegisterView,
    meta: {
      hideForAuth: true,
    },
  },
  {
    path: "/todo",
    name: "Todo",
    component: TodoView,
    meta: {
      requiresAuth: true,
    },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // Check if the route requires authentication
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (userStore.token === "") {
      // If not logged in, redirect to login
      next({ name: "Login" });
    } else {
      // If logged in, proceed to the requested route
      next();
    }
  }
  // Check if the route is for logged-in users only (e.g., login or register)
  else if (to.matched.some((record) => record.meta.hideForAuth)) {
    if (userStore.token !== "") {
      // If logged in, redirect to login
      next({ name: "Login" });
    } else {
      // If not logged in, proceed to the requested route
      next();
    }
  } else {
    // For all other routes, proceed normally
    next();
  }
});

export default router;
