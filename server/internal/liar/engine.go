package liar

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
	"uooobarry/liar-groundhog/internal/types"
)

type Card string

const (
	Jack        Card = "jack"
	Queen       Card = "queen"
	King        Card = "king"
	Ace         Card = "Ace"
	BigJoker    Card = "big_joker"
	LittleJoker Card = "little_joker"
)

const TotalCardCount = 26

var CardPartition = map[Card]int{
	Jack:        6,
	Queen:       6,
	King:        6,
	Ace:         6,
	BigJoker:    1,
	LittleJoker: 1,
}

type Engine struct {
	State types.GameState
	Cards []Card
}

func New() Engine {
	cards := make([]Card, 0, TotalCardCount)
	for card, count := range CardPartition {
		cards = append(cards, slices.Repeat([]Card{card}, count)...)
	}

	return Engine{State: types.StatePreparing, Cards: cards}
}

func (e *Engine) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(e.Cards) - 1; i > 0; i-- {
		// Randomly select an index j between 0 and i (inclusive)
		j := r.Intn(i + 1)

		// Swap the elements at i and j
		e.Cards[i], e.Cards[j] = e.Cards[j], e.Cards[i]
	}
}

func (e *Engine) StartGame() error {
	if e.State != types.StatePreparing {
		return errors.New("cannot start game: game is not in the preparing state")
	}

	e.Shuffle()
	e.State = types.StateInGame
	fmt.Println("Game has started")
	fmt.Println(e.Cards)
	return nil
}

func (e *Engine) EndGame() error {
	if e.State != types.StateInGame {
		return errors.New("cannot end game: game is not in progress")
	}
	e.State = types.StateSettlement
	fmt.Println("Game is now in settlement state")
	return nil
}

func (e *Engine) ResetGame() error {
	if e.State != types.StateSettlement {
		return errors.New("cannot reset game: game is not in the settlement state")
	}
	e.State = types.StatePreparing
	fmt.Println("Game has been reset to preparing state")
	return nil
}

func (e *Engine) GetState() types.GameState {
	return e.State
}
