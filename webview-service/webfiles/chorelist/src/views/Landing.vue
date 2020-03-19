<template>
  <b-container fluid>
    <b-navbar toggleable="lg" fixed>
      <b-navbar-brand href="#">
        <img src="@/assets/logo.png" alt="ChoreList" height="28px" />
        ChoreList
      </b-navbar-brand>
      <b-navbar-nav class="ml-auto">
        <b-navbar-nav>
          <b-nav-item href="#/login">
            <b-button>Login</b-button>
          </b-nav-item>
          <b-nav-item href="#/signup">
            <b-button variant="primary">Signup</b-button>
          </b-nav-item>
        </b-navbar-nav>
      </b-navbar-nav>
    </b-navbar>

    <b-row>
      <b-col class="bg-primary">
        <b-jumbotron
          bg-variant="primary"
          text-variant="light+"
          fluid
          container-fluid
          header="ChoreList"
          header-level="4"
          lead="I'm a Chore List"
        ></b-jumbotron>
      </b-col>
    </b-row>

    <b-row>
      <b-col>
        Other Content that's formatted nicely.
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
export default {
  name: "Landing",
  beforeMount() {
    // // If we have a valid token, redirect to the tasks page.
    // this.$store.dispatch("fetchToken");
    // if (this.$store.getters.getTokenValid) {
    //   this.$router.replace("/tasks");
    // }
  }
};
</script>
