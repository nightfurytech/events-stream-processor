package connection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Handler struct {
	messageCh chan []byte
	clients   map[*websocket.Conn]bool // A set to keep track of all connected clients
	mu        sync.Mutex
}

func NewHandler(ch chan []byte) *Handler {
	return &Handler{
		messageCh: ch,
		clients:   make(map[*websocket.Conn]bool),
		mu:        sync.Mutex{},
	}
}

// Upgrader is used to upgrade an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	// Add the connection to the set of clients
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	fmt.Println("Client connected!")

	// Listen for messages from this client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
			break
		}
		h.messageCh <- msg
	}
}
