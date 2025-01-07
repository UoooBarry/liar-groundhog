package errors

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

// Helper to create a ClientError
func NewClientError(message string) error {
	return &ClientError{Message: message}
}
