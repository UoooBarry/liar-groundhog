import { defineStore } from 'pinia'
import { ref, computed } from 'vue';

export const useSessionsStore = defineStore('sessions', () => {
  const username = ref("")
  const gameState = ref("init")

  const isLoggedIn = computed(() => {
	return username.value !== "" && username.value && gameState.value !== "init"
  })
  const login = (inputUsername) => {
    username.value = inputUsername
    gameState.value = "loggedIn"
  }

  return {username, gameState, isLoggedIn, login}
})
