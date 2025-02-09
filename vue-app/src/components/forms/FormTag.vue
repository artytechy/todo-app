<template>
  <form
    @submit.prevent="submit"
    ref="myForm"
    :event="event"
    :action="action"
    autocomplete="off"
    :method="method"
    class="needs-validation"
    novalidate
  >
    <slot />
  </form>
</template>

<script>
import { ref } from "vue";

export default {
  name: "FormTag",
  props: ["action", "method", "name", "event"],

  setup(props, ctx) {
    let myForm = ref(null);

    function submit() {
      if (myForm.value.checkValidity()) {
        ctx.emit(props["event"]);
      }
      myForm.value.classList.add("was-validated");
    }

    return {
      myForm,
      submit,
    };
  },
};
</script>
