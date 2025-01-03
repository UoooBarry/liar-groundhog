package main

import (
	"fmt"
	"log"
	"net/http"
	"uooobarry/liar-groundhog/internal/ws"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello go!")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", ws.HandleWebSocket)
	port := "8080"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
