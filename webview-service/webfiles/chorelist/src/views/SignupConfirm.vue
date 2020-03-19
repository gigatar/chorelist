<template>
  <b-container fluid>
    <b-row align-h="center">
      <b-col lg="5">
        <b-container fluid>
          <img src="@/assets/logo.png" />
          <b-container v-if="error">
            We encounterd an error validating your signup: {{ error }}
          </b-container>
        </b-container></b-col
      ></b-row
    >
  </b-container>
</template>

<script>
export default {
  name: "SignupConfirm",
  data: () => ({
    code: null,
    error: null
  }),

  beforeMount() {
    this.getCode();
    this.confirmSignup();
  },
  methods: {
    /**
     * Get the code from the route.  Should maybe move this to a param
     * but that would require a backend change in the email-service.
     */
    getCode() {
      this.code = this.$route.query.code;
    },
    /**
     * Submit the signup confirm and then navigate to the login page.
     */
    confirmSignup() {
      this.$store
        .dispatch("signupConfirm", this.code)
        .then(() => {
          this.$router.push({ name: "Login", params: { signupSuccess: true } });
        })
        .catch(error => {
          if (error === 404) {
            this.error = "Signup code not found.";
          } else {
            this.error = error;
          }
        });
    }
  }
};
</script>
