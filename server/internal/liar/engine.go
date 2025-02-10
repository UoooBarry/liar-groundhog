package liar

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
	"uooobarry/liar-groundhog/internal/utils"
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
	State          GameState
	Cards          []Card
	CurrentAction  GameAction
	LastPlaceCards []Card
	BetCard        Card
}

func New() Engine {
	return Engine{State: StatePreparing,
		Cards:         newPackOfCards(),
		CurrentAction: PlaceCards,
		BetCard:       Ace}
}

func newPackOfCards() []Card {
	cards := make([]Card, 0, TotalCardCount)
	for card, count := range CardPartition {
		cards = append(cards, slices.Repeat([]Card{card}, count)...)
	}

	return cards
}

func (e *Engine) ValidStateAndAction(action GameAction) error {
	if e.CurrentAction != action {
		return errors.New("Not a valid action")
	}

	if e.State != StateInGame {
		return errors.New("The game is not running")
	}

	return nil
}

func (e *Engine) nextAction() {
	if e.CurrentAction == Doubt {
		e.CurrentAction = PlaceCards
	} else {
		e.CurrentAction = Doubt
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

func (e *Engine) DealCards(num int) ([]Card, error) {
	if len(e.Cards) < num {
		return e.Cards, errors.New("not enough cards to deal")
	}

	dealedCards := e.Cards[:num]
	e.Cards = e.Cards[num:]
	return dealedCards, nil
}

func (e *Engine) StartGame() error {
	if e.State != StatePreparing {
		return errors.New("cannot start game: game is not in the preparing state")
	}

	e.Shuffle()
	e.State = StateInGame
	fmt.Println("Game has started")
	fmt.Println(e.Cards)
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
	e.Cards = newPackOfCards()
	fmt.Println("Game has been reset to preparing state")
	return nil
}

func (e *Engine) GetState() GameState {
	return e.State
}

func (e *Engine) PlaceCard(holdingCards []Card, placedCards []Card) ([]Card, error) {
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

func (e *Engine) Declare(doubt bool) DeclareResult {
	var result DeclareResult
	if !doubt {
		result = Skip
	}

	// Every cards it claims are the goal card
	if utils.SliceIsAll(e.LastPlaceCards, func(c Card) bool {
		return c == e.BetCard || c == BigJoker || c == LittleJoker
	}) {
		result = Truthful
	} else {
		result = Lied
	}

	e.nextAction()
	return result
}
