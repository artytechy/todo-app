import { createApp } from "vue";
import { createPinia } from "pinia";
import { createVfm } from "vue-final-modal";
import "vue-final-modal/style.css";
import Vue3ConfirmDialog from "vue3-confirm-dialog";
import "vue3-confirm-dialog/style";

import App from "./App.vue";
import router from "@/router";

import "./assets/main.scss";

const app = createApp(App);
const pinia = createPinia();
const vfm = createVfm();
app.use(router);
app.use(pinia);
app.use(vfm);
app.use(Vue3ConfirmDialog);
app.component("vue3-confirm-dialog", Vue3ConfirmDialog.default);
window.$Confirm = app.config.globalProperties.$confirm;
app.mount("#app");
