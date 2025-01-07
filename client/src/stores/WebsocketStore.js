import { defineStore } from 'pinia';
import { ref, onUnmounted } from 'vue';

export const useWebSocketStore = defineStore('webSocket', () => {
  const socket = ref(null);
  const messages = ref([]); // Store incoming messages
  const isConnected = ref(false);

  const connect = () => {
    if (socket.value) {
      console.warn('WebSocket is already connected.');
      return;
    }

    socket.value = new WebSocket('ws://localhost:8080/ws');

    socket.value.addEventListener('open', (event) => {
      console.log('WebSocket connected');
      isConnected.value = true;
    });

    socket.value.addEventListener('message', (event) => {
      console.log('Message received:', event.data);
      messages.value.push(event.data);
    });

    socket.value.addEventListener('close', () => {
      console.log('WebSocket disconnected');
      isConnected.value = false;
      socket.value = null;
    });

    socket.value.addEventListener('error', (error) => {
      console.error('WebSocket error:', error);
    });
    console.log(socket.value);
  };

  const disconnect = () => {
    if (socket.value) {
      socket.value.close();
      socket.value = null;
    } else {
      console.warn('WebSocket is not connected.');
    }
  };

  const sendMessage = (message) => {
    if (socket.value && isConnected.value) {
      socket.value.send(message);
      console.log('Message sent:', message);
    } else {
      console.warn('Cannot send message: WebSocket is not connected.');
    }
  };

  // Clean up on component unmount
  onUnmounted(() => {
    disconnect();
  });

  return {
    messages,
    isConnected,
    connect,
    disconnect,
    sendMessage,
  };
});
