package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myHttp "github.com/ZetoOfficial/chat-server/internal/chat/transport/http"
	"github.com/ZetoOfficial/chat-server/internal/chat/usecase"
	"github.com/ZetoOfficial/chat-server/internal/redis"
	"github.com/ZetoOfficial/chat-server/internal/websocket"
	"github.com/gorilla/mux"
)

const (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 60 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisClient := redis.NewRedisClient("localhost:6379")
	hub := websocket.NewHub()
	chatUseCase := usecase.NewChatUseCase(redisClient, hub)
	chatUseCase.Run(ctx)
	chatHandler := myHttp.NewHandler(chatUseCase, hub)

	router := mux.NewRouter()
	chatHandler.RegisterRoutes(router)

	staticDir := "./frontend"
	fs := http.FileServer(http.Dir(staticDir))
	router.PathPrefix("/").Handler(fs)

	server := &http.Server{
		Addr:         ":8090",
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер запущен на %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка запуска сервера: %v", err)
			cancel()
		}
	}()

	sig := <-stop
	log.Printf("Получен сигнал завершения (%s), остановка сервера...", sig)

	ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Ошибка при завершении сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен")
}
