package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ZetoOfficial/chat-server/internal/chat/domain"
	"github.com/ZetoOfficial/chat-server/internal/chat/usecase"
	mywebsocket "github.com/ZetoOfficial/chat-server/internal/websocket"
	"github.com/ZetoOfficial/chat-server/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const historyMessagesCount = 50

type Handler struct {
	UseCase *usecase.ChatUseCase
	Hub     *mywebsocket.Hub
}

func NewHandler(usecase *usecase.ChatUseCase, hub *mywebsocket.Hub) *Handler {
	return &Handler{
		UseCase: usecase,
		Hub:     hub,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws", h.HandleWebSocket)
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsHandler := mywebsocket.NewHandler()
	conn, err := wsHandler.Upgrade(w, r)
	if err != nil {
		log.Println("Ошибка при обновлении WebSocket:", err)
		return
	}

	username, err := utils.GenerateUsername()
	if err != nil {
		log.Println("Ошибка при генерации имени пользователя:", err)
		return
	}

	h.Hub.Register(conn, username)
	defer h.cleanupConnection(conn, username)

	h.handleUserJoin(username)

	h.sendChatHistory(conn)

	h.Hub.BroadcastUserList()

	h.handleMessages(conn, username)
}

func (h *Handler) cleanupConnection(conn *websocket.Conn, username string) {
	h.Hub.Unregister(conn)
	leaveMsg := domain.Message{
		User:    "System",
		Content: username + " покинул чат",
	}
	if err := h.UseCase.SendMessage(context.Background(), leaveMsg); err != nil {
		log.Println("Ошибка при отправке сообщения об отключении:", err)
	}
	h.Hub.BroadcastUserList()
}

func (h *Handler) handleUserJoin(username string) {
	joinMsg := domain.Message{
		User:    "System",
		Content: username + " присоединился к чату",
	}
	if err := h.UseCase.SendMessage(context.Background(), joinMsg); err != nil {
		log.Println("Ошибка при отправке сообщения о подключении:", err)
	}
}

func (h *Handler) sendChatHistory(conn *websocket.Conn) {
	history, err := h.UseCase.GetHistory(context.Background(), historyMessagesCount)
	if err != nil {
		log.Println("Ошибка при получении истории сообщений:", err)
		return
	}

	historyMsg := domain.BroadcastMessage{
		Type: "history",
		Data: history,
	}
	data, err := json.Marshal(historyMsg)
	if err != nil {
		log.Println("Ошибка при маршалинге истории сообщений:", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Ошибка при отправке истории сообщений клиенту:", err)
	}
}

func (h *Handler) handleMessages(conn *websocket.Conn, username string) {
	for {
		_, msgData, err := conn.ReadMessage()
		if err != nil {
			h.handleConnectionError(err, username)
			break
		}

		h.processMessage(msgData, username)
	}
}

func (h *Handler) handleConnectionError(err error, username string) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		log.Printf("Неожиданная ошибка чтения сообщения: %v", err)
	} else {
		log.Printf("Клиент %s отключился: %v", username, err)
	}
}

func (h *Handler) processMessage(msgData []byte, username string) {
	var msg domain.Message
	if err := json.Unmarshal(msgData, &msg); err != nil {
		log.Println("Ошибка при разборе сообщения:", err)
		return
	}
	msg.User = username

	if err := h.UseCase.SendMessage(context.Background(), msg); err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
