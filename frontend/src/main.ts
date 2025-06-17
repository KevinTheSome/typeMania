import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";

import "./style.css";
import App from "./App.vue";
import Game from "./Game.vue";
import Lb from "./Lb.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: Game },
    { path: "/lb", component: Lb },
  ],
});

const app = createApp(App);

app.use(router);
app.component("App", App);
app.component("Lb", Lb);

app.mount("#app");
