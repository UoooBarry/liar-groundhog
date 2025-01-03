package liar

import (
	"errors"
	"fmt"

	"uooobarry/liar-groundhog/internal/session"
)

const MAX_PLAYERS  = 4

type Engine struct {
	Players []session.Session
	RoomID string
}

func validToJoin(engine *Engine, playerUUID string) (session.Session, error) {
	player, exist := session.FindSession(playerUUID)

	if !exist {
		return player, fmt.Errorf("Player session not exist '%s'", playerUUID)
	}

	if len(engine.Players) >= MAX_PLAYERS {
		return player, errors.New("The current game room is full.")
	}

	return player, nil
}

func (e *Engine) AddPlayer(playerUUID string) error {
	player, error := validToJoin(e, playerUUID)

	if error != nil {
		return error
	}
    e.Players = append(e.Players, player)
	return nil
}
