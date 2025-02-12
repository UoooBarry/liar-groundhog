package liar

// GameState represents the state of the game.
type GameState string
type GameAction string

const (
	StatePreparing  GameState = "preparing"
	StateInGame     GameState = "in_game"
	StateSettlement GameState = "settlement"
)

const (
	Doubt      GameAction = "doubt"
	PlaceCards GameAction = "place_cards"
)

type GameEngine interface {
	StartGame() error
	ResetGame() error
	EndGame() error
	GetState() GameState
}

type Card string

const (
	Jack        Card = "jack"
	Queen       Card = "queen"
	King        Card = "king"
	Ace         Card = "Ace"
	BigJoker    Card = "big_joker"
	LittleJoker Card = "little_joker"
)

type DeclareResult string

const (
	Truthful DeclareResult = "truthful"
	Lied     DeclareResult = "lied"
	Skip     DeclareResult = "skip"
)
