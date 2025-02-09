<template #default="{ close }">
  <VueFinalModal
    content-class="modal modal-card pt-5"
    @click-outside="{close}"
  >
    <!-- Modal Card -->
    <div class="modal-card">
      <!-- Modal Header -->
      <header class="modal-card-head">
        <p class="modal-card-title">{{ title }}</p>
        <button
          class="delete"
          aria-label="close"
          @click="$emit('close')"
        ></button>
      </header>

      <!-- Modal Body -->
      <section class="modal-card-body">
        <slot />
      </section>

      <!-- Modal Footer -->
      <footer class="modal-card-foot">
        <div class="buttons">
          <button class="button is-success" @click="$emit('confirm')">
            Save changes
          </button>
          <button class="button" @click="$emit('close')">Cancel</button>
        </div>
      </footer>
    </div>
  </VueFinalModal>
</template>

<script>
import { VueFinalModal } from "vue-final-modal";
export default {
  props: ["title"],
  emits: ["confirm", "close", "click-outside"],
  components: {
    VueFinalModal,
  },

  setup(props, ctx) {
    function clickOutside() {
        console.log(11)
    }

    return {
        clickOutside
    }
  } 
};
</script>
