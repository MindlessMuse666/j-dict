# j~dict!^w^

![Vue.js](https://img.shields.io/badge/vuejs-%2335495e.svg?style=for-the-badge&logo=vuedotjs&logoColor=%234FC08D)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-000?style=for-the-badge&logo=apachekafka)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E)
![Vite](https://img.shields.io/badge/vite-%23646CFF.svg?style=for-the-badge&logo=vite&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Makefile](https://img.shields.io/badge/Makefile-GNU-blue?style=for-the-badge)

Японско-русский словарь с возможностью поиска, добавления слов и управления личным кабинетом.

## 🏗 Архитектура

```mermaid
graph TD
    User((Пользователь))
    
    subgraph Frontend [Frontend (Vue.js)]
        UI[Интерфейс]
        AuthStore[Auth Store]
        WordStore[Word Store]
    end
    
    subgraph Backend [Backend (Go)]
        API[REST API (Gin)]
        AuthService[Auth Service]
        WordService[Word Service]
        ImportService[Import Service]
        KafkaProducer[Kafka Producer]
    end
    
    subgraph Infrastructure
        DB[(PostgreSQL)]
        Kafka{Apache Kafka}
        Zookeeper[Zookeeper]
    end
    
    User <--> UI
    UI <--> API
    
    API --> AuthService
    API --> WordService
    API --> ImportService
    
    AuthService --> DB
    WordService --> DB
    ImportService --> KafkaProducer
    KafkaProducer --> Kafka
    Kafka -.-> Zookeeper
```

## 🚀 Запуск проекта

Для запуска проекта необходимы установленные **Docker** и **Docker Compose**.

### Быстрый старт (Makefile)

Самый простой способ запустить проект — использовать `make`. Команды выполняются из корня проекта.

```bash
# Запуск всех сервисов (в фоновом режиме, с пересборкой)
make up

# Просмотр логов всех сервисов
make logs

# Остановка сервисов
make down
```

### Ручной запуск (Docker Compose)

Если у вас нет `make`, используйте команды docker-compose напрямую:

```bash
docker-compose up --build -d
```

Приложение будет доступно по адресу: [http://localhost:5173](http://localhost:5173) (или порт 80 в зависимости от конфигурации Docker).

## 📚 API Документация

API задокументировано с помощью Swagger. После запуска бэкенда интерактивная документация доступна по адресу:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Основные эндпоинты:

- **Auth**:
  - `POST /api/auth/register` - Регистрация
  - `POST /api/auth/login` - Вход
  - `GET /api/auth/me` - Получение текущего пользователя
- **Words**:
  - `GET /api/words` - Получение списка слов
  - `POST /api/words` - Добавление слова
  - `PATCH /api/words/{id}` - Обновление слова
  - `DELETE /api/words/{id}` - Удаление слова
- **Import**:
  - `POST /api/words/import` - Загрузка слов из CSV через Kafka

## 🛠 Технологический стек

- **Frontend**: Vue 3, Pinia, Vue Router, TailwindCSS, Vite
- **Backend**: Go (Golang), Gin Gonic, JWT Auth
- **Database**: PostgreSQL 18
- **Message Broker**: Apache Kafka (для асинхронной обработки импорта)
- **DevOps**: Docker, Docker Compose, Makefile

## 📂 Структура проекта

```
.
├── backend/            # Исходный код бэкенда (Go)
│   ├── cmd/            # Точка входа (main.go)
│   ├── internal/       # Бизнес-логика, хендлеры, сервисы
│   ├── db/migrations/  # SQL миграции
│   └── docs/           # Swagger документация
├── frontend/           # Исходный код фронтенда (Vue.js)
│   ├── src/            # Компоненты, сторы, представления
│   └── dist/           # Статическая сборка (для GitHub Pages)
├── docker-compose.yml  # Описание сервисов и их связей
└── Makefile            # Команды для управления проектом
```

## 🌐 GitHub Pages (Демонстрация)

В папке `frontend/dist` находится сгенерированная статическая версия фронтенда, которую можно использовать для демонстрации внешнего вида (UI) на GitHub Pages. Обратите внимание, что без запущенного локально бэкенда функционал API (авторизация, загрузка слов) работать не будет.
