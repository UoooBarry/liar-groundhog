package session

import (
	"fmt"
	"log"
	"sync"

	"uooobarry/liar-groundhog/internal/errors"
	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/types"
	"uooobarry/liar-groundhog/internal/utils"

	"github.com/google/uuid"
)

const MAX_PLAYERS = 4
const MIN_PLAYERS_TO_START =4

var rooms = struct {
	sync.Mutex
	data map[string]*Room // uuid -> Session
}{
	data: make(map[string]*Room),
}

type Room struct {
	RoomUUID string
	Players  []Session
	Engine   types.GameEngine
    OwnerUUID string
}

func CreateRoom(ownerUUID string, gameEngine *liar.Engine) (*Room, error) {
	rooms.Lock()
	defer rooms.Unlock()

	owner, exist := FindSession(ownerUUID)

	if !exist {
		return nil, fmt.Errorf("Player session not exist '%s'", ownerUUID)
	}
	uuid := uuid.NewString()
	room := &Room{RoomUUID: uuid, Engine: gameEngine}
	rooms.data[uuid] = room
	room.Players = append(room.Players, *owner)
    room.OwnerUUID = *&owner.SessionUUID

	log.Printf("Created room UUID '%s'", uuid)
	return room, nil
}

func FindRoom(uuid *string) (*Room, bool) {
    if (uuid == nil) {
        return nil, false
    }
	room, exist := rooms.data[*uuid]

	return room, exist
}

func (room *Room) FindPlayerInRoom(username *string) (*Session, bool) {
    if (username == nil) {
        return nil, false
    }

	for _, player := range room.Players {
		if player.Username == *username {
			return &player, true
		}
	}

	return nil, false
}

func (room *Room) FindPlayerInRoomByUUID(uuid *string) (*Session, bool) {
    if (uuid == nil) {
        return nil, false
    }

	for _, player := range room.Players {
		if player.SessionUUID == *uuid {
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
		return nil, errors.NewLoggableError(fmt.Sprintf("Player session not exist '%s'", playerUUID), errors.ERROR)
	}

	if len(room.Players) >= MAX_PLAYERS {
		return player, errors.NewClientError("The current game room is full.")
	}

	if _, inRoom := room.FindPlayerInRoom(&player.Username); inRoom {
        return player, errors.NewClientError(fmt.Sprintf("A player name '%s' is already in this room", player.Username))
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
        GameState: room.Engine.GetState(),
	}
}

func (room *Room) TryStartGame(playerUUID *string) error {
    player, exist := room.FindPlayerInRoomByUUID(playerUUID)
    if !exist || room.OwnerUUID != player.SessionUUID {
        return errors.NewClientError("Invalid player")
    }
    if err := room.Engine.StartGame(); err != nil {
        return err
    }
    if room.PlayerCount() < MIN_PLAYERS_TO_START {
        return errors.NewClientError("Require at least 4 players to start the game.")
    }

    room.PublishRoomInfo()

    return nil
}
