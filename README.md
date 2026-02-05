<div align="center">
    <img src="./docs/assets/logo.svg" width="200" height="200" alt="logo.svg" />
    <h1>Документация проекта <a href="https://github.com/MindlessMuse666/jp-ru-dict/blob/main/README.md">j~dict!^w^</a></h1>
    <p><b><i>Личный японский словарь (∩^o^)⊃━☆</i></b></p>
    <br>

![Vue.js](https://img.shields.io/badge/vue-%2335495e.svg?style=for-the-badge&logo=vuedotjs&logoColor=%234FC08D) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E) ![Vite](https://img.shields.io/badge/vite-%23646CFF.svg?style=for-the-badge&logo=vite&logoColor=yellow) ![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) <br> ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Apache Kafka](https://img.shields.io/badge/Kafka-000?style=for-the-badge&logo=apachekafka) ![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white) ![Makefile](https://img.shields.io/badge/Makefile-%23404d59.svg?style=for-the-badge&logo=make&logoColor=red)
</div>

## 📖 [j~dict!\^w\^](https://github.com/MindlessMuse666/jp-ru-dict/blob/main/README.md) это...

...небольшое веб-приложение для ведения личного японского словаря. В нём вы можете создавать, читать, обновлять и удалять слова, а также организовывать слова по тегам и даже... управлять личным кабинетом и настраивать свою аватарку!😜

Этот пет-проект резюмирует практические навыки, полученные в процессе прохождения производстенной практики от [Колледжа связи №54 им. П. М. Вострухина](https://www.ks54.ru "Колледж связи №54 им. П. М. Вострухина") в организации [Инфотекс Интернет Траст](https://iitrust.ru "Инфотекс Интернет Траст").

## 🏗 Архитектура

### Диаграмма контейнеров

Основные части системы взаимодействуют следующим образом:

1. **Пользователь** работает с интерфейсом (**Фронтенда**);
2. **Фронтенд** отправляет запросы к **Бэкенду** (REST API);
3. **Бэкенд** сохраняет/читает данные из **PostgreSQL** и отправляет события изменений в **Kafka**.

```mermaid
graph LR
    User((Пользователь))

    subgraph FE["Фронтенд"]
        UI["Пользовательский интерфейс"]
        AuthStore["Хранилище авторизации"]
        WordStore["Хранилище слов"]
    end

    subgraph BE["Бэкенд"]
        API["REST API"]
        AuthService["Сервис авторизации"]
        WordService["Сервис словаря"]
        ImportService["Сервис импорта"]
        KafkaProducer["Продюсер Kafka"]
    end

    subgraph INF["Инфраструктура"]
        DB[(PostgreSQL)]
        Kafka{{Kafka}}
        Zookeeper["Zookeeper"]
    end

    User <--> UI
    UI --> API

    API --> AuthService
    API --> WordService
    API --> ImportService

    AuthService --> DB
    WordService --> DB
    ImportService --> WordService
    WordService --> KafkaProducer
    KafkaProducer --> Kafka
    Kafka -.-> Zookeeper

    classDef frontend fill:#E3F2FD,stroke:#1E88E5,stroke-width:2px
    classDef backend fill:#E8F5E9,stroke:#43A047,stroke-width:2px
    classDef infra fill:#FFFDE7,stroke:#F9A825,stroke-width:2px
    classDef user fill:#FCE4EC,stroke:#C2185B,stroke-width:2px

    class UI,AuthStore,WordStore frontend
    class API,AuthService,WordService,ImportService,KafkaProducer backend
    class DB,Kafka,Zookeeper infra
    class User user
```

### Диаграмма компонентов

Изнутри ИС настроена следующим образом:

* **Фронтенд:** Состоит из Vue-компонентов (UI), Pinia (хранение данных в памяти) и Axios (общение с сервером);
* **Бэкенд:** Разделен на слои:
  * *Handlers* (обработчики) - принимают HTTP-запросы;
  * *Services* (сервисы) - содержат бизнес-логику (правила работы);
  * *Repositories* (репозитории) - работают напрямую с базой данных.
* **База данных и Kafka:** Внешние системы для хранения данных и очереди событий.

```mermaid
classDiagram
    direction LR

    class Фронтенд {
        Компоненты интерфейса
        Управление состоянием
        HTTP-запросы
    }

    class Бэкенд {
        REST API
        Бизнес-логика
        Доступ к данным
    }

    class База данных {
        Пользователи
        Слова
    }

    class Kafka {
        События изменений слов
    }

    Фронтенд ..> Бэкенд : HTTP / JSON
    Бэкенд ..> PostgreSQL : SQL
    Бэкенд ..> Kafka : Публикация событий

    classDef frontend fill:#E3F2FD,stroke:#1E88E5,stroke-width:2px
    classDef backend fill:#E8F5E9,stroke:#43A047,stroke-width:2px
    classDef infra fill:#FFFDE7,stroke:#F9A825,stroke-width:2px
```

## 🚀 Запуск проекта

Для запуска проекта необходимы установленные **Docker** и **docker-compose**.

### Быстрый старт через Makefile

Самый простой способ запустить проект - использовать `make`. Команды выполняются из корня проекта.

```bash
# Запуск всех сервисов (в фоновом режиме, с пересборкой)
make up

# Просмотр логов всех сервисов
make logs

# Остановка сервисов
make down
```

### Ручной запуск через [docker-compose.yml](./docker-compose "docker-compose")

Используйте команды `docker-compose` из корня проекта напрямую:

```bash
docker-compose up --build -d
```

Приложение будет доступно по адресу: [http://localhost:5173](http://localhost:5173 "Развёрнутое через docker-compose веб-приложение j~dict!^w^").

## 📚 API-Документация

API задокументировано с помощью [Swagger](./openapi.yaml "Swagger API-Документация"). После запуска бэкенда Swagger-документация доступна по адресу:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html "Swagger API-Документация")

### Основные эндпоинты

* [Auth](./backend/internal/handler/auth.go "AuthHandler"):
  * `POST /api/auth/register` - Регистрация
  * `POST /api/auth/login` - Вход
  * `GET /api/auth/me` - Получение текущего пользователя
* [Words](./backend/internal/handler/word.go "WordHandler"):
  * `GET /api/words` - Получение списка слов
  * `POST /api/words` - Добавление слова
  * `PATCH /api/words/{id}` - Обновление слова
  * `DELETE /api/words/{id}` - Удаление слова
* [Import](./backend/internal/handler/import.go "ImportHandler"):
  * `POST /api/words/import` - Загрузка слов из CSV через Kafka

## 🛠 Технологический стек

| Область | Технология |
| ---- | --- |
| **Бэкенд** | Go (Golang), Gin Gonic, JWT-Auth |
| **Фронтенд** | Vue 3, Pinia, Vue Router, TailwindCSS, Vite |
| **База данных** | PostgreSQL 18 |
| **Брокер сообщений** | Apache Kafka (для асинхронной обработки импорта) |
| **Инфраструктура** | Docker, docker-compose, Makefile |
| **Документация** | Swagger (OpenAPI) |

## 📂 Структура проекта

```text
├── backend/            # Код бэкенда
│   ├── cmd/            # Точка входа (main.go)
│   ├── internal/       # Бизнес-логика, хендлеры, сервисы
│   ├── db/migrations/  # SQL-миграции
│   └── docs/           # Swagger-документация
├── frontend/           # Код фронтенда
│   ├── src/            # Компоненты, сторы, представления
├── docker-compose.yml  # Описание сервисов и их связей
└── Makefile            # Команды для управления проектом
```

## 🌐 GitHub Pages (Демо)

В этой ветке находится статическая версия фронтенда, которая используется для демо интрефейса проекта на GitHub Pages.

**ВАЖНО!** Без запущенного локально бэкенда функционал API (авторизация, загрузка слов) работать **не будет**.

---

<div align="center">
    <a href="#jdictw">
        <img src="./docs/assets/logo.svg" alt="Logo" width="200" height="200">
    </a>
    <br>
    <sub><b>Веб-приложение // j~dict!^w^</b></sub>
    <br>
    <sup><i>Made with love by <a href="https://github.com/MindlessMuse666" target="_blank" title="MindlessMuse666">MindlessMuse666</a></i></sup>
</div>
