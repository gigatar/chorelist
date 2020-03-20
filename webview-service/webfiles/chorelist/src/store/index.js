import Vue from "vue";
import Vuex from "vuex";
import Signup from "@/store/signup";
import User from "@/store/user";
import Family from "@/store/family";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    signup: Signup,
    user: User,
    family: Family
  }
});
