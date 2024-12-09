// internal/websocket/hub.go
package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/ZetoOfficial/chat-server/internal/chat/domain"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]string
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]string),
	}
}

func (h *Hub) Register(conn *websocket.Conn, username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = username
	log.Printf("Клиент %s подключился", username)
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	username, exists := h.clients[conn]
	if exists {
		delete(h.clients, conn)
		log.Printf("Клиент %s отключился", username)
	}
}

func (h *Hub) BroadcastMessage(msg domain.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(domain.BroadcastMessage{
		Type: "message",
		Data: msg,
	})
	if err != nil {
		log.Println("Ошибка при маршалинге сообщения:", err)
		return
	}

	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения клиенту %s: %v", h.clients[conn], err)
			go h.Unregister(conn)
		}
	}
}

func (h *Hub) BroadcastUserList() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := []string{}
	for _, username := range h.clients {
		users = append(users, username)
	}

	data, err := json.Marshal(domain.BroadcastMessage{
		Type: "users",
		Data: users,
	})
	if err != nil {
		log.Println("Ошибка при маршалинге списка пользователей:", err)
		return
	}

	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Ошибка при отправке списка пользователям клиенту %s: %v", h.clients[conn], err)
			go h.Unregister(conn)
		}
	}
}
