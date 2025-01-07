import { defineStore } from 'pinia'
import { ref, computed } from 'vue';
import { useWebSocketStore } from './WebsocketStore';

export const useSessionsStore = defineStore('sessions', () => {
  const username = ref("")
  const gameState = ref("init")
  const uuid = ref("")
  const websocketStore = useWebSocketStore();

  const isLoggedIn = computed(() => {
	return username.value !== "" && username.value && gameState.value !== "init"
  })
  const login = (inputUsername) => {
    username.value = inputUsername
    gameState.value = "loggedIn"
    websocketStore.sendMessage({
      type: "login",
      username: inputUsername
    });
  }

  return {username, gameState, isLoggedIn, login}
})
