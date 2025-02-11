package session

import (
	"github.com/gorilla/websocket"
)

// Player represents a user's session
type Player struct {
	SessionUUID string
	Username    string
	RoomUUID    string
	Conn        *websocket.Conn
	Alive       bool
}

func FindSession(uuid string) (*Player, bool) {
	session, exist := sessionManger.FindSession(uuid)
	return session, exist
}

// CreateSession generates a new session for a username
func CreateSession(conn *websocket.Conn, username string) *Player {
	session := sessionManger.CreateSession(conn, username)
	return session
}

// RemoveSession deletes a session by UUID
func RemoveSession(uuid string) error {
	err := sessionManger.RemoveSession(uuid)
	return err
}
