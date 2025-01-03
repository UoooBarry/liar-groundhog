package ws

import (
	"log"
	"net/http"
	"sync"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Session struct {
	Username     string // User's username
	RoomUUID string // UUID of the game room the user is in
}

type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g., "login")
	Username string `json:"username,omitempty"` // Username for login
	UUID    string `json:"uuid,omitempty"`    // UUID generated for the user
	Content string `json:"content,omitempty"` // Additional content
}

// Define a global upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

// In-memory storage for user sessions
var sessions = struct {
	sync.Mutex
	data map[string]Session // uuid -> Session
}{
	data: make(map[string]Session),
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

	userUUID := createSession(msg.Username)
	
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

func createSession(username string) string {
	sessions.Lock()
	defer sessions.Unlock()

	uuid := uuid.NewString()
	
	session := Session{
		Username: username,
	}
	sessions.data[uuid] = session

	return uuid
}

// removeSession removes a username from the sessions map
func removeSession(uuid string) error {
	sessions.Lock()
	defer sessions.Unlock()
	
	session, exist := sessions.data[uuid]
	if !exist {
		return fmt.Errorf("Session '%s' is not exist", uuid)
	} else {
		delete(sessions.data, uuid)
		log.Printf("Session for user '%s' has been released", session.Username)
	}
	return nil
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
		error := removeSession(userUUID)
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
