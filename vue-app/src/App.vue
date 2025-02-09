<template>
  <vue3-confirm-dialog />
  <div class="has-background-info">
    <Header @notification="setNotification" />
    <div class="container">
      <div
        v-if="notificationMessage != '' && notificationType != ''"
        class="notification mt-3 mb-0"
        :class="`is-${notificationType}`"
      >
        <button @click="clearNotification" class="delete"></button>
        {{ notificationMessage }}
      </div>
      <RouterView
        v-slot="{ Component }"
        :key="componentKey"
        @notification="setNotification"
        @forceUpdate="forceUpdate"
        ><component :is="Component"
      /></RouterView>
    </div>
    <Footer />
  </div>
</template>

<script>
import { onBeforeMount, ref } from "vue";
import { storeToRefs } from "pinia";
import { useUserStore } from "@/stores/user.js";
import { useCookies } from "./composables/cookies.js";
import Header from "@/components/layout/Header.vue";
import Footer from "@/components/layout/Footer.vue";

export default {
  name: "App",
  components: {
    Header,
    Footer,
  },

  setup() {
    const userStore = useUserStore();
    const { user, token } = storeToRefs(userStore);
    const cookies = useCookies();
    const notificationMessage = ref("");
    const notificationType = ref("");
    const componentKey = ref(0);

    function setNotification(data) {
      notificationMessage.value = data.message;
      notificationType.value = data.type;
    }

    function clearNotification() {
      notificationMessage.value = "";
      notificationType.value = "";
    }

    function forceUpdate() {
      componentKey.value += 1;
    }

    onBeforeMount(() => {
      // Check for a cookie
      let data = cookies.getCookie("_SITE_DATA");

      if (data !== "") {
        let cookieData = JSON.parse(data);

        // update store
        token.value = cookieData.token;
        user.value = {
          id: cookieData.user.id,
          name: cookieData.user.name,
          email: cookieData.user.email,
        };
      }
    });

    return {
      notificationMessage,
      notificationType,
      componentKey,
      setNotification,
      clearNotification,
      forceUpdate,
    };
  },
};
</script>

<style scoped>
.container {
  min-height: 90vh;
}
</style>
