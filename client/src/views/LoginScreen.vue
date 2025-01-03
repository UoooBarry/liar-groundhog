<template>
  <div class="login-container">
    <div class="input-group">
      <PixelatedInput placeholder="username" v-model="username"></PixelatedInput>
      <PixelatedButton class="login-btn" @click="login()">Login</PixelatedButton>
    </div>
    <div>{{ sessionStore.username }}</div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import PixelatedButton from '../components/PixelatedButton.vue';
import PixelatedInput from '../components/PixelatedInput.vue';
import { useSessionsStore } from '../stores/SessionStores';

const sessionStore = useSessionsStore();
const username = ref("");
const login = () => {
  sessionStore.login(username.value);
}
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  justify-content: center;
}
.login-btn{
  margin-top: 30px;
}
</style>
