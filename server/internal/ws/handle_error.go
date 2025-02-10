package ws

import (
	"log"
	"uooobarry/liar-groundhog/internal/errors"
	"uooobarry/liar-groundhog/internal/message"

	"github.com/gorilla/websocket"
)

func HandleError(conn *websocket.Conn, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *errors.ClientError:
		message.SendError(conn, e.Message)
	case *errors.LoggableError:
		switch e.Severity {
		case errors.INFO:
			log.Println("INFO:", e.Message)
		case errors.WARN:
			log.Println("WARNING:", e.Message)
		case errors.ERROR:
			log.Println("ERROR:", e.Message)
		}
	default:
		message.SendError(conn, e.Error())
	}
}
