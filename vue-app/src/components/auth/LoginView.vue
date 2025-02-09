<template>
  <div class="box mt-3">
    <p class="is-size-4">Login</p>
    <hr />
    <FormTag @myevent="submitHandler" name="myform" event="myevent">
      <TextInput
        v-model="email"
        label="Email"
        name="email"
        type="email"
        :help="errors.email ? errors.email : ''"
        :feedbackClass="errors.email ? 'is-danger' : ''"
      />
      <TextInput
        v-model="password"
        label="Password"
        name="password"
        type="password"
        :help="errors.password ? errors.password : ''"
        :feedbackClass="errors.password ? 'is-danger' : ''"
      />
      <input type="submit" class="button is-primary" value="Login" />
    </FormTag>
  </div>
</template>

<script>
import { ref } from "vue";
import { storeToRefs } from "pinia";
import { useApi } from "@/composables/api";
import { useValidation } from "@/composables/validation";
import { useUserStore } from "@/stores/user.js";
import { useCookies } from "@/composables/cookies.js";
import router from "@/router";
import FormTag from "../forms/FormTag.vue";
import TextInput from "../forms/TextInput.vue";

export default {
  name: "Login",
  props: {},
  emits: ["notification"],
  components: {
    FormTag,
    TextInput,
  },

  setup(props, ctx) {
    const api = useApi();
    const validation = useValidation();
    const userStore = useUserStore();
    const { user, token } = storeToRefs(userStore);
    const cookies = useCookies();
    const errors = ref({});
    const email = ref("");
    const password = ref("");

    async function submitHandler() {
      const payload = {
        email: email.value,
        password: password.value,
      };

      // Validate the payload
      errors.value = validation.validate(payload);

      if (Object.keys(errors.value).length > 0) {
        ctx.emit("notification", {
          message: "There was an issue with the validation process.",
          type: "danger",
        });
        return;
      }

      try {
        const response = await api.post(
          `${import.meta.env.VITE_API_URL}/users/login`,
          payload
        );

        // Handle errors in the response
        if (response.error) {
          if (response.data?.errors) {
            errors.value = response.data.errors;
          }
          ctx.emit("notification", {
            message: response.message || "Login failed.",
            type: "danger",
          });
          return;
        }

        // Check for valid user & token
        if (response.data?.user && response.data?.token) {
          user.value = response.data.user;
          token.value = response.data.token;

          cookies.storeCookieWithUserData(response.data);

          ctx.emit("notification", {
            message: response.message || "Login successful!",
            type: "success",
          });

          router.push("/todo");
        }
      } catch (error) {
        console.error("Login error:", error);
        ctx.emit("notification", {
          message: "An unexpected error occurred. Please try again.",
          type: "danger",
        });
      }
    }

    return {
      errors,
      email,
      password,
      submitHandler,
    };
  },
};
</script>
