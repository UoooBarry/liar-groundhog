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
    GameState GameState `json:"game_state"`
}
