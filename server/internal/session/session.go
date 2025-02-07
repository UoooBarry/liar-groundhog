package session

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"uooobarry/liar-groundhog/internal/errors"
)

// Player represents a user's session
type Player struct {
	SessionUUID string
	Username    string
	RoomUUID    string
	Conn        *websocket.Conn
	Alive       bool
}

var sessions = struct {
	sync.Mutex
	data map[string]*Player // uuid -> Session
}{
	data: make(map[string]*Player),
}

func FindSession(uuid string) (*Player, bool) {
	session, exist := sessions.data[uuid]
	if !exist {
		return nil, false
	}
	return session, exist
}

// CreateSession generates a new session for a username
func CreateSession(conn *websocket.Conn, username string) Player {
	sessions.Lock()
	defer sessions.Unlock()

	uuid := uuid.NewString()
	session := Player{
		Username:    username,
		SessionUUID: uuid,
		Conn:        conn,
        Alive:       true,
	}
	sessions.data[uuid] = &session
	log.Printf("Created session for user '%s' with UUID '%s'", username, uuid)
	return session
}

// RemoveSession deletes a session by UUID
func RemoveSession(uuid string) error {
	sessions.Lock()
	defer sessions.Unlock()

	if uuid == "" {
		return errors.NewLoggableError("Trying to remove a empty UUID session", errors.WARN)
	}

	session, exists := sessions.data[uuid]
	if !exists {
		return errors.NewLoggableError(fmt.Sprintf("session '%s' does not exist", uuid), errors.ERROR)
	}

	delete(sessions.data, uuid)
	log.Printf("Session for user '%s' has been removed", session.Username)
	return nil
}
