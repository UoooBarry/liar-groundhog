package message

import (
	"uooobarry/liar-groundhog/internal/liar"
)

type MessageInterface interface{}

type PublicPlayerMessage struct {
	Username string `json:"username"`
}

type Message struct {
	Type     string `json:"type"`               // Type of the message (e.g., "login")
	Username string `json:"username,omitempty"` // Username for login
	Content  string `json:"content,omitempty"`  // Additional content
}

type LoginMessage struct {
	Message
	Username string `json:"username,omitempty"` // Username for login
}

type LoginSuccessMessage struct {
	Message
	Username    string `json:"username,omitempty"`    // Username for login
	SessionUUID string `json:"sessionuuid,omitempty"` // UUID generated for the user
}

type RoomOpMessage struct {
	Message
	SessionUUID string `json:"sessionuuid,omitempty"` // UUID generated for the user
	RoomUUID    string `json:"roomuuid,omitempty"`
}

type RoomCreateMessage struct {
	Message
	SessionUUID string `json:"sessionuuid"` // UUID generated for the user
}

type PlayerActionMessage struct {
	Message            // Embeds the Message struct for common fields
	SessionUUID string `json:"sessionuuid,omitempty"` // UUID generated for the user
	ActionType  string `json:"action_type"`           // Type of action
	RoomUUID    string `json:"roomuuid,omitempty"`
}

type RoomInfoMessage struct {
	Type        string                `json:"type"` // Type of the message (e.g., "login")
	PlayerCount int                   `json:"player_count"`
	PlayerList  []PublicPlayerMessage `json:"player_list"`
	GameState   liar.GameState        `json:"game_state"`
	Surivals    int                   `json:"surivals"`
}

type PlayerHoldingCardsMessage struct {
	Type         string      `json:"type"` // Type of the message (e.g., "player_holding_cards")
	HoldingCards []liar.Card `json:"holding_cards"`
	SessionUUID  string      `json:"sessionuuid,omitempty"` // UUID generated for the user
	Username     string      `json:"username,omitempty"`    // Username for login
}

type PlayerPlaceCardsMessage struct {
	PlayerActionMessage
	Cards []liar.Card `json:"cards"`
}

type PlayerDeclareLiarMessage struct {
	PlayerActionMessage
	Doubt bool `json:"doubt"`
}

type RoomBoardCastDeclareMessage struct {
	Type    string             `json:"type"`
	Refname string             `json:"ref_name"`
	Suspect string             `json:"suspect"`
	Result  liar.DeclareResult `json:"declare_result"`
}

type RoomBoardPlayerPlaceCardMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Number   int    `json:"number"`
}
