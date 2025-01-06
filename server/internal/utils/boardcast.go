package utils

import (
	"github.com/gorilla/websocket"
	"log"
	"uooobarry/liar-groundhog/internal/types"
)

func SendError(conn *websocket.Conn, errMsg string) {
	response := types.Message{
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
