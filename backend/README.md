<div align="center">
    <img src="../frontend/public/logo.svg" width="200" height="200" alt="Logo" />
    <h1>Бэкенд проекта <a href="https://github.com/MindlessMuse666/jp-ru-dict/tree/main/backend">j~dict!^w^</a></h1>
    <p><b><i>Серверная часть на Go (Golang) (✿◠‿◠)</i></b></p>
    <br>

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Apache Kafka](https://img.shields.io/badge/Kafka-000?style=for-the-badge&logo=apachekafka) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) <br> ![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white) ![Makefile](https://img.shields.io/badge/Makefile-%23404d59.svg?style=for-the-badge&logo=make&logoColor=red)
</div>

## 🚀 Запуск

### Через Makefile

Из корня проекта:

```bash
# Запуск всех сервисов
make up

# Просмотр логов
make logs-backend
```

### Локальная разработка

```bash
cd backend

# Установка зависимостей
go mod download

# Запуск зависимостей (PostgreSQL, Kafka)
docker-compose up -d postgres kafka zookeeper

# Выполнение миграций
go run cmd/server/main.go --migrate

# Запуск сервера
go run cmd/server/main.go
```

## 🛠 Стек технологий

| Область | Технология |
| --- | --- |
| **Язык** | Go 1.25.5 |
| **Фреймворк** | Gin Gonic |
| **Базы данных** | PostgreSQL 18 |
| **Брокер сообщений** | Apache Kafka |
| **Логирование** | Zerolog |
| **Миграции** | Golang-migrate |
| **Документация** | Swagger / OpenAPI |

## 📂 Структура

- `cmd/server/` - Точка входа (main.go)
- `internal/` - Приватный код (handlers, services, repositories)
- `pkg/` - Публичные пакеты (logger)
- `docs/` - Swagger-документация

## 🔌 API

Документация доступна по адресу `http://localhost:8080/swagger/index.html` после запуска сервера.

Основные группы эндпоинтов:

- `/api/auth` - Авторизация (JWT);
- `/api/words` - CRUD-операции со словами;
- `/api/health` - Проверка состояния.

---

<div align="center">
    <a href="#-запуск">
        <img src="../frontend/public/logo.svg" alt="Logo" width="200" height="200">
    </a>
    <br>
    <sub><b>Веб-приложение // j~dict!^w^</b></sub>
    <br>
    <sup><i>Made with love by <a href="https://github.com/MindlessMuse666" target="_blank" title="MindlessMuse666">MindlessMuse666</a></i></sup>
</div>
