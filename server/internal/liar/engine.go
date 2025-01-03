package liar

import (
	"errors"
	"fmt"

	"uooobarry/liar-groundhog/internal/session"
)

const MAX_PLAYERS  = 4

// GameState represents the state of the game.
type GameState string

const (
	StatePreparing GameState = "preparing"
	StateInGame    GameState = "in_game"
	StateSettlement GameState = "settlement"
)

type Engine struct {
	Room session.Room
    State GameState
}

func New() Engine {
	_, room := session.CreateRoom()
    return Engine{Room: room, State: StatePreparing}
}

func validToJoin(engine *Engine, playerUUID string) (session.Session, error) {
	player, exist := session.FindSession(playerUUID)

	if !exist {
		return player, fmt.Errorf("Player session not exist '%s'", playerUUID)
	}

	if len(engine.Room.Players) >= MAX_PLAYERS {
		return player, errors.New("The current game room is full.")
	}

	if _, inRoom := engine.Room.FindPlayerInRoom(player.Username); inRoom {
		return player, fmt.Errorf("Player '%s' is already in the room", player.Username)
	}

	return player, nil
}

func (e *Engine) AddPlayer(playerUUID string) error {
	player, error := validToJoin(e, playerUUID)

	if error != nil {
		return error
	}
	e.Room.Players = append(e.Room.Players, player)
	player.RoomUUID = e.Room.RoomUUID
	return nil
}

func (e *Engine) StartGame() error {
	if e.State != StatePreparing {
		return errors.New("cannot start game: game is not in the preparing state")
	}
	e.State = StateInGame
	fmt.Println("Game has started")
	return nil
}

func (e *Engine) EndGame() error {
	if e.State != StateInGame {
		return errors.New("cannot end game: game is not in progress")
	}
	e.State = StateSettlement
	fmt.Println("Game is now in settlement state")
	return nil
}

func (e *Engine) ResetGame() error {
	if e.State != StateSettlement {
		return errors.New("cannot reset game: game is not in the settlement state")
	}
	e.State = StatePreparing
	fmt.Println("Game has been reset to preparing state")
	return nil
}
