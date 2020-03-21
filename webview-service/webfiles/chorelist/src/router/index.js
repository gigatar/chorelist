import Vue from "vue";
import VueRouter from "vue-router";
import store from "@/store/index.js";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Landing",
    component: () => import("@/views/Landing")
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login")
  },
  {
    path: "/signup",
    name: "Signup",
    component: () => import("@/views/Signup")
  },
  {
    path: "/signup/confirm",
    name: "SignupConfirm",
    component: () => import("@/views/SignupConfirm")
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/Dashboard"),
    beforeEnter: requireAuth
  },
  {
    path: "/family",
    name: "Family",
    component: () => import("@/views/Family"),
    beforeEnter: requireAuth
  },
  {
    path: "*",
    redirect: "/"
  }
  // {
  //   path: "/about",
  //   name: "About",
  //   // route level code-splitting
  //   // this generates a separate chunk (about.[hash].js) for this route
  //   // which is lazy-loaded when the route is visited.
  //   component: () =>
  //     import(/* webpackChunkName: "about" */ "../views/About.vue")
  // }
];

const router = new VueRouter({
  routes
});

export default router;

/**
 * Middleware for route authentication requirements.
 *
 * @param {string} to
 * @param {string} from
 * @param {string} next
 */
function requireAuth(to, from, next) {
  store.dispatch("fetchToken");
  if (!store.state.user.accessToken && to.fullPath !== "/") {
    next("/");
  } else {
    if (to.path === "/" && store.state.user.accessToken) {
      next("/tasks");
    } else {
      next();
    }
  }
}
