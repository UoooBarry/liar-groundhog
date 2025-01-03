import { defineStore } from 'pinia'
import { ref } from 'vue';

export const useSessionsStore = defineStore('sessions', () => {
  const username = ref("")
  const gameState = ref("init")

  const isLoggedIn = computed(() => username !== "" && username && gameState !== "init")
  const login = (inputUsername) => {
    username = inputUsername
    gameState = "loggedIn"
  }

  return {username, gameState, isLoggedIn, login}
})