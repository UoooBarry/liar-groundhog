package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	appErrors "uooobarry/liar-groundhog/internal/errors"
	"uooobarry/liar-groundhog/internal/message"
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

func validateUsername(msg message.LoginMessage) error {
	if msg.Username == "" {
		return appErrors.NewClientError("Username is required for login")
	}

	return nil
}

// handleLogin handles the login request
func handleLogin(conn *websocket.Conn, msg message.LoginMessage) (string, error) {
	if err := validateUsername(msg); err != nil {
		return "", err
	}

	user := session.CreateSession(conn, msg.Username)

	// Send response to the client
	response := message.LoginSuccessMessage{
		Message: message.Message{
			Content: "Login successful",
			Type:    "login",
		},
		Username:    msg.Username,
		SessionUUID: user.SessionUUID,
	}
	message.SendResponse(conn, response)

	return user.SessionUUID, nil
}

// Helper function to handle errors in a consistent manner
func handleMessageError(conn *websocket.Conn, err error) {
	if err != nil {
		HandleError(conn, err)
	}
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
			handleMessageError(conn, error)
		}
	}()

	// Loop to read and write messages
	for {
        _, rawMsg, err := conn.ReadMessage()
        if err != nil {
            log.Println("ReadMessage Error:", err)
            return
        }

		var msg message.Message
        if err := json.Unmarshal(rawMsg, &msg); err != nil {
            log.Println("JSON Unmarshal Error:", err)
            return
        }
        msg.RawData = rawMsg
        
		switch msg.Type {
		case "login":
			var loginMsg message.LoginMessage
			if err := json.Unmarshal(msg.RawData, &loginMsg); err != nil {
				HandleError(conn, appErrors.NewLoggableError("Parse Error: "+err.Error(), appErrors.WARN))
			}

			userUUID, err = handleLogin(conn, loginMsg)
			handleMessageError(conn, err)
		case "room_create":
            var roomCreateMessage message.RoomCreateMessage
			if err := json.Unmarshal(msg.RawData, &roomCreateMessage); err != nil {
				HandleError(conn, appErrors.NewLoggableError("Parse Error: "+err.Error(), appErrors.WARN))
			}
            log.Println(roomCreateMessage)
            log.Println(roomCreateMessage.SessionUUID)

			err = handleRoomCreate(conn, roomCreateMessage)
			handleMessageError(conn, err)
		case "room_join":
			var roomJoinMsg message.RoomOpMessage
			if err := json.Unmarshal(msg.RawData, &roomJoinMsg); err != nil {
				HandleError(conn, appErrors.NewLoggableError("Parse Error: "+err.Error(), appErrors.WARN))
			}

			err = handleRoomJoin(conn, roomJoinMsg)
			handleMessageError(conn, err)
		case "room_start":
            var roomMsg message.RoomOpMessage
			if err := json.Unmarshal(msg.RawData, &roomMsg); err != nil {
				HandleError(conn, appErrors.NewLoggableError("Parse Error: "+err.Error(), appErrors.WARN))
			}

			err = handleRoomStart(conn, roomMsg)
			handleMessageError(conn, err)
		case "player_action":
            var actionMsg message.PlayerActionMessage
			if err := json.Unmarshal(msg.RawData, &actionMsg); err != nil {
				HandleError(conn, appErrors.NewLoggableError("Parse Error: "+err.Error(), appErrors.WARN))
			}
			err = handlePlayerAction(conn, actionMsg)
			handleMessageError(conn, err)
		default:
			handleMessageError(conn, appErrors.NewClientError("Unknown message type"))
		}
	}
}
