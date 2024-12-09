package usecase

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ZetoOfficial/chat-server/internal/chat/domain"
	"github.com/ZetoOfficial/chat-server/internal/redis"
	"github.com/ZetoOfficial/chat-server/internal/websocket"
)

type ChatUseCase struct {
	redisRepo *redis.Client
	hub       *websocket.Hub
}

func NewChatUseCase(redisRepo *redis.Client, hub *websocket.Hub) *ChatUseCase {
	return &ChatUseCase{
		redisRepo: redisRepo,
		hub:       hub,
	}
}

func (c *ChatUseCase) Run(ctx context.Context) {
	pubsub := c.redisRepo.Subscribe(ctx, "chat_messages")
	go func() {
		defer pubsub.Close()
		for {
			msgi, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("Ошибка при получении сообщения из Redis:", err)
				return
			}
			var msg domain.Message
			if err := json.Unmarshal([]byte(msgi.Payload), &msg); err != nil {
				log.Println("Ошибка при разборе сообщения:", err)
				continue
			}
			c.hub.BroadcastMessage(msg)
			// Обновление списка пользователей, если сообщение от System
			if msg.User == "System" {
				c.hub.BroadcastUserList()
			}
		}
	}()
}

func (c *ChatUseCase) SendMessage(ctx context.Context, msg domain.Message) error {
	if err := c.redisRepo.Publish(ctx, "chat_messages", msg); err != nil {
		return err
	}
	if err := c.redisRepo.AddToHistory(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (c *ChatUseCase) GetHistory(ctx context.Context, count int64) ([]domain.Message, error) {
	return c.redisRepo.GetHistory(ctx, count)
}
