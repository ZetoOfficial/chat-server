package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ZetoOfficial/chat-server/internal/chat/domain"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	return &Client{client: client}
}

func (r *Client) Publish(ctx context.Context, channel string, msg domain.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

func (r *Client) AddToHistory(ctx context.Context, msg domain.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, "chat_history", data).Err()
}

func (r *Client) GetHistory(ctx context.Context, count int64) ([]domain.Message, error) {
	data, err := r.client.LRange(ctx, "chat_history", 0, count-1).Result()
	if err != nil {
		return nil, err
	}
	var messages []domain.Message
	for i := len(data) - 1; i >= 0; i-- {
		var msg domain.Message
		if err := json.Unmarshal([]byte(data[i]), &msg); err != nil {
			log.Println("Ошибка при разборе истории сообщений:", err)
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
