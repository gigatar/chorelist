import Vue from "vue";
import Vuex from "vuex";
import Signup from "@/store/signup";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    signup: Signup
  }
});
