<template>
  <b-container>
    <SecureArea />
    Welcome to the secure area {{ name }}
  </b-container>
</template>

<script>
import SecureArea from "@/components/SecureArea";
export default {
  name: "Dashboard",
  components: {
    SecureArea
  },
  computed: {
    name() {
      return this.$store.getters.getUserName;
    }
  }
};
</script>
