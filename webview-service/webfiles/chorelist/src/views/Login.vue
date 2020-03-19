<template>
  <b-container fluid>
    <b-alert :variant="alert.variant" v-model="alert.show" dismissible>{{
      alert.text
    }}</b-alert>
    <b-row align-h="center">
      <b-col lg="3">
        <b-container fluid>
          <img src="@/assets/logo.png" />
          <form @submit.prevent="noop">
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="mail"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="email-input"
                  v-model="loginData.email"
                  required
                  type="text"
                  placeholder="Enter Email Address"
                  autocomplete="email"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-form-group class="mb-1">
              <b-input-group size="md">
                <b-input-group-prepend is-text>
                  <octicon name="lock"></octicon>
                </b-input-group-prepend>
                <b-form-input
                  id="password-input"
                  v-model="loginData.password"
                  required
                  type="password"
                  placeholder="Enter Password"
                  autocomplete="current-password"
                ></b-form-input>
              </b-input-group>
            </b-form-group>
            <b-button
              block
              type="submit"
              variant="primary"
              :disabled="loginDisabled"
              @click="login"
              >Login</b-button
            >
          </form>
        </b-container>
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
export default {
  name: "Login",
  data: () => ({
    loginDisabled: false,
    loginData: { email: null, password: null },
    alert: { show: false, variant: "danger", text: "Test" }
  }),
  beforeMount() {
    this.checkSignup();
  },
  methods: {
    /**
     * Check signup reads a route param and if it's set
     * displays the signup message.
     */
    checkSignup() {
      if (this.$route.params.signupSuccess === true) {
        this.alert = { show: true, variant: "success", text: "Signup Success" };
      }
    },

    /**
     * Ummm... it's a noop..
     */
    noop() {
      /*Prevent form submit and bug*/
    },

    /**
     * Login to the app.
     */
    login() {
      this.loginDisabled = true;
      this.alert.show = false;

      this.$store
        .dispatch("login", this.loginData)
        .then(() => {
          this.$router.push("/dashboard");
        })
        .catch(error => {
          if (error.status === 400) {
            this.alert = {
              show: true,
              variant: "danger",
              text: "Invalid Username and/or Password"
            };
          } else {
            this.alert = {
              show: true,
              variant: "danger",
              text: "Unknown error: " + error.data.message
            };
          }
          this.loginDisabled = false;
        });
    }
  }
};
</script>
