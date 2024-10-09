import "vue-d3-network-graph/dist/style.css";

import { createApp } from 'vue'
import App from './App.vue'
import plugin from "vue-d3-network-graph";

const app = createApp(App)
app.use(plugin);
app.mount('#app');
