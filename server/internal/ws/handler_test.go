package ws_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"uooobarry/liar-groundhog/internal/types"
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
	loginMessage := types.Message{
		Type:     "login",
		Username: "test_user",
	}
	err = client.WriteJSON(loginMessage)
	assert.NoError(t, err, "Failed to send login message")

	// Read the response
	var response types.Message
	err = client.ReadJSON(&response)
	assert.NoError(t, err, "Failed to read login response")
	assert.Equal(t, "login", response.Type, "Unexpected message type")
	assert.Equal(t, "test_user", response.Username, "Unexpected username")
	assert.Equal(t, "Login successful", response.Content, "Unexpected content")

	// Test Room Creation
	roomCreateMessage := types.Message{
		Type:        "room_create",
		SessionUUID: response.SessionUUID, // Use the session from the login response
	}
	err = client.WriteJSON(roomCreateMessage)
	assert.NoError(t, err, "Failed to send room creation message")

	err = client.ReadJSON(&response)
	assert.NoError(t, err, "Failed to read room creation response")
	assert.Equal(t, "room_create", response.Type, "Unexpected message type")
	assert.NotEmpty(t, response.RoomUUID, "RoomUUID should not be empty")
	assert.Equal(t, "Room create successful", response.Content, "Unexpected content")

	// Test Room Join
	// Test Login Message
	loginTwo := types.Message{
		Type:     "login",
		Username: "test_user_2",
	}
	err = client.WriteJSON(loginTwo)
	assert.NoError(t, err)
	var loginTwoResponse types.Message
    var roomJoinResponse types.Message
	err = client.ReadJSON(&loginTwoResponse)
	assert.NoError(t, err)

	roomJoinMessage := types.Message{
		Type:        "room_join",
		SessionUUID: loginTwoResponse.SessionUUID,
		RoomUUID:    response.RoomUUID, // Use the created room UUID
	}
	err = client2.WriteJSON(roomJoinMessage)
	assert.NoError(t, err)

	err = client2.ReadJSON(&roomJoinResponse)
	assert.NoError(t, err)
	fmt.Printf("rs %v", roomJoinResponse)
	assert.Equal(t, "room_join", roomJoinResponse.Type)
	assert.Equal(t, "Room join successful", roomJoinResponse.Content)
}
