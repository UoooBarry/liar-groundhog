package ws_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"uooobarry/liar-groundhog/internal/message"
	"uooobarry/liar-groundhog/internal/ws"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	// Create an HTTP test server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(ws.HandleWebSocket))
	defer server.Close()

	// Convert the server URL to a WebSocket URL
	wsURL := "ws:" + server.URL[len("http:"):]
	fmt.Printf("ws: %s", wsURL)

	// Dial the WebSocket server
	client, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to WebSocket server")
	defer client.Close()

	client2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err, "Failed to connect to WebSocket server")
	defer client2.Close()

	// Test Login Message
	loginMessage := message.Message{
		Type:     "login",
		Username: "test_user",
	}
	err = client.WriteJSON(loginMessage)
	assert.NoError(t, err, "Failed to send login message")

	// Read the response
	var loginMsg message.LoginSuccessMessage
	err = client.ReadJSON(&loginMsg)
	assert.NoError(t, err, "Failed to read login response")
	assert.Equal(t, "login", loginMsg.Type, "Unexpected message type")
	assert.Equal(t, "test_user", loginMsg.Username, "Unexpected username")
	assert.Equal(t, "Login successful", loginMsg.Content, "Unexpected content")

	// Test Room Creation
	roomCreateMessage := message.RoomCreateMessage{
		SessionUUID: loginMsg.SessionUUID, // Use the session from the login response
		Message: message.Message{
			Type: "room_create",
		},
	}
	err = client.WriteJSON(roomCreateMessage)
	assert.NoError(t, err, "Failed to send room creation message")

	var roomMsg message.RoomOpMessage
	err = client.ReadJSON(&roomMsg)
	assert.NoError(t, err, "Failed to read room creation response")
	resNoError(t, roomMsg.Message)
	log.Println("Room MSG:", roomMsg)
	assert.NotEmpty(t, roomMsg.RoomUUID, "RoomUUID should not be empty")

	// Test Room Join
	// Test Login Message
	loginTwo := message.LoginMessage{
		Message: message.Message{
			Type: "login",
		},
		Username: "test_user_2",
	}
	err = client.WriteJSON(loginTwo)
	assert.NoError(t, err)
	var loginTwoResponse message.LoginSuccessMessage
	var roomJoinResponse message.RoomOpMessage
	err = client.ReadJSON(&loginTwoResponse)
	resNoError(t, loginTwoResponse.Message)
	assert.NoError(t, err)

	roomJoinMessage := message.RoomOpMessage{
		Message: message.Message{
			Type: "room_join",
		},
		SessionUUID: loginTwoResponse.SessionUUID,
		RoomUUID:    roomMsg.RoomUUID, // Use the created room UUID
	}
	err = client2.WriteJSON(roomJoinMessage)
	assert.NoError(t, err)

	err = client2.ReadJSON(&roomJoinResponse)
	assert.NoError(t, err)
	fmt.Printf("rs %v", roomJoinResponse)
	resNoError(t, roomJoinMessage.Message)
	assert.Equal(t, "room_join", roomJoinResponse.Type)
	assert.Equal(t, "Room join successful", roomJoinResponse.Content)
}

func resNoError(t *testing.T, res message.Message) {
	assert.NotEqual(t, "error", res.Type)
}
