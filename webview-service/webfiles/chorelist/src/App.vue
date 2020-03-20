<template>
  <b-container fluid id="app">
    <vue-headful :title="title" />
    <router-view :key="$route.fullpath" />
  </b-container>
</template>

<script>
export default {
  computed: {
    // Set the title based on the route name.
    title: function() {
      return "ChoreList [" + this.$route.name + "]";
    }
  }
};
</script>
<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
