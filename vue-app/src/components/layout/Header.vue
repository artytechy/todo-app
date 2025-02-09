<template>
  <nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <RouterLink to="/" class="navbar-item">
        <img src="@/assets/pngwing.com.png" alt="logo" />
      </RouterLink>

      <a
        role="button"
        class="navbar-burger"
        aria-label="menu"
        aria-expanded="false"
        data-target="navbarBasicExample"
      >
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>

    <div id="navbarBasicExample" class="navbar-menu">
      <div class="navbar-start">
        <RouterLink
          to="/"
          class="navbar-item"
          :class="{ 'is-active': $route.path === '/' }"
          >Home</RouterLink
        >
        <template v-if="userStore.token !== ''">
          <RouterLink
            to="/todo"
            class="navbar-item"
            :class="{ 'is-active': $route.path === '/todo' }"
            >Todo</RouterLink
          >
        </template>
      </div>

      <div class="navbar-end">
        <template v-if="userStore.token !== ''">
          <div class="navbar-item has-dropdown is-hoverable">
            <a href="javascript:void(0);" class="navbar-link">{{
              userStore.user.name
            }}</a>

            <div class="navbar-dropdown">
              <a @click="logout" href="javascript:void(0);" class="navbar-item"
                >Logout</a
              >
            </div>
          </div>
        </template>

        <template v-else>
          <div class="navbar-item">
            <div class="buttons">
              <RouterLink to="/register" class="button is-primary"
                >Register</RouterLink
              >
              <RouterLink to="/login" class="button is-light">Login</RouterLink>
            </div>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>

<script>
import { onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useApi } from "@/composables/api";
import { useUserStore } from "@/stores/user.js";
import { useCookies } from "@/composables/cookies.js";
import { useRoute } from "vue-router";
import router from "@/router";

export default {
  name: "Header",
  props: {},
  emits: ["notification"],

  setup(props, ctx) {
    const userStore = useUserStore();
    const cookies = useCookies();
    const api = useApi();
    const { user, token } = storeToRefs(userStore);
    const route = useRoute();

    async function logout() {
      try {
        const response = await api.post(
          `${import.meta.env.VITE_API_URL}/users/logout`,
          { token: token.value }
        );

        if (response.error) {
          ctx.emit("notification", {
            message: response.message || "Logout failed.",
            type: "danger",
          });
          return;
        }

        // Clear user data and token
        user.value = {};
        token.value = "";

        cookies.deleteCookie();

        ctx.emit("notification", {
          message: response.message || "Logout successful!",
          type: "success",
        });

        router.push("/login");
      } catch (error) {
        console.error("Logout error:", error);
        ctx.emit("notification", {
          message: "An unexpected error occurred. Please try again.",
          type: "danger",
        });
      }
    }

    onMounted(() => {
      document.addEventListener("DOMContentLoaded", () => {
        // Get all "navbar-burger" elements
        const $navbarBurgers = Array.prototype.slice.call(
          document.querySelectorAll(".navbar-burger"),
          0
        );

        // Add a click event on each of them
        $navbarBurgers.forEach((el) => {
          el.addEventListener("click", () => {
            // Get the target from the "data-target" attribute
            const target = el.dataset.target;
            const $target = document.getElementById(target);

            // Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
            el.classList.toggle("is-active");
            $target.classList.toggle("is-active");
          });
        });
      });
    });

    return {
      userStore,
      logout,
      route,
    };
  },
};
</script>

<style scoped>
.logo {
  max-width: 100px;
  height: auto;
}
</style>
