import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vueHeadful from "vue-headful";
import BootstrapVue from "bootstrap-vue";
import Octicon from "vue-octicon/components/Octicon.vue";
import "vue-octicon/icons";
import axios from "axios";

import "@/assets/app.scss";

Vue.use(BootstrapVue);
Vue.component("octicon", Octicon);
Vue.component("vue-headful", vueHeadful);

axios.defaults.baseURL = "https://chorelist.gigatar.net";

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
