<template>
  <div class="field">
    <label :for="name" class="label">{{ label }}</label>
    <div class="control">
      <input
      :name="name"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :min="min"
      :max="max"
      :value="modelValue"
      :autocomplete="name + '-new'"
      @input="$emit('update:modelValue', $event.target.value)"
      class="input"
      :class="feedbackClass"
    />
    </div>
    <div v-if="help != ''" class="help" :class="feedbackClass">{{ help }}</div>
  </div>
</template>

<script>
export default {
  name: "TextInput",
  props: {
    name: String,
    type: String,
    label: String,
    placeholder: String,
    required: String,
    min: String,
    max: String,
    modelValue: String,
    help: String,
    feedbackClass: String,
  },
};
</script>
