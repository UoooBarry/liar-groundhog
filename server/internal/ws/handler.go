package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}


// Define a global upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

// WebSocket handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Loop to read and write messages
	for {
		// Read a message from the WebSocket connection
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		// Log the received message
		log.Printf("Received: %+v", msg)

		// Respond with a JSON message
		response := Message{
			Type:    "response",
			Content: "Message received: " + msg.Content,
		}
		if err := conn.WriteJSON(response); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
