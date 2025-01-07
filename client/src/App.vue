<template>
    <div class="login-container">
      <LoginScreen v-if="!sessionStore.isLoggedIn" />
    </div>

</template>

<script setup>
import LoginScreen from './views/LoginScreen.vue';
import { useSessionsStore } from './stores/SessionStores';
import { useWebSocketStore } from './stores/WebsocketStore';

const sessionStore = new useSessionsStore();
const websocketStore = new useWebSocketStore();

websocketStore.connect();
</script>

<style scoped>
.login-container {
  margin: 0 auto;
  height: 500px;
  width: 700px;
  border: 2px solid rgb(99, 99, 175);
  border-radius: 20px;
  font-family: 'VT323';
  font-size: 25px;
  background-color: white;
}
</style>
