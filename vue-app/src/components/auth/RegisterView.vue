<template>
  <div class="box mt-3">
    <p class="is-size-4">Register</p>
    <hr>
    <FormTag @myevent="submitHandler" name="myform" event="myevent">
      <TextInput
        v-model="name"
        label="Name"
        name="name"
        type="text"
        :help="errors.name ? errors.name : ''"
        :feedbackClass="errors.name ? 'is-danger' : ''"
      />
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
      <TextInput
        v-model="confirm_password"
        label="Confirm password"
        name="confirm_password"
        type="password"
        :help="errors.confirm_password ? errors.confirm_password : ''"
        :feedbackClass="errors.confirm_password ? 'is-danger' : ''"
      />
      <input type="submit" class="button is-primary" value="Register" />
    </FormTag>
  </div>
</template>

<script>
import { ref } from "vue";
import { useApi } from "@/composables/api";
import { useValidation } from "@/composables/validation";
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
    const errors = ref({});
    const name = ref("");
    const email = ref("");
    const password = ref("");
    const confirm_password = ref("");

    async function submitHandler() {
  // Prepare the payload
  const payload = {
    name: name.value,
    email: email.value,
    password: password.value,
    confirm_password: confirm_password.value,
  };

  // Validate the payload
  errors.value = validation.validate(payload);

  // Ensure passwords match
  if (password.value !== "" && confirm_password.value !== "") {
    const passwordError = validation.validatePasswordMatch(password.value, confirm_password.value).confirm_password;
    if (passwordError) {
      errors.value.confirm_password = passwordError;
    }
  }

  // If there are validation errors, notify and stop execution
  if (Object.keys(errors.value).length > 0) {
    ctx.emit("notification", {
      message: "There was an issue with the validation process.",
      type: "danger",
    });
    return;
  }

  try {
    // Send registration request
    const response = await api.post(`${import.meta.env.VITE_API_URL}/users/register`, payload);

    // Handle server errors
    if (response.error) {
      if (response.data?.errors) {
        errors.value = response.data.errors;
      }
      ctx.emit("notification", {
        message: response.message || "Registration failed.",
        type: "danger",
      });
      return;
    }

    // Success message
    ctx.emit("notification", {
      message: response.message,
      type: "success",
    });

    // Redirect to login
    router.push("/login");
  } catch (error) {
    console.error("Registration error:", error);
    ctx.emit("notification", {
      message: "An unexpected error occurred. Please try again.",
      type: "danger",
    });
  }
}


    return {
      errors,
      name,
      email,
      password,
      confirm_password,
      submitHandler,
    };
  },
};
</script>
