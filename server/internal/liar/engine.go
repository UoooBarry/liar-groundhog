package liar

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
	"uooobarry/liar-groundhog/internal/types"
	"uooobarry/liar-groundhog/internal/utils"
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
	State          types.GameState
	Cards          []types.Card
	CurrentAction  types.GameAction
	LastPlaceCards []types.Card
	BetCard        types.Card
}

func New() Engine {
	return Engine{State: types.StatePreparing,
		Cards:         newPackOfCards(),
		CurrentAction: types.PlaceCards,
		BetCard:       types.Ace}
}

func newPackOfCards() []types.Card {
	cards := make([]types.Card, 0, TotalCardCount)
	for card, count := range CardPartition {
		cards = append(cards, slices.Repeat([]types.Card{card}, count)...)
	}

	return cards
}

func (e *Engine) ValidStateAndAction(action types.GameAction) error {
	if e.CurrentAction != action {
		return errors.New("Not a valid action")
	}

	if e.State != types.StateInGame {
		return errors.New("The game is not running")
	}

	return nil
}

func (e *Engine) nextAction() {
	if e.CurrentAction == types.Doubt {
		e.CurrentAction = types.PlaceCards
	} else {
		e.CurrentAction = types.Doubt
	}
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

func (e *Engine) PlaceCard(holdingCards []types.Card, placedCards []types.Card) ([]types.Card, error) {
	for _, card := range placedCards {
		i := slices.Index(holdingCards, card)
		if i < 0 {
			return holdingCards, fmt.Errorf("You do not have %s", card)
		}

		// put cards back to the public
		e.Cards = append(e.Cards, holdingCards[i])
        holdingCards = append(holdingCards[:i], holdingCards[i+1:]...)
	}

	e.LastPlaceCards = placedCards
    e.nextAction()
	return holdingCards, nil
}

func (e *Engine) Declare(doubt bool) types.DeclareResult {
    var result types.DeclareResult
    if (!doubt) {
        result = types.Skip
    }

	// Every cards it claims are the goal card
	if utils.SliceIsAll(e.LastPlaceCards, func(c types.Card) bool {
		return c == e.BetCard || c == types.BigJoker || c == types.LittleJoker
	}) {
		result = types.Truthful
	} else {
		result = types.Lied
	}

    e.nextAction()
    return result
}
