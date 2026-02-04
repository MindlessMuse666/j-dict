<div align="center">
    <img src="./frontend/public/logo.svg" width="120" height="120" alt="Logo" />
    <h1>Отчет о разработке информационной системы</h1>
    <h2>"Японский Словарь" (j~dict!^w^)</h2>
    <p><b><i>Курсовой проект /ᐠ｡ꞈ｡ᐟ\</i></b></p>
</div>

## Введение

В рамках производственной практики была спроектирована и реализована информационная система (ИС) для ведения личного японского словаря. Система предназначена для изучающих японский язык и позволяет создавать, редактировать, удалять и искать слова с поддержкой различных чтений (онъёми, кунъёми) и примеров использования.

Целью проекта являлось создание веб-приложения с использованием микросервисных паттернов (в упрощенном виде) и контейнеризации.

## 1. Архитектура и технологии

### 1.1. Технологический стек

Для реализации ИС был выбран следующий стек технологий:

* **Бэкенд:** Go 1.25.5 - для реализации REST API;
* **Фронтенд:** Vue.js 3 + Vite + TailwindCSS - для создания UI;
* **База данных:** PostgreSQL 18 - основное реляционное хранилище данных;
* **Брокер сообщений:** Apache Kafka - для асинхронной обработки тяжелых задач (импорт словарей);
* **Деплой:** Docker и docker-compose - для оркестрации сервисов;
* **Документация:** Swagger (OpenAPI 3.0) - для документирования API.

### 1.2. Архитектура приложения

ИС построена по клиент-серверной архитектуре. Клиентская часть (SPA) взаимодействует с сервером через REST API.

#### Диаграмма контейнеров

Взаимодействие основных частей системы:

1. **Пользователь** работает с интерфейсом (**Frontend**).
2. **Frontend** отправляет запросы к **Backend** (REST API).
3. **Backend** сохраняет/читает данные из **PostgreSQL** и отправляет задачи на импорт в **Kafka**.

```mermaid
graph LR
    User((Пользователь))

    subgraph FE["Frontend"]
        UI["Пользовательский интерфейс<br/>Vue 3"]
        AuthStore["Хранилище авторизации<br/>Pinia"]
        WordStore["Хранилище слов<br/>Pinia"]
    end

    subgraph BE["Backend (Go)"]
        API["REST API<br/>Gin"]
        AuthService["Сервис авторизации"]
        WordService["Сервис словаря"]
        ImportService["Сервис импорта"]
        KafkaProducer["Продюсер Kafka"]
    end

    subgraph INF["Инфраструктура"]
        DB[(PostgreSQL)]
        Kafka{{Apache Kafka}}
        Zookeeper["Zookeeper"]
    end

    User <--> UI
    UI --> API

    API --> AuthService
    API --> WordService
    API --> ImportService

    AuthService --> DB
    WordService --> DB
    ImportService --> KafkaProducer
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

#### Диаграмма компонентов

Детальное устройство системы изнутри:

* **Frontend:** Состоит из Vue-компонентов (отображение), Pinia (хранение данных в памяти) и Axios (общение с сервером).
* **Backend:** Разделен на слои:
  * *Handlers* (обработчики) - принимают HTTP-запросы.
  * *Services* (сервисы) - содержат бизнес-логику (правила работы).
  * *Repositories* (репозитории) - работают напрямую с базой данных.
* **Database & Kafka:** Внешние системы для хранения данных и очереди задач.

```mermaid
classDiagram
    direction LR

    class Frontend {
        Компоненты интерфейса
        Управление состоянием
        HTTP-запросы
    }

    class Backend {
        REST API
        Бизнес-логика
        Доступ к данным
    }

    class Database {
        Пользователи
        Слова
    }

    class Kafka {
        Очередь импорта слов
    }

    Frontend ..> Backend : HTTP / JSON
    Backend ..> Database : SQL
    Backend ..> Kafka : Асинхронные задачи

    classDef frontend fill:#E3F2FD,stroke:#1E88E5,stroke-width:2px
    classDef backend fill:#E8F5E9,stroke:#43A047,stroke-width:2px
    classDef infra fill:#FFFDE7,stroke:#F9A825,stroke-width:2px
```

## 2. Реализация функционала

### 2.1. База данных

Используется PostgreSQL. Реализована схема данных с пользователями и словами. Для слов используются массивы (`text[]`) для хранения вариантов перевода, чтений и тегов, что позволяет гибко искать по любому вхождению.

<div align="center">
    <img src="./docs/screenshot/PgAdmin_DB_Scheme.png" width="auto" height="auto" alt="Схема БД в PostgreSQL" />
    <p><i>Рис. 1. Схема базы данных в PostgreSQL</i></p>
</div>

### 2.2. Бэкенд (API)

Реализован RESTful API на Go с использованием фреймворка Gin.

Основные возможности API:

* **Auth**: Регистрация и аутентификация через JWT.
* **Words**: CRUD-операции со словами, курсорная пагинация, сложный поиск.
* **System**: Health-checks, метрики.

Документация API доступна через Swagger UI:
> Ссылка: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

<div align="center">
    <img src="./docs/screenshot/SwaggerPage.png" width="auto" height="auto" alt="Swagger UI" />
    <p><i>Рис. 2. Интерактивная документация Swagger UI</i></p>
</div>

### 2.3. Фронтенд

Интерфейс разработан на Vue 3 (Composition API). Реализована адаптивная верстка, модальные окна для редактирования, "умный" поиск с дебаунсом.

> Локальная версия: [http://localhost:5173](http://localhost:5173)
> Демонстрация (GitHub Pages): [https://mindlessmuse666.github.io/jp-ru-dict/](https://mindlessmuse666.github.io/jp-ru-dict/)

<div align="center">
    <img src="./docs/screenshot/MainPage.png" width="auto" height="auto" alt="Главная страница" />
    <p><i>Рис. 3. Главная страница приложения</i></p>
    <br>
    <img src="./docs/screenshot/Modal_NewWord.png" width="auto" height="auto" alt="Модальное окно: Добавление нового слова" />
    <p><i>Рис. 4. Модальное окно: Добавление нового слова</i></p>
    <br>
    <img src="./docs/screenshot/Modal_EditWord.png" width="auto" height="auto" alt="Модальное окно: Редактирование слова" />
    <p><i>Рис. 5. Модальное окно: Редактирование слова</i></p>
    <br>
    <img src="./docs/screenshot/Modal_DeleteWord.png" width="auto" height="auto" alt="Модальное окно: Удаление слова" />
    <p><i>Рис. 6. Модальное окно: Удаление слова</i></p>
</div>

## 3. Инструменты разработки и администрирования

Для управления проектом используется `Makefile`, автоматизирующий рутинные задачи (запуск, миграции, логи).

Для администрирования БД подключен **pgAdmin 4**:
> Ссылка: [http://localhost:5050](http://localhost:5050)

## Заключение

В ходе работы была успешно разработана и протестирована ИС личного словаря. Проект демонстрирует использование подходов к проектировеб-разработке и готов к дальнейшему расширению.

---

<div align="center">
    <a href="#введение">
        <img src="./frontend/public/logo.svg" alt="Logo" width="100" height="100">
    </a>
    <br>
    <sub><b>Веб-приложение // j~dict!^w^</b></sub>
    <br>
    <sup><i>Made with love by <a href="https://github.com/MindlessMuse666" target="_blank" title="MindlessMuse666">MindlessMuse666</a></i></sup>
</div>
