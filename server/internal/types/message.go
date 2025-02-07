package types

type PublicPlayerMessage struct {
	Username string `json:"username"`
}

type Message struct {
	Type        string `json:"type"`                  // Type of the message (e.g., "login")
	Username    string `json:"username,omitempty"`    // Username for login
	SessionUUID string `json:"sessionuuid,omitempty"` // UUID generated for the user
	Content     string `json:"content,omitempty"`     // Additional content
	RoomUUID    string `json:"roomuuid,omitempty"`
	Cards       []Card `json:"cards,omitempty"`
}

type ActionMessage struct {
	Message           // Embeds the Message struct for common fields
	ActionType string `json:"action_type"` // Type of action
	Details    string `json:"details"`     // Additional details about the action
}

type RoomInfoMessage struct {
	Type        string                `json:"type"` // Type of the message (e.g., "login")
	PlayerCount int                   `json:"player_count"`
	PlayerList  []PublicPlayerMessage `json:"player_list"`
	GameState   GameState             `json:"game_state"`
	Surivals    int                   `json:"surivals"`
}

type PlayerHoldingCardsMessage struct {
	Type         string `json:"type"` // Type of the message (e.g., "player_holding_cards")
	HoldingCards []Card `json:"holding_cards"`
	SessionUUID  string `json:"sessionuuid,omitempty"` // UUID generated for the user
	Username     string `json:"username,omitempty"`    // Username for login
}

type MessageInterface interface{}

type PlayerPlaceCardsMessage struct {
	Type     string `json:"type"`               // Type of the message (e.g., "player_holding_cards")
	Username string `json:"username,omitempty"` // Username for login
	Number   int    `json:"number"`
}

type PlayerActionMessage struct {
	Type        string `json:"type"` // Type of the message (e.g., "player_holding_cards")
	SessionUUID string `json:"sessionuuid,omitempty"`
	Cards       []Card `json:"cards"`
	RoomUUID    string `json:"roomuuid,omitempty"`
}

type PlayerDeclareLiarMessage struct {
	Type        string `json:"type"`
	SessionUUID string `json:"sessionuuid,omitempty"`
}

type RoomBoardCastDeclareMessage struct {
	Type    string        `json:"type"`
	Refname string        `json:"username"`
	Suspect string        `json:"suspect"`
	Result  DeclareResult `json:"declare_result"`
}
