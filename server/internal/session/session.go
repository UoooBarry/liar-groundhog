package session

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Session represents a user's session
type Session struct {
	SessionUUID string
	Username    string
	RoomUUID    string
	Conn        *websocket.Conn
}

var sessions = struct {
	sync.Mutex
	data map[string]Session // uuid -> Session
}{
	data: make(map[string]Session),
}

func FindSession(uuid string) (Session, bool) {
	session, exist := sessions.data[uuid]
	return session, exist
}

// CreateSession generates a new session for a username
func CreateSession(conn *websocket.Conn, username string) (string, Session) {
	sessions.Lock()
	defer sessions.Unlock()

	uuid := uuid.NewString()
	session := Session{
		Username:    username,
		SessionUUID: uuid,
		Conn:        conn,
	}
	sessions.data[uuid] = session
	log.Printf("Created session for user '%s' with UUID '%s'", username, uuid)
	return uuid, session
}

// RemoveSession deletes a session by UUID
func RemoveSession(uuid string) error {
	sessions.Lock()
	defer sessions.Unlock()

	session, exists := sessions.data[uuid]
	if !exists {
		return fmt.Errorf("session '%s' does not exist", uuid)
	}

	delete(sessions.data, uuid)
	log.Printf("Session for user '%s' has been removed", session.Username)
	return nil
}
