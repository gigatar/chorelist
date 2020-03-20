<template>
  <b-navbar toggleable="lg" type="light" variant="light" fixed="top">
    <b-navbar-brand href="#">
      <img src="@/assets/logo.png" alt="Task Tracker" height="28px" />
    </b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item href="#/dashboard" active>Dashboard</b-nav-item>
        <b-nav-item href="#">Chores</b-nav-item>
        <b-nav-item href="#">Family</b-nav-item>
      </b-navbar-nav>

      <!-- Right aligned -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown right block>
          <template slot="button-content">
            <!-- Add avatar here and no-caret to dropdown -->
            {{ userName }}
          </template>
          <b-dropdown-item @click="userModal">Profile</b-dropdown-item>
          <b-dropdown-item @click="logout">Logout</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
    <UserModal ref="child"></UserModal>
  </b-navbar>
</template>
<script>
import UserModal from "@/components/UserModal";
export default {
  name: "Navigation",
  components: {
    UserModal
  },
  beforeMount() {
    this.checkLogin();
  },
  mounted() {
    this.tokenValidationInterval();
  },
  computed: {
    userName() {
      return this.$store.getters.getUserName;
    }
  },
  methods: {
    /**
     * Run checkLogin() every 30s to ensure
     * our token hasn't expired.
     */
    tokenValidationInterval() {
      setInterval(() => {
        this.$store.dispatch("fetchToken");
        if (!this.$store.getters.getTokenValid) {
          this.$store.dispatch("logout");
        }
      }, 30000);
    },
    /**
     * Fetch token and check to make sure it's valid,
     * logout if not.
     */
    checkLogin() {
      this.$store.dispatch("fetchToken");
      if (!this.$store.getters.getTokenValid) {
        this.$store.dispatch("logout");
      }
    },
    /**
     * Logout of the application.
     */
    logout() {
      this.$store.dispatch("logout");
    },
    /**
     * Launch the user profile modal.
     */
    userModal() {
      this.$refs.child.createModal();
    }
  }
};
</script>
