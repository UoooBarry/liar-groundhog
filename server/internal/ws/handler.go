package ws

import (
	"log"
	"net/http"
	"sync"

	appErrors "uooobarry/liar-groundhog/internal/errors"
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
		return appErrors.NewClientError("Username is required for login")
	}

	return nil
}

// handleLogin handles the login request
func handleLogin(conn *websocket.Conn, msg types.Message) (string, error) {
	if err := validateUsername(msg); err != nil {
		return "", err
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

	return user.SessionUUID, nil
}

// Helper function to handle errors in a consistent manner
func handleMessageError(conn *websocket.Conn, err error) {
	if err != nil {
		appErrors.HandleError(conn, err)
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
		var msg types.Message
		var err error
		// Read JSON message from client
		formatErr := conn.ReadJSON(&msg)
		if formatErr != nil {
			appErrors.HandleError(conn, appErrors.NewLoggableError("Read Error: "+formatErr.Error(), appErrors.WARN))
			break
		}

		switch msg.Type {
		case "login":
			userUUID, err = handleLogin(conn, msg)
			handleMessageError(conn, err)
		case "room_create":
			err = handleRoomCreate(conn, msg)
			handleMessageError(conn, err)
		case "room_join":
			err = handleRoomJoin(conn, msg)
			handleMessageError(conn, err)
        case "room_start":
            err = handleRoomStart(conn, msg)
            handleMessageError(conn, err)
        case "game_place_card":
            err = handleRoomPlaceCard(conn, msg)
            handleMessageError(conn, err)
		default:
			handleMessageError(conn, appErrors.NewClientError("Unknown message type"))
		}
	}
}
