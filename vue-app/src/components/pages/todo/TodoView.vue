<template>
  <nav class="panel mt-3 has-background-white">
    <p class="panel-heading">Todo</p>
    <FormTag @myevent="submitHandler" name="myform" event="myevent">
      <div class="panel-block">
        <div class="control field has-addons">
          <p class="control">
            <span class="select" :class="errors.priority_id ? 'is-danger' : ''">
              <select v-model="priorityID">
                <option :value="null" disabled>Select</option>
                <option
                  v-for="priority in priorities"
                  :value="priority.id"
                  :key="priority.id"
                >
                  {{ priority.name }}
                </option>
              </select>
              <span v-if="errors.priority_id" class="help is-danger text-wrap">
                {{ errors.priority_id }}
              </span>
            </span>
          </p>
          <p class="control is-expanded">
            <input
              v-model="text"
              name="text"
              type="text"
              placeholder="Task"
              class="input"
              :class="errors.text ? 'is-danger' : ''"
            />
            <span v-if="errors.text" class="help is-danger">
              {{ errors.text }}
            </span>
          </p>
          <template v-if="todoID != 0">
            <p class="control">
              <input
                @click="cancel"
                type="button"
                class="button is-danger"
                value="Cancel"
              />
            </p>
          </template>
          <p class="control">
            <input type="submit" class="button" value="Save" />
          </p>
        </div>
      </div>
    </FormTag>

    <p class="panel-tabs">
      <a
        @click="setFilter(0)"
        :class="currentFilter === 0 ? 'is-active' : ''"
        href="javascript:void(0);"
        >All</a
      >
      <a
        v-for="priority in priorities"
        @click="setFilter(priority.id)"
        :key="priority.id"
        :class="currentFilter === priority.id ? 'is-active' : ''"
        href="javascript:void(0);"
        >{{ priority.name }}</a
      >
    </p>
    <template v-if="Object.keys(todos).length">
      <span v-for="todo in todos" :key="todo.id">
        <template
          v-if="todo.priority_id === currentFilter || currentFilter === 0"
        >
          <a
            href="javascript:void(0);"
            class="panel-block is-justify-content-space-between"
          >
            <span
              ><button
                @click="handleDelete(todo.id)"
                class="button is-danger mr-1"
              >
                <i class="fa-solid fa-trash"></i>
              </button>
              <button @click="editTodo(todo)" class="button is-warning mr-5">
                <i class="fa-solid fa-pen"></i>
              </button>
              <span>{{ todo.text }}</span></span
            >
            <span class="tag" :class="todo.priority.badge">{{
              todo.priority.name
            }}</span>
          </a>
        </template>
      </span>
    </template>
    <template v-else>
      <span class="panel-block"> No todos found. </span>
    </template>
    <div class="panel-block">
      <button class="button is-link is-outlined is-fullwidth">
        Reset all filters
      </button>
    </div>
  </nav>
</template>

<script>
import { ref, onBeforeMount } from "vue";
import { useApi } from "@/composables/api";
import { useValidation } from "@/composables/validation";
import FormTag from "@/components/forms/FormTag.vue";
import TextInput from "@/components/forms/TextInput.vue";

export default {
  name: "Todo",
  props: {},
  emits: ["notification", "forceUpdate"],
  components: {
    FormTag,
    TextInput,
  },

  setup(props, ctx) {
    const api = useApi();
    const validation = useValidation();
    const priorities = ref([]);
    const todos = ref([]);
    const errors = ref({});
    const currentFilter = ref(0);
    const todoID = ref(0);
    const priorityID = ref(null);
    const text = ref("");

    async function submitHandler() {
      const payload = {
        id: todoID.value,
        priority_id: priorityID.value,
        text: text.value,
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
          `${import.meta.env.VITE_API_URL}/todo/save`,
          payload
        );

        // Handle server errors
        if (response.error) {
          if (response.data?.errors) {
            errors.value = response.data.errors;
          }
          ctx.emit("notification", {
            message: response.message,
            type: "danger",
          });
          return;
        }

        // Clear input fields after successful creation
        priorityID.value = null;
        text.value = "";

        // Success message
        ctx.emit("notification", {
          message: response.message,
          type: "success",
        });
        ctx.emit("forceUpdate");
      } catch (error) {
        console.error("Login error:", error);
        ctx.emit("notification", {
          message: "An unexpected error occurred. Please try again.",
          type: "danger",
        });
      }
    }

    function editTodo(todo) {
      todoID.value = todo.id;
      priorityID.value = todo.priority_id;
      text.value = todo.text;
    }

    function handleDelete(id) {
      window.$Confirm({
        message: "Are you sure?",
        button: {
          no: "No",
          yes: "Yes",
        },
        /**
         * Callback Function
         * @param {Boolean} confirm
         */
        callback: async (confirm) => {
          if (confirm) {
            const payload = {
              id: id,
            };
            try {
              const response = await api.post(
                `${import.meta.env.VITE_API_URL}/todo/delete`,
                payload
              );
              // Handle errors in the response
              if (response.error) {
                ctx.emit("notification", {
                  message: response.message,
                  type: "danger",
                });
                return;
              }

              ctx.emit("notification", {
                message: response.message,
                type: "success",
              });
              ctx.emit("forceUpdate");
            } catch (error) {
              ctx.emit("notification", {
                message: "An unexpected error occurred. Please try again.",
                type: "danger",
              });
            }
          }
        },
      });
    }

    function cancel() {
      todoID.value = 0;
      priorityID.value = null;
      text.value = "";
    }

    function setFilter(filter) {
      currentFilter.value = filter;
    }

    onBeforeMount(async () => {
      try {
        const response = await api.get(
          `${import.meta.env.VITE_API_URL}/priorities`
        );
        // Handle errors in the response
        if (response.error) {
          ctx.emit("notification", {
            message: response.message,
            type: "danger",
          });
          return;
        }

        // If priorities data exists in the response, update the priorities state
        if (response.data?.priorities) {
          priorities.value = response.data.priorities;
        }
      } catch (error) {
        ctx.emit("notification", {
          message: "An unexpected error occurred. Please try again.",
          type: "danger",
        });
      }

      try {
        const response = await api.get(`${import.meta.env.VITE_API_URL}/todo`);
        // Handle errors in the response
        if (response.error) {
          ctx.emit("notification", {
            message: response.message,
            type: "danger",
          });
          return;
        }

        // If priorities data exists in the response, update the priorities state
        if (response.data?.todos) {
          todos.value = response.data.todos;
        }
      } catch (error) {
        ctx.emit("notification", {
          message: "An unexpected error occurred. Please try again.",
          type: "danger",
        });
      }
    });

    return {
      priorities,
      todos,
      errors,
      currentFilter,
      todoID,
      priorityID,
      text,
      submitHandler,
      editTodo,
      handleDelete,
      cancel,
      setFilter,
    };
  },
};
</script>
