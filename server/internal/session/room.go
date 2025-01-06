package session

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/types"
	"uooobarry/liar-groundhog/internal/utils"

	"github.com/google/uuid"
)

const MAX_PLAYERS = 4

var rooms = struct {
	sync.Mutex
	data map[string]*Room // uuid -> Session
}{
	data: make(map[string]*Room),
}

type Room struct {
	RoomUUID string
	Players  []Session
	Engine   liar.Engine
}

func CreateRoom(ownerUUID string) (*Room, error) {
	rooms.Lock()
	defer rooms.Unlock()

    owner, exist := FindSession(ownerUUID)

	if !exist {
		return nil, fmt.Errorf("Player session not exist '%s'", ownerUUID)
	}
	uuid := uuid.NewString()
    room := &Room{RoomUUID: uuid, Engine: liar.New()}
	rooms.data[uuid] = room
    room.Players = append(room.Players, *owner)

	log.Printf("Created room UUID '%s'", uuid)
	return room, nil
}

func FindRoom(uuid string) (*Room, bool) {
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

func (room *Room) PublishRoomInfo() {
	for _, player := range room.Players {
		conn := player.Conn
		if conn == nil {
			continue
		}
		utils.SendResponse(conn, room.GetInfoMessage())
	}
}

func validPlayerJoin(room *Room, playerUUID string) (*Session, error) {
	player, exist := FindSession(playerUUID)

	if !exist {
		return nil, fmt.Errorf("Player session not exist '%s'", playerUUID)
	}

	if len(room.Players) >= MAX_PLAYERS {
		return player, errors.New("The current game room is full.")
	}

	if _, inRoom := room.FindPlayerInRoom(player.Username); inRoom {
		return player, fmt.Errorf("Player '%s' is already in the room", player.Username)
	}

	return player, nil
}

func (room *Room) AddPlayer(playerUUID string) error {
	player, error := validPlayerJoin(room, playerUUID)

	if error != nil {
		return error
	}
	room.Players = append(room.Players, *player)
	player.RoomUUID = room.RoomUUID
    room.PublishRoomInfo()
	return nil
}

func (room *Room) PlayerCount() int {
	return len(room.Players)
}

func (room *Room) GetInfoMessage() types.RoomInfoMessage {
	playerListInfo := utils.MapSlice(room.Players, func(p Session) types.PublicPlayerMessage {
		return types.PublicPlayerMessage{
			Username: p.Username,
		}
	})
	return types.RoomInfoMessage{
		Type:        "room_info",
		PlayerCount: room.PlayerCount(),
		PlayerList:  playerListInfo,
	}
}
