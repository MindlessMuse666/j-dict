# Этап сборки
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для сборки с использованием зеркал
RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/main" > /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/community" >> /etc/apk/repositories && \
    apk update --no-cache && \
    apk add --no-cache git build-base

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

# Этап запуска
FROM alpine:3.20

WORKDIR /app

# Устанавливаем необходимые пакеты с указанием конкретной версии
RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/main" > /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/community" >> /etc/apk/repositories && \
    apk update --no-cache && \
    apk add --no-cache ca-certificates postgresql-client curl su-exec

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

# Копируем скрипт запуска
COPY docker-entrypoint.sh .

# Исправляем окончания строк (CRLF -> LF) и делаем скрипт исполняемым
RUN sed -i 's/\r$//' docker-entrypoint.sh && \
    chmod +x docker-entrypoint.sh

# Создаем пользователя для запуска приложения
RUN adduser -D -g '' appuser

# Создаем директории для логов и загрузок
RUN mkdir -p /app/logs /app/uploads/avatars

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]