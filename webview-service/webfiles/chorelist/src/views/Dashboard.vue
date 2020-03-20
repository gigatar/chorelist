<template>
  <b-container fluid>
    <Navigation />
    Welcome to the secure area {{ name }}
  </b-container>
</template>

<script>
import Navigation from "@/components/Navigation";
export default {
  name: "Dashboard",
  components: {
    Navigation
  },
  computed: {
    name() {
      return this.$store.getters.getUserName;
    }
  }
};
</script>
