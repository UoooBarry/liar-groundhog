package session

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Session represents a user's session
type Session struct {
	Username string
	RoomUUID string
}

var sessions = struct {
	sync.Mutex
	data map[string]Session // uuid -> Session
}{
	data: make(map[string]Session),
}

// CreateSession generates a new session for a username
func CreateSession(username string) string {
	sessions.Lock()
	defer sessions.Unlock()

	uuid := uuid.NewString()
	sessions.data[uuid] = Session{Username: username}
	log.Printf("Created session for user '%s' with UUID '%s'", username, uuid)
	return uuid
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
