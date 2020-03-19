import axios from "axios";

export default {
  actions: {
    /**
     * Signup a new user.
     *
     * @param {*} context
     * @param {object} signupPayload
     *
     * @return {Promise}
     */
    signup(context, signupPayload) {
      return new Promise((resolve, reject) => {
        axios
          .post("/rest/v1/signups", signupPayload)
          .then(({ status }) => {
            resolve(status);
          })
          .catch(function(error) {
            if (typeof error.response !== "undefined") {
              reject(error.response.status);
            } else {
              reject(500);
            }
          });
      });
    },

    /**
     * Confirm a new users signup code.
     *
     * @param {*} context
     * @param {string} code
     *
     * @return {Promise}
     */
    signupConfirm(context, code) {
      return new Promise((resolve, reject) => {
        axios
          .get("/rest/v1/signups/" + code)
          .then(({ status }) => {
            resolve(status);
          })
          .catch(function(error) {
            if (typeof error.response !== "undefined") {
              reject(error.response.status);
            } else {
              reject(500);
            }
          });
      });
    }
  }
};
