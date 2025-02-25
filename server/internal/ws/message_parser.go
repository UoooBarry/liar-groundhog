package ws

import (
	"encoding/json"
	"uooobarry/liar-groundhog/internal/message"
)

type MessageHandler func([]byte) (interface{}, error)

// 生成解析函数的工厂函数
func createParser(msg interface{}) MessageHandler {
	return func(data []byte) (interface{}, error) {
		err := json.Unmarshal(data, msg)
		return msg, err
	}
}

var messageParsers = map[string]MessageHandler{
	"login":         createParser(&message.LoginMessage{}),
	"room_create":   createParser(&message.RoomCreateMessage{}),
	"room_join":     createParser(&message.RoomOpMessage{}),
	"room_start":    createParser(&message.RoomOpMessage{}),
	"player_action": createParser(&message.PlayerActionMessage{}),
}

var actionParsers = map[string]MessageHandler{
	"player_place_cards": createParser(&message.PlayerActionMessage{}),
}
