import axios from "axios";
import router from "@/router";

export default {
  state: {
    email: null,
    userName: null,
    accessToken: null,
    accessTokenExpiration: 0,
    accessTokenTTL: 0
  },
  getters: {
    getAuthToken(state) {
      return state.accessToken;
    },
    getTokenExpiration(state) {
      return state.accessTokenExpiration;
    },
    getTokenTTL(state) {
      return state.accessTokenTTL;
    },
    getTokenValid(state) {
      if (!state.accessToken) {
        return false;
      }
      return true;
    },
    getEmail(state) {
      if (state.email === null) {
        state.email = sessionStorage.getItem("email");
      }
      return state.email;
    },
    getUserName(state) {
      if (state.userName === null) {
        state.userName = sessionStorage.getItem("userName");
      }
      return state.userName;
    }
  },
  mutations: {
    updateAccessToken: (state, data) => {
      if (!data) {
        state.accessToken = null;
        return;
      }
      state.accessToken = data;

      // Get expiration time from token
      var exipration = JSON.parse(window.atob(state.accessToken.split(".")[1]));
      state.accessTokenTTL = exipration.TTL;

      // Check the token validity before replying.
      var currentTime = Math.round(new Date().getTime() / 1000);
      if (currentTime > exipration.Claims.exp) {
        state.accessToken = null;
      }
    },
    removeAccessToken: state => {
      state.accessToken = null;
    },
    removeEmail: state => {
      state.email = null;
    },
    setEmail: (state, email) => {
      state.email = email;
    },
    removeUserName: state => {
      state.userName = null;
    },
    setUserName: (state, userName) => {
      state.userName = userName;
    }
  },
  actions: {
    /**
     * Fetch the token from the browser sessionStorage and
     * place it in the state.
     *
     * @param {*} context
     */
    fetchToken(context) {
      context.commit(
        "updateAccessToken",
        sessionStorage.getItem("accessToken")
      );
    },

    /**
     * Login to the application and set all state and browser variables.
     *
     * Note: Currently using browser session storage which means each tab
     *       will get a unique session and require login.  This should
     *       probbably be switched to the localStorage.
     *
     * @param {*} context
     * @param {object} payload JSON login information
     *
     * @return {Promise}
     */
    login(context, payload) {
      return new Promise((resolve, reject) => {
        axios
          .post("/rest/v1/users/login", payload)
          .then(({ data, status, headers }) => {
            if (status === 200) {
              sessionStorage.setItem("accessToken", headers.authorization);
              sessionStorage.setItem("email", data.email);
              sessionStorage.setItem("userName", data.name);
              context.commit("updateAccessToken", headers.authorization);
              context.commit("setUserName", data.name);
              context.commit("setEmail", data.email);
              resolve(true);
            }
          })
          .catch(function(error) {
            reject(error.response);
          });
      });
    },

    /**
     * Log out of the application by removing all variables and
     * forcing a non-history navigation to "/".
     *
     * @param {*} context
     */
    logout(context) {
      sessionStorage.removeItem("accessToken");
      sessionStorage.removeItem("email");
      sessionStorage.removeItem("displayName");
      context.commit("removeAccessToken");
      context.commit("removeEmail");
      context.commit("removeDisplayName");

      router.go("/");
    },

    changeUserName(context, newName) {
      return new Promise((resolve, reject) => {
        axios
          .patch("rest/v1/users/name", newName, {
            headers: {
              Authorization: "Bearer " + context.getters.getAuthToken
            }
          })
          .then(({ status }) => {
            sessionStorage.setItem("userName", newName.name);
            context.commit("setUserName", newName.name);
            resolve(status);
          })
          .catch(function(error) {
            reject(error.response);
          });
      });
    },

    changePassword(context, passwordChange) {
      return new Promise((resolve, reject) => {
        axios
          .patch("/rest/v1/users/password", passwordChange, {
            headers: {
              Authorization: "Bearer " + context.getters.getAuthToken
            }
          })
          .then(({ status }) => {
            resolve(status);
          })
          .catch(function(error) {
            reject(error.response);
          });
      });
    },
    getUser(context, userID) {
      return new Promise((resolve, reject) => {
        axios
          .get("/rest/v1/users/" + userID, {
            headers: {
              Authorization: "Bearer " + context.getters.getAuthToken
            }
          })
          .then(({ data }) => {
            resolve(data);
          })
          .catch(function(error) {
            reject(error.response);
          });
      });
    }
  }
};
