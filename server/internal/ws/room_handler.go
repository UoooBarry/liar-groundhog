package ws

import (
	"fmt"

	appErrors "uooobarry/liar-groundhog/internal/errors"
	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/message"
	"uooobarry/liar-groundhog/internal/session"

	"github.com/gorilla/websocket"
)

func handleRoomCreate(conn *websocket.Conn, msg message.RoomCreateMessage) error {
	engine := liar.New()
	room, err := session.CreateRoom(msg.SessionUUID, &engine)
	if err != nil {
		return err
	}
	response := message.RoomOpMessage{
		Message: message.Message{
			Type:    "room_create",
			Content: "Room create successful",
		},
		RoomUUID:    room.RoomUUID,
		SessionUUID: msg.SessionUUID,
	}
	message.SendResponse(conn, response)
	return nil
}

func handleRoomJoin(conn *websocket.Conn, msg message.RoomOpMessage) error {
	room, exist := session.FindRoom(msg.RoomUUID)
	if !exist {
		return fmt.Errorf("Room ID '%s' does not exist", msg.RoomUUID)
	}
	if err := room.AddPlayer(msg.SessionUUID); err != nil {
		return err
	}

	response := message.RoomOpMessage{
		Message: message.Message{
			Type:    "room_join",
			Content: "Room join successful",
		},
		RoomUUID: room.RoomUUID,
	}
	message.SendResponse(conn, response)
	return nil
}

func handleRoomStart(conn *websocket.Conn, msg message.RoomOpMessage) error {
	room, exist := session.FindRoom(msg.RoomUUID)
	if !exist {
		return appErrors.NewClientError("Room not existed")
	}

	if err := room.TryStartGame(&msg.SessionUUID); err != nil {
		return err
	}
	response := message.RoomOpMessage{
		Message: message.Message{
			Type:    "room_start",
			Content: "Room start successful",
		},
		RoomUUID:    room.RoomUUID,
		SessionUUID: msg.SessionUUID,
	}
	message.SendResponse(conn, response)

	return nil
}

func handlePlayerAction(conn *websocket.Conn, msg message.PlayerActionMessage, raw []byte) error {
	room, exist := session.FindRoom(msg.RoomUUID)
	if !exist {
		return appErrors.NewClientError(fmt.Sprintf("Room ID '%s' does not exist", msg.RoomUUID))
	}
	parser, exist := actionParsers[msg.ActionType]
	if exist {
		parsed, err := parser(raw)
		if err != nil {
			return appErrors.NewLoggableError(err.Error(), appErrors.ERROR)
		}

		switch v := parsed.(type) {
		case message.PlayerPlaceCardsMessage:
			rsp, err := handleRoomPlaceCard(room, v)
			if err != nil {
				return err
			}
			message.SendResponse(conn, rsp)
		}
	}

	return nil
}

func handleRoomPlaceCard(room *session.Room, msg message.PlayerPlaceCardsMessage) (*message.PlayerPlaceCardsMessage, error) {
	if err := room.PlayerPlaceCard(msg.SessionUUID, msg.Cards); err != nil {
		return nil, err
	}
	response := message.PlayerPlaceCardsMessage{
		PlayerActionMessage: message.PlayerActionMessage{
			ActionType: "player_place_cards",
			RoomUUID:   room.RoomUUID,
			Message: message.Message{
				Type:    "player_action",
				Content: "Success",
			},
			SessionUUID: msg.SessionUUID,
		},
	}
	return &response, nil
}
