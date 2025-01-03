package ws

type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g., "login")
	Username string `json:"username,omitempty"` // Username for login
	UUID    string `json:"uuid,omitempty"`    // UUID generated for the user
	Content string `json:"content,omitempty"` // Additional content
}
