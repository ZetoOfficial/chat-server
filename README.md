# WebSocket Chat App

Это простой чат с написанный на Go, Redis и WebSocket. Проект позволяет нескольким пользователям общаться в режиме реального времени.

## Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/ZetoOfficial/chat-server.git
   cd chat-server
   ```

2. Установите зависимости:

   ```bash
   go mod download
   ```

---

## Сборка и запуск

### Запуск Parser

1. Соберите бинарный файл:

   ```bash
   go build -o server cmd/server/main.go
   ```

2. Запустите server:

   ```bash
   ./server
   ```

---

## Запуск Redis в Docker

1. Запустите Redis с помощью Docker:

   ```bash
   docker run --name redis-chat -d -p 6379:6379 redis
   ```

   **Описание команды**:

   - `--name redis-chat` — задает имя контейнера.
   - `-d` — запускает контейнер в фоновом режиме.
   - `-p 6379:6379` — пробрасывает порт 6379, чтобы Redis был доступен локально.
