package errors

type Severity int

const (
	INFO Severity = iota
	WARN
	ERROR
)

type LoggableError struct {
	Message  string
	Severity Severity
}

func (e *LoggableError) Error() string {
	return e.Message
}

func NewLoggableError(message string, severity Severity) error {
	return &LoggableError{Message: message, Severity: severity}
}
