package liar_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"uooobarry/liar-groundhog/internal/liar"
)

func TestEngine_StateTransistions(t *testing.T) {
	engine := liar.New()

	// Verify initial state is 'preparing'
	if engine.State != liar.StatePreparing {
		t.Errorf("expected initial state to be 'preparing', got '%s'", engine.State)
	}

	// Test start the game
	t.Run("Able to start the game", func(t *testing.T) {
		err := engine.StartGame()
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
