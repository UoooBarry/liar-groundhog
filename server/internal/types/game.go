package types

// GameState represents the state of the game.
type GameState string

const (
	StatePreparing  GameState = "preparing"
	StateInGame     GameState = "in_game"
	StateSettlement GameState = "settlement"
)

type GameEngine interface {
    StartGame() error 
    ResetGame() error 
    EndGame() error
    GetState() GameState
}
