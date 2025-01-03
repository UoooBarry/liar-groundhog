<template>
  <link href='https://fonts.googleapis.com/css?family=VT323' rel='stylesheet' type='text/css'>
  <div class="container">
    <LoginScreen v-if="!sessionStore.isLoggedIn"/>
  </div>
</template>

<script setup>
import LoginScreen from './views/LoginScreen.vue';
import { useSessionsStore } from './stores/SessionStores';

const sessionStore = useSessionsStore();
</script>

<style scoped>
.container{
  margin: 0 auto;
  height: 600px;
  width: 800px;
  border: 2px solid black;
  font-family: 'VT323';
  font-size: 25px;
}
</style>
