package liar_test

import (
	"testing"
	"strconv"
	"fmt"

	"github.com/stretchr/testify/assert"
	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/session"
)

func TestEngine_AddPlayer(t *testing.T) {
	// Mock the session creation
	playerUsername := "testplayer"
	playerUUID, _ := session.CreateSession(nil, playerUsername)

	engine := liar.New()

	t.Run("Add player successfully", func(t *testing.T) {
		err := engine.AddPlayer(playerUUID)
		fmt.Println(engine.Room.Players)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(engine.Room.Players))
		assert.Equal(t, playerUsername, engine.Room.Players[0].Username)
	})

	t.Run("Player already in the room", func(t *testing.T) {
		err := engine.AddPlayer(playerUUID) // Add the player again
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Player 'testplayer' is already in the room")
	})

	t.Run("Room is full", func(t *testing.T) {
		// Add more players to fill the room
		for i := 1; i < liar.MAX_PLAYERS; i++ {
			name := "player-" + strconv.Itoa(i)
			sessionUUID, _ := session.CreateSession(nil, name)
			err := engine.AddPlayer(sessionUUID)
			assert.NoError(t, err)
		}

		// Try adding one more player
		newPlayerUUID, _ := session.CreateSession(nil, "extraplayer")
		err := engine.AddPlayer(newPlayerUUID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The current game room is full.")
	})

	t.Run("Player session not found", func(t *testing.T) {
		nonExistentUUID := "non-existent-uuid"
		err := engine.AddPlayer(nonExistentUUID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Player session not exist")
	})
}

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
    t.Run("Call StartGame when state is 'in_game'", func (t *testing.T)  {
        err := engine.StartGame()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "cannot start game: game is not in the preparing state")
    })

    t.Run("Call RestGame when state is 'in_game'", func (t *testing.T)  {
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
