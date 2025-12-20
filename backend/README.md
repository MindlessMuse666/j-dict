# Бэкенд проекта jp-ru-dict (✿◠‿◠)

Бэкенд часть приложения для личного русско-японского словаря.

## Инструкция по запуску

### Требования

- Docker и Docker Compose
- Go 1.25.5+ (для локальной разработки)

### Запуск через Docker

```bash
# Клонирование репозитория
git clone https://github.com/MindlessMuse666/jp-ru-dict.git
cd jp-ru-dict

# Настройка окружения
cp .env.example .env
# Отредактируйте .env при необходимости

# Запуск всех сервисов
docker-compose up --build -d

# Проверка статуса
docker-compose ps
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

## Структура проекта

```md
backend/
├── cmd/server/             # Точка входа приложения
├── internal/               # Внутренние пакеты
│   ├── config/             # Конфигурация (через переменные окружения)
│   ├── handler/            # HTTP-обработчики (REST API)
│   ├── middleware/         # JWT аутентификация, CORS
│   ├── model/              # Модели данных
│   ├── repository/         # Слой доступа к данным (PostgreSQL)
│   ├── service/            # Бизнес-логика
│   └── kafka/              # Интеграция с Kafka
├── pkg/logger/             # Структурированное логирование (zerolog)
├── db/migrations/          # Миграции базы данных
├── docs/                   # Swagger-документация
├── Dockerfile              # Docker образ бэкенда
└── docker-entrypoint.sh    # Скрипт запуска в контейнере
```

## Модель базы данных

### Таблица `users`

| Поле | Тип | Описание |
| --- | --- | --- |
| **id** | SERIAL | Уникальный идентификатор |
| **email** | VARCHAR(255) | Email пользователя (уникальный) |
| **password_hash** | VARCHAR(255) | Хэш пароля (bcrypt) |
| **name** | VARCHAR(100) | Имя пользователя |
| **created_at** | TIMESTAMPTZ | Время регистрации |
| **updated_at** | TIMESTAMPTZ | Время обновления |

### Таблица `words`

| Поле | Тип Описание |
| --- | --- | --- |
| **id** | SERIAL | Уникальный идентификатор |
| **user_id** | INTEGER | Владелец слова (FK к `users.id`) |
| **jp** | TEXT[] | Написания на японском (обязательно) |
| **ru** | TEXT[] | Переводы на русский (обязательно) |
| **on** | TEXT[] | Онъёми чтения (опционально) |
| **kun** | TEXT[] | Кунъёми чтения (опционально) |
| **ex_jp** | TEXT[] | Примеры на японском (макс. 3) |
| **ex_ru** | TEXT[] | Переводы примеров (макс. 3) |
| **tags** | TEXT[] | Теги для категоризации (макс. 5) |
| **created_at** | TIMESTAMPTZ | Время создания |
| **updated_at** | TIMESTAMPTZ | Время обновления |

**Индексы:**

- `idx_words_user_id` - поиск по пользователю;
- `idx_words_jp`, `idx_words_ru` - GIN-индексы для массивов;
- `idx_words_tags` - поиск по тегам;
- `idx_words_onyomi`, `idx_words_kunyomi` - частичные индексы для чтений.

## API Endpoints

### Аутентификация

| Эндпоинт | Описание |
| --- | --- |
| `POST /api/auth/register` | регистрация пользователя |
| `POST /api/auth/login` | вход в систему |
| `GET /api/auth/me` | информация о текущем пользователе |
| `POST /api/auth/logout` | выход из системы |

### Управление словами

| Эндпоинт | Описание |
| --- | --- |
| `GET /api/words` | список слов с пагинацией |
| `GET /api/words/{id}` | слово по ID |
| `POST /api/words` | создание слова |
| `PATCH /api/words/{id}` | частичное обновление слова |
| `DELETE /api/words/{id}` | удаление слова |
| `GET /api/words/search` | поиск по различным критериям |

### Системные

| Эндпоинт | Описание |
| --- | --- |
| `GET /api/health` | проверка работоспособности сервиса |
| `GET /swagger/*` | Swagger UI документация |

## Технологический стек

| Технология | Функционал |
| --- | --- |
| **Go 1.25.5** | основной язык бэкенда |
| **Gin** | HTTP-фреймворк |
| **PostgreSQL 18** | основная база данных |
| **Apache Kafka** | асинхронная обработка событий |
| **JWT** | аутентификация |
| **Docker** | контейнеризация |
| **Swagger/OpenAPI 3.0** | документация API |

## Особенности реализации

| Особенность | Описание | Примечание |
| --- | --- | --- |
| **Пагинация** | курсорная пагинация на основе поля `id`: | параметры пагинации:<ul><li>`limit` (default=20, max=100);</li><li>`cursor` (опциональный ID для следующей страницы).</li></ul> |
| **Поиск** | <ol><li>по подстроке: `ILIKE` по массивам `jp` и `ru`;</li><li>по чтениям: точное совпадение элемента массива `on` и/или `kun`;<li>по тегам: точное совпадение элемента массива `tags`;</li></ol> | поиск поддерживает комбинированные фильтры |
| **Валидация** | обязательность и ограничения параметров | <ul><li>обязательные поля: `jp`, `ru`, `user_id`;</li><li>ограничения массивов: примеры (max=3), теги (max=5).</li></ul>для частичного обновления все поля опциональны (см. эндпоинт PATCH) |
| **Безопасность** | JWT-токены для аутентификации | <ul><li>bcrypt для хэширования паролей;</li><li>изоляция данных по `user_id`;</li><li>CORS-политика.</li></ul> |

## Документация API

### Swagger UI

После запуска сервера откройте: <http://localhost:8080/swagger/index.html>
OpenAPI спецификация

- JSON: `docs/swagger.json`
- YAML: `docs/swagger.yaml` или корневой `openapi.yaml`

### Тестирование API

Импортируйте коллекцию Insomnia/Postman:

- Коллекция Insomnia: [тут](docs/insomnia/insomnia-collection.json "Коллекция Insomnia")
- Коллекция Postman: [тут](docs/postman/postman-collection.json "Коллекция Postman")

## Отладка и мониторинг

### Логи

Структурированные логи в формате JSON:

```bash
docker-compose logs backend -f
```

### База данных

```bash
# Подключение через psql
docker-compose exec postgres psql -U app_user -d jp_ru_dict

# PgAdmin (опционально)
open http://localhost:5050
```

### Health checks

```bash
curl http://localhost:8080/api/health
```

## Миграции базы данных

Управление через `golang-migrate`:

```bash
# Применение миграций
go run cmd/server/main.go --migrate

# Через Docker
docker-compose run --rm backend ./main --migrate

# Создание новой миграции
migrate create -ext sql -dir db/migrations -seq migration_name
```

## Конфигурация

Основные переменные окружения (`.env`):

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=password
DB_NAME=jp_ru_dict
DB_SSL_MODE=disable

# Server
SERVER_PORT=8080
LOG_LEVEL=info

# JWT
JWT_SECRET=secret-jwt-key
JWT_EXPIRY=24h

# Kafka
KAFKA_BROKER=kafka:9092

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

## Тестирование

```bash
# Запуск всех тестов
cd backend && go test ./...

# Тестирование с покрытием
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Полезные команды

```bash
# Полная очистка
make clean

# Пересоздание базы данных
make reset-db

# Просмотр логов
make logs-backend
make logs-postgres

# Статус сервисов
make status
```

## Лицензия

Проект под распространяется под лицензией [MIT](../LICENCE)
