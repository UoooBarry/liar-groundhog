package ws

import (
	"log"
	"net/http"
	"sync"

	"uooobarry/liar-groundhog/internal/session"

	"github.com/gorilla/websocket"
)

// Define a global upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

var rooms = struct {
	sync.Mutex
	data map[string][]string // roomUUID -> list of usernames
}{
	data: make(map[string][]string),
}

// handleLogin handles the login request
func handleLogin(conn *websocket.Conn, msg Message) string {
	if msg.Username == "" {
		sendError(conn, "Username is required for login")
		return ""
	}

	userUUID := session.CreateSession(msg.Username)
	
	// Send response to the client
	response := Message{
		Type:    "login",
		Username: msg.Username,
		UUID:    userUUID,
		Content: "Login successful",
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Println("Write error:", err)
	}

	return userUUID
}

// WebSocket handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	var userUUID string
	defer conn.Close()
	defer func() {
		error := session.RemoveSession(userUUID)
		if error != nil {
			sendError(conn, error.Error())
		}
		log.Println(userUUID)
	}()

	// Loop to read and write messages
	for {
		var msg Message
		// Read JSON message from client
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		switch msg.Type {
		case "login":
			userUUID = handleLogin(conn, msg)
		default:
			sendError(conn, "Unknown message type")
		}
	}
}

func sendError(conn *websocket.Conn, errMsg string) {
	response := Message{
		Type:    "error",
		Content: errMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Println("Write error:", err)
	}
}
