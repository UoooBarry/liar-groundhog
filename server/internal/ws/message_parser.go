package ws

import (
	"encoding/json"
	"uooobarry/liar-groundhog/internal/message"
)

type MessageHandler func([]byte) (interface{}, error)

var messageParsers = map[string]MessageHandler{
	"login": func(data []byte) (interface{}, error) {
		var msg message.LoginMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
	"room_create": func(data []byte) (interface{}, error) {
		var msg message.RoomCreateMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
	"room_join": func(data []byte) (interface{}, error) {
		var msg message.RoomOpMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
	"room_start": func(data []byte) (interface{}, error) {
		var msg message.RoomOpMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
	"player_action": func(data []byte) (interface{}, error) {
		var msg message.PlayerActionMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
}

var actionParsers = map[string]MessageHandler{
	"player_place_cards": func(data []byte) (interface{}, error) {
		var msg message.PlayerActionMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	},
}
