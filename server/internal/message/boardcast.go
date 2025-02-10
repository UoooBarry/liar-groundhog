package message

import (
	"log"
	"github.com/gorilla/websocket"
)

func SendError(conn *websocket.Conn, errMsg string) {
	response := Message{
		Type:    "error",
		Content: errMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Println("Write error:", err)
	}
}

func SendResponse(conn *websocket.Conn, msg any) {
	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Write error:", err)
	}
}
