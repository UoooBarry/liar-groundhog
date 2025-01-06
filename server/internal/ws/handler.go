package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"uooobarry/liar-groundhog/internal/session"
	"uooobarry/liar-groundhog/internal/types"
	"uooobarry/liar-groundhog/internal/utils"

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

func validateUsername(msg types.Message) error {
	if msg.Username == "" {
		return errors.New("Username is required for login")
	}

	return nil
}

// handleLogin handles the login request
func handleLogin(conn *websocket.Conn, msg types.Message) *string {
	if err := validateUsername(msg); err != nil {
		utils.SendError(conn, err.Error())
		return nil
	}

	user := session.CreateSession(conn, msg.Username)

	// Send response to the client
	response := types.Message{
		Type:        "login",
		Username:    msg.Username,
		SessionUUID: user.SessionUUID,
		Content:     "Login successful",
	}
	utils.SendResponse(conn, response)

	return &user.SessionUUID
}

func handleRoomCreate(conn *websocket.Conn, msg types.Message) {
	room, err := session.CreateRoom(msg.SessionUUID)
	if err != nil {
		utils.SendError(conn, err.Error())
		return
	}
	response := types.Message{
		Type:     "room_create",
		RoomUUID: room.RoomUUID,
		Content:  "Room create successful",
	}
	utils.SendResponse(conn, response)
}

func hanldeRoomJoin(conn *websocket.Conn, msg types.Message) {
	room, exist := session.FindRoom(msg.RoomUUID)
	if !exist {
		utils.SendError(conn, fmt.Sprintf("Room ID '%s' does not exist", msg.RoomUUID))
		return
	}
	if err := room.AddPlayer(msg.SessionUUID); err != nil {
		utils.SendError(conn, err.Error())
		return
	}

	response := types.Message{
		Type:     "room_join",
		RoomUUID: room.RoomUUID,
		Content:  "Room join successful",
	}
	utils.SendResponse(conn, response)
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
			utils.SendError(conn, error.Error())
		}
		log.Println(userUUID)
	}()

	// Loop to read and write messages
	for {
		var msg types.Message
		// Read JSON message from client
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		switch msg.Type {
		case "login":
			userUUID = *handleLogin(conn, msg)
		case "room_create":
			handleRoomCreate(conn, msg)
		case "room_join":
			hanldeRoomJoin(conn, msg)
		default:
			utils.SendError(conn, "Unknown message type")
		}
	}
}
