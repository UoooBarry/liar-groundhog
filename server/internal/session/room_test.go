package session_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"uooobarry/liar-groundhog/internal/session"
)

func TestRoom_AddPlayer(t *testing.T) {
	// Mock the session creation
	playerUsername := "testplayer"
	player := session.CreateSession(nil, playerUsername)
	playerUUID := player.SessionUUID

	room, err := session.CreateRoom(playerUUID)

	t.Run("Added owner successfully", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, 1, room.PlayerCount())
	})

	t.Run("Add player successfully", func(t *testing.T) {
		newPlayer := session.CreateSession(nil, "new player")
		err := room.AddPlayer(newPlayer.SessionUUID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(room.Players))
		assert.Equal(t, "new player", room.Players[1].Username)
	})

	t.Run("Player already in the room", func(t *testing.T) {
		err := room.AddPlayer(playerUUID) // Add the player again
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "A player name 'testplayer' is already in this room")
	})

	t.Run("Room is full", func(t *testing.T) {
		// Add more players to fill the room
		for i := 1; i < (session.MAX_PLAYERS - 1); i++ {
			name := "player-" + strconv.Itoa(i)
			session := session.CreateSession(nil, name)
			sessionUUID := session.SessionUUID
			err := room.AddPlayer(sessionUUID)
			assert.NoError(t, err)
		}

		// Try adding one more player
		newPlayer := session.CreateSession(nil, "extraplayer")
		newPlayerUUID := newPlayer.SessionUUID
		err := room.AddPlayer(newPlayerUUID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The current game room is full.")
	})

	t.Run("Player session not found", func(t *testing.T) {
		nonExistentUUID := "non-existent-uuid"
		err := room.AddPlayer(nonExistentUUID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Player session not exist")
	})
}
