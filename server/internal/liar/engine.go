package liar

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
	"uooobarry/liar-groundhog/internal/types"
)

const TotalCardCount = 26

var CardPartition = map[types.Card]int{
	types.Jack:        6,
	types.Queen:       6,
	types.King:        6,
	types.Ace:         6,
	types.BigJoker:    1,
	types.LittleJoker: 1,
}

type Engine struct {
	State types.GameState
	Cards []types.Card
}

func New() Engine {
	return Engine{State: types.StatePreparing, Cards: newPackOfCards()}
}

func newPackOfCards() []types.Card {
	cards := make([]types.Card, 0, TotalCardCount)
	for card, count := range CardPartition {
		cards = append(cards, slices.Repeat([]types.Card{card}, count)...)
	}

	return cards
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

func (e *Engine) DealCards(num int) ([]types.Card, error) {
	if len(e.Cards) < num {
		return e.Cards, errors.New("not enough cards to deal")
	}

    dealedCards := e.Cards[:num]
	e.Cards = e.Cards[num:]
	return dealedCards, nil
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
	e.Cards = newPackOfCards()
	fmt.Println("Game has been reset to preparing state")
	return nil
}

func (e *Engine) GetState() types.GameState {
	return e.State
}
