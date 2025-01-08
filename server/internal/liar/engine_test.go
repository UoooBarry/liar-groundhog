package liar_test

import (
	"log"
	"testing"

	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestEngine_StateTransistions(t *testing.T) {
	engine := liar.New()

	// Verify initial state is 'preparing'
	if engine.State != types.StatePreparing {
		t.Errorf("expected initial state to be 'preparing', got '%s'", engine.State)
	}

	// Test start the game
	t.Run("Able to start the game", func(t *testing.T) {
		err := engine.StartGame()
		log.Println(engine.Cards)
		assert.NoError(t, err)
	})

	// Test invalid transistion call StartGame() when 'in_game'
	t.Run("Call StartGame when state is 'in_game'", func(t *testing.T) {
		err := engine.StartGame()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot start game: game is not in the preparing state")
	})

	t.Run("Call RestGame when state is 'in_game'", func(t *testing.T) {
		err := engine.ResetGame()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot reset game: game is not in the settlement state")
	})

	t.Run("Able to end the game", func(t *testing.T) {
		err := engine.EndGame()
		assert.NoError(t, err)
	})

	t.Run("Able to reset the game", func(t *testing.T) {
		err := engine.ResetGame()
		assert.NoError(t, err)
	})
}

func TestEngine_Cards(t *testing.T) {
	engine := liar.New()

	t.Run("Able to generate cards", func(t *testing.T) {
		// Create a map to count the occurrences of each card in the Engine's Cards slice
		cardCount := make(map[types.Card]int)

		// Count the occurrences of each card
		for _, card := range engine.Cards {
			cardCount[card]++
		}

		for card, expectedCount := range liar.CardPartition {
			actualCount, exists := cardCount[card]
			if !exists {
				t.Errorf("Card %v not found in the Engine's Cards", card)
			} else if actualCount != expectedCount {
				t.Errorf("Expected %v cards of type %v, but found %v", expectedCount, card, actualCount)
			}
		}
	})
}

func TestEngine_DealCards(t *testing.T) {
	// Initialize the engine
	engine := liar.New()

	// Initial deck size
	initialDeckSize := len(engine.Cards)

	// Number of cards to deal
	cardsPerPlayer := 5
	numPlayers := 4

	// Ensure the deck has enough cards
	t.Run("Able to deal cards to players", func(t *testing.T) {
		// Deal 4 cards to each player (or simulated players in this case)
		_, err := engine.DealCards(cardsPerPlayer * numPlayers)
		assert.NoError(t, err) // Ensure no errors occurred while dealing cards

		// Verify that the deck size is reduced
		remainingDeckSize := len(engine.Cards)
		expectedRemainingDeckSize := initialDeckSize - cardsPerPlayer*numPlayers
		assert.Equal(t, expectedRemainingDeckSize, remainingDeckSize, "Deck size after dealing cards is incorrect")
	})

	// Handle case where there are not enough cards in the deck
	t.Run("Handle case where there are not enough cards", func(t *testing.T) {
		// Let's assume there are fewer than 100 cards in the deck
		_, err := engine.DealCards(100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enough cards to deal")
	})

	// Verify the deck size is correctly updated after dealing
	t.Run("Deck size after dealing cards", func(t *testing.T) {
        // Check the remaining deck size
		remainingDeckSize := len(engine.Cards)
		expectedRemainingDeckSize := initialDeckSize - cardsPerPlayer*numPlayers
		assert.Equal(t, expectedRemainingDeckSize, remainingDeckSize, "Deck size after second deal is incorrect")
	})
}
