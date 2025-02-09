import { ref } from "vue";
import { defineStore } from "pinia";

export const useUserStore = defineStore("user", () => {
  const user = ref({});
  const token = ref("");

  function logout() {
    user.value = {};
    token.value = "";
  }

  return { user, token, logout };
});
