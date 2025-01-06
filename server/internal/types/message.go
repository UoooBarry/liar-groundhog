package types

type Message struct {
    Type    string `json:"type"`    // Type of the message (e.g., "login")
    Username string `json:"username,omitempty"` // Username for login
    UUID    string `json:"uuid,omitempty"`    // UUID generated for the user
    Content string `json:"content,omitempty"` // Additional content
}

type ActionMessage struct {
    Message    // Embeds the Message struct for common fields
    ActionType string `json:"action_type"` // Type of action
    Details    string `json:"details"`     // Additional details about the action
}
