package errors

import (
	"github.com/gorilla/websocket"
	"log"
	"uooobarry/liar-groundhog/internal/utils"
)

func HandleError(conn *websocket.Conn, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *ClientError:
		utils.SendError(conn, e.Message)
	case *LoggableError:
		switch e.Severity {
		case INFO:
			log.Println("INFO:", e.Message)
		case WARN:
			log.Println("WARNING:", e.Message)
		case ERROR:
			log.Println("ERROR:", e.Message)
		}
	default:
		log.Println("Unexpected error:", err)
	}
}
