package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
	appErrors "uooobarry/liar-groundhog/internal/errors"
	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/session"
	"uooobarry/liar-groundhog/internal/types"
	"uooobarry/liar-groundhog/internal/utils"
)

func handleRoomCreate(conn *websocket.Conn, msg types.Message) error {
	engine := liar.New()
	room, err := session.CreateRoom(msg.SessionUUID, &engine)
	if err != nil {
		return err
	}
	response := types.Message{
		Type:     "room_create",
		RoomUUID: room.RoomUUID,
		Content:  "Room create successful",
	}
	utils.SendResponse(conn, response)
	return nil
}

func handleRoomJoin(conn *websocket.Conn, msg types.Message) error {
	room, exist := session.FindRoom(&msg.RoomUUID)
	if !exist {
		return appErrors.NewClientError(fmt.Sprintf("Room ID '%s' does not exist", msg.RoomUUID))
	}
	if err := room.AddPlayer(msg.SessionUUID); err != nil {
		return err
	}

	response := types.Message{
		Type:     "room_join",
		RoomUUID: room.RoomUUID,
		Content:  "Room join successful",
	}
	utils.SendResponse(conn, response)
	return nil
}

func handleRoomStart(conn *websocket.Conn, msg types.Message) error {
    room, exist := session.FindRoom(&msg.RoomUUID)
	if !exist {
		return appErrors.NewClientError(fmt.Sprintf("Room ID '%s' does not exist", msg.RoomUUID))
	}

    if err := room.TryStartGame(&msg.SessionUUID); err != nil {
        return err
    }

    response := types.Message{
    	Type:     "room_start",
		RoomUUID: room.RoomUUID,
		Content:  "Room start successful",
        SessionUUID: msg.SessionUUID,
    }
    utils.SendResponse(conn, response)

    return nil
}
