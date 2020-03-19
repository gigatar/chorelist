<template>
  <b-container></b-container>
</template>
<script>
export default {
  name: "SecureArea",
  beforeMount() {
    this.checkLogin();
  },
  mounted() {
    this.tokenValidationInterval();
  },
  methods: {
      /**
       * Run checkLogin() every 30s to ensure
       * our token hasn't expired.
       */
    tokenValidationInterval() {
      setInterval(() => {
        this.checkLogin();
      }, 30000);
    },
    checkLogin() {
        this.$store.dispatch("fetchToken");
      if (!this.$store.getters.getTokenValid) {
        this.$store.dispatch("logout");
      }
    }
  }
};
</script>
