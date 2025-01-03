package session

import (
	"sync"
	"log"

    "uooobarry/liar-groundhog/internal/types"

	"github.com/google/uuid"
)

var rooms = struct {
	sync.Mutex
	data map[string]Room // uuid -> Session
}{
	data: make(map[string]Room),
}

type Room struct {
	RoomUUID string
	Players []Session
}

func CreateRoom() (string, Room) {
	rooms.Lock()
	defer rooms.Unlock()

	uuid := uuid.NewString()
	room := Room{RoomUUID: uuid}
	rooms.data[uuid] = room
	log.Printf("Created room UUID '%s'", uuid)
	return uuid, room
}

func FindRoom(uuid string) (Room, bool) {
	room, exist := rooms.data[uuid]
	return room, exist
}

func (room *Room) FindPlayerInRoom(username string) (*Session, bool) {
	for _, player := range room.Players {
		if player.Username == username {
			return &player, true
		}
	}

	return nil, false
}

func (room *Room) SendPublicPlayerAction(player Session, msg types.ActionMessage) {}
