package liar

import (
	"errors"
	"fmt"
)

// GameState represents the state of the game.
type GameState string

const (
	StatePreparing  GameState = "preparing"
	StateInGame     GameState = "in_game"
	StateSettlement GameState = "settlement"
)

type Engine struct {
	State GameState
}

func New() Engine {
	return Engine{State: StatePreparing}
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
