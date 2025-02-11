package session

import (
	"fmt"
	"log"
	"sync"
	"uooobarry/liar-groundhog/internal/errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionManger struct {
	sync.Mutex
	players map[string]*Player // uuid -> Player
}

var sessionManger = &SessionManger{
	players: make(map[string]*Player),
}

func (sessionManger *SessionManger) CreateSession(conn *websocket.Conn, username string) *Player {
	sessionManger.Lock()
	defer sessionManger.Unlock()

	uuid := uuid.NewString()
	session := Player{
		Username:    username,
		SessionUUID: uuid,
		Conn:        conn,
		Alive:       true,
	}
	sessionManger.players[uuid] = &session
	log.Printf("Created session for user '%s' with UUID '%s'", username, uuid)
	return &session
}

func (SessionManger *SessionManger) FindSession(uuid string) (*Player, bool) {
	sessionManger.Lock()
	defer sessionManger.Unlock()

	session, exist := sessionManger.players[uuid]
	if !exist {
		return nil, false
	}
	return session, exist
}

func (SessionManger *SessionManger) RemoveSession(uuid string) error {
	sessionManger.Lock()
	defer sessionManger.Unlock()

	if uuid == "" {
		return errors.NewLoggableError("Trying to remove a empty UUID session", errors.WARN)
	}

	session, exists := sessionManger.players[uuid]
	if !exists {
		return errors.NewLoggableError(fmt.Sprintf("session '%s' does not exist", uuid), errors.ERROR)
	}

	delete(sessionManger.players, uuid)
	log.Printf("Session for user '%s' has been removed", session.Username)

	return nil
}
