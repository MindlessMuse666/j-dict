# Этап сборки
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git gcc musl-dev

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

# Этап запуска
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk add --no-cache ca-certificates postgresql-client

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

# Копируем скрипт запуска
COPY docker-entrypoint.sh .

# Делаем скрипт исполняемым
RUN chmod +x docker-entrypoint.sh

# Создаем пользователя для запуска приложения
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]