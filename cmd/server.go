package main

import (
	"fmt"
	"github.com/nightfurytech/events-stream-processor/internal/connection"
	"github.com/nightfurytech/events-stream-processor/internal/process"
	"net/http"
)

func main() {
	messageCh := make(chan []byte, 1000)
	h := connection.NewHandler(messageCh)
	db := connection.Create()
	p := process.NewProcessor(db, messageCh)
	go http.HandleFunc("/ws", h.Handle)
	go p.EventProcessor()
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
