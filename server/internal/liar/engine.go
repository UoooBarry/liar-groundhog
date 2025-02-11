package liar

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"
	"uooobarry/liar-groundhog/utils"
)

const (
	TotalCardCount = 26
	MinPlayers     = 2
	MaxPlayers     = 6
)

var (
	CardPartition = map[Card]int{
		Jack:        6,
		Queen:       6,
		King:        6,
		Ace:         6,
		BigJoker:    1,
		LittleJoker: 1,
	}

	ErrInvalidAction    = errors.New("invalid action for current state")
	ErrGameNotRunning   = errors.New("game is not running")
	ErrNotEnoughCards   = errors.New("not enough cards to deal")
	ErrInvalidGameState = errors.New("invalid game state")
	ErrCardNotHeld      = errors.New("player does not hold specified card")
)

type Engine struct {
	State          GameState
	Cards          []Card
	CurrentAction  GameAction
	LastPlaceCards []Card
	BetCard        Card
	randSource     rand.Source
}

func New() *Engine {
	return &Engine{
		State:         StatePreparing,
		Cards:         newPackOfCards(),
		CurrentAction: PlaceCards,
		BetCard:       Ace,
		randSource:    rand.NewSource(time.Now().UnixNano()),
	}
}

func NewWithSeed(seed int64) *Engine {
	return &Engine{
		State:         StatePreparing,
		Cards:         newPackOfCards(),
		CurrentAction: PlaceCards,
		BetCard:       Ace,
		randSource:    rand.NewSource(seed),
	}
}

func newPackOfCards() []Card {
	cards := make([]Card, 0, TotalCardCount)
	for card, count := range CardPartition {
		cards = append(cards, slices.Repeat([]Card{card}, count)...)
	}

	return cards
}

func (e *Engine) validStateAndAction(action GameAction) error {
	if e.State != StateInGame {
		return ErrGameNotRunning
	}
	if e.CurrentAction != action {
		return ErrInvalidAction
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
	r := rand.New(e.randSource)
	r.Shuffle(len(e.Cards), func(i, j int) {
		e.Cards[i], e.Cards[j] = e.Cards[j], e.Cards[i]
	})
}

func (e *Engine) DealCards(num int) ([]Card, error) {
	if num <= 0 {
		return nil, fmt.Errorf("must deal at least 1 card")
	}
	if len(e.Cards) < num {
		return nil, fmt.Errorf("%w: requested %d, available %d", ErrNotEnoughCards, num, len(e.Cards))
	}

	dealtCards := make([]Card, num)
	copy(dealtCards, e.Cards[:num])
	e.Cards = e.Cards[num:]
	return dealtCards, nil
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
	if err := e.validStateAndAction(PlaceCards); err != nil {
		return holdingCards, err
	}
	if len(placedCards) == 0 {
		return holdingCards, fmt.Errorf("must place at least one card")
	}

	for _, card := range placedCards {
		i := slices.Index(holdingCards, card)
		if i < 0 {
			return holdingCards, fmt.Errorf("%w: %s", ErrCardNotHeld, card)
		}
	}

	// Move cards to public deck
	e.Cards = append(e.Cards, placedCards...)
	
	// Remove from player's hand
	newHand := make([]Card, 0, len(holdingCards)-len(placedCards))
	for _, card := range holdingCards {
		if !slices.Contains(placedCards, card) {
			newHand = append(newHand, card)
		}
	}

	e.LastPlaceCards = placedCards
	e.nextAction()
	return newHand, nil
}

func (e *Engine) Declare(doubt bool) (DeclareResult, error) {
	if err := e.validStateAndAction(Doubt); err != nil {
		return Skip, err
	}

	if !doubt {
		e.nextAction()
		return Skip, nil
	}

	result := Lied
	if utils.SliceIsAll(e.LastPlaceCards, func(c Card) bool {
		return c == e.BetCard || c == BigJoker || c == LittleJoker
	}) {
		result = Truthful
	}

	e.nextAction()
	return result, nil
}
