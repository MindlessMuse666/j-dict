<div align="center">
    <img src="../frontend/public/logo.svg" width="200" height="200" alt="Logo" />
    <h1>Отчет о разработке информационной системы</h1>
    <h2>"Японский Словарь" (j~dict!^w^)</h2>
    <p><b><i>Проект производственной практики /ᐠ｡ꞈ｡ᐟ\</i></b></p>
</div>

## Введение

В рамках производственной практики была спроектирована и реализована информационная система (ИС) для ведения личного японского словаря. Система предназначена для изучающих японский язык и позволяет создавать, редактировать, удалять и искать слова с поддержкой различных чтений (онъёми, кунъёми) и примеров использования.

Целью проекта являлось создание веб-приложения с использованием микросервисных паттернов (в упрощенном виде) и контейнеризации.

## 1. Архитектура и технологии

### 1.1. Стек

Для реализации ИС был выбран следующий стек технологий:

* **Бэкенд:** Go 1.25.5 - для реализации REST API;
* **Фронтенд:** Vue.js 3 + Vite + TailwindCSS - для создания UI;
* **База данных:** PostgreSQL 18 - основное реляционное хранилище данных;
* **Брокер сообщений:** Apache Kafka - для публикации событий изменений данных (CDC-паттерн);
* **Деплой:** Docker и docker-compose - для оркестрации сервисов;
* **Документация:** Swagger (OpenAPI 3.0) - для документирования API.

### 1.2. Архитектура приложения

ИС построена по клиент-серверной архитектуре. Клиентская часть (SPA) взаимодействует с сервером через REST API.

#### Диаграмма контейнеров

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

#### Диаграмма компонентов

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

## 2. Реализация функционала

### 2.1. Бэкенд (API)

Реализован RESTful API на Go с использованием фреймворка Gin.

Основные возможности API:

* **Auth**: Регистрация и аутентификация через JWT.
* **Words**: CRUD-операции со словами, курсорная пагинация, сложный поиск.
* **System**: Health-checks, метрики.

Документация API доступна через Swagger UI:
> Ссылка: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

<div align="center">
    <img src="./screenshot/service/swagger.png" width="auto" height="auto" alt="Swagger UI" />
    <p><i>Рис. 1. Документация Swagger UI</i></p>
</div>

### 2.2. База данных

Используется PostgreSQL. Реализована схема данных с пользователями и словами. Для слов используются массивы (`text[]`) для хранения вариантов перевода, чтений и тегов, что позволяет гибко искать по любому вхождению.

<div align="center">
    <img src="./screenshot/service/pgadmin_scheme.png" width="auto" height="auto" alt="Схема БД в PostgreSQL" />
    <p><i>Рис. 2. Схема базы данных в PostgreSQL</i></p>
</div>

### 2.3. Фронтенд

UI разработан на Vue 3 (Composition API). Реализована адаптивная верстка, модальные окна для редактирования, компактный и подробный режимы и поиск по фильтрам.

> Локальная версия: [http://localhost:5173](http://localhost:5173)
> Демонстрация (GitHub Pages): [https://mindlessmuse666.github.io/jp-ru-dict/](https://mindlessmuse666.github.io/jp-ru-dict/)

#### Превью основного интерфейса

<div align="center">
    <img src="./screenshot/main/basic_view.png" width="auto" height="auto" alt="Главная страница (базовый вид)" />
    <p><i>Рис. 3. Главная страница приложения (базовый вид)</i></p>
    <br>
    <img src="./screenshot/main/full_view.png" width="auto" height="auto" alt="Главная страница (расширенный вид)" />
    <p><i>Рис. 4. Главная страница приложения (расширенный вид)</i></p>
    <br>
    <img src="./screenshot/main/search_view.png" width="auto" height="auto" alt="Главная страница (поиск)" />
    <p><i>Рис. 5. Главная страница приложения (поиск)</i></p>
    <br>
    <img src="./screenshot/profile/profile.png" width="auto" height="auto" alt="Профиль пользователя" />
    <p><i>Рис. 6. Профиль пользователя</i></p>
</div>

#### Превью функционала

<div align="center">
    <img src="./screenshot/main/modal/new_word.png" width="auto" height="auto" alt="Модальное окно: Добавление нового слова" />
    <p><i>Рис. 7. Модальное окно: Добавление нового слова</i></p>
    <br>
    <img src="./screenshot/main/modal/edit_word.png" width="auto" height="auto" alt="Модальное окно: Редактирование слова" />
    <p><i>Рис. 8. Модальное окно: Редактирование слова</i></p>
    <br>
    <img src="./screenshot/main/modal/delete_word.png" width="auto" height="auto" alt="Модальное окно: Удаление слова" />
    <p><i>Рис. 9. Модальное окно: Удаление слова</i></p>
    <br>
    <img src="./screenshot/main/modal/import_word.png" width="auto" height="auto" alt="Модальное окно: Импорт слов" />
    <p><i>Рис. 10. Модальное окно: Импорт слов</i></p>
</div>

### 2.4. Брокер сообщений

Используется Kafka для публикации событий изменений данных (создание, обновление, удаление слов). Это позволяет в будущем подключать новые сервисы (например, аналитику или уведомления) без изменения основного бэкенда.

<div align="center">
    <img src="./screenshot/service/kafka_ui.png" width="auto" height="auto" alt="Kafka UI" />
    <p><i>Рис. 11. Интерфейс Kafka UI</i></p>
</div>

## 3. Инструменты разработки и администрирования

Для управления проектом используется [Makefile](../Makefile "Makefile"), который автоматизирует запуск и просмотр логов.

Для администрирования БД подключен **pgAdmin 4**:
> Ссылка: [http://localhost:5050](http://localhost:5050)

## 4. Небольшая фишка

Пользователь на странице профиля может управлять своей аватаркой!
Есть два интересных сценария и сценарий удаления аватара (сброса до дефолтного).

### Сценарий 1: Загрузка пользовательской аватарки

<div align="center">
    <img src="./screenshot/profile/avatar/avatar_00.png" width="auto" height="auto" alt="Шаг 1: Страница профиля" />
    <p><i>Рис. 12. Шаг 1: Страница профиля</i></p>
    <br>
    <img src="./screenshot/profile/avatar/avatar_load_01.png" width="auto" height="auto" alt="Шаг 2: Окно загрузки аватара" />
    <p><i>Рис. 13. Шаг 2: Окно загрузки аватара</i></p>
    <br>
    <img src="./screenshot/profile/avatar/avatar_load_02.png" width="auto" height="auto" alt="Шаг 3: Обрезка аватара" />
    <p><i>Рис. 14. Шаг 3: Обрезка аватара</i></p>
    <br>
    <img src="./screenshot/profile/avatar/avatar_load_03.png" width="auto" height="auto" alt="Шаг 4: Аватар успешно загружен + Тост" />
    <p><i>Рис. 15. Шаг 4: Аватар успешно загружен + Тост</i></p>
</div>

### Сценарий 2: Установка аватарки из предлагаемых

<div align="center">
    <img src="./screenshot/profile/avatar/avatar_00.png" width="auto" height="auto" alt="Шаг 1: Страница профиля" />
    <p><i>Рис. 16. Шаг 1: Страница профиля</i></p>
    <br>
    <img src="./screenshot/profile/avatar/avatar_preset_load_01.png" width="auto" height="auto" alt="Шаг 2: Окно установки предустановленного аватара" />
    <p><i>Рис. 17. Шаг 2: Окно установки предустановленного аватара</i></p>
    <br>
    <img src="./screenshot/profile/avatar/preset_load_02.png" width="auto" height="auto" alt="Шаг 4: Аватар успешно установлен + Тост" />
    <p><i>Рис. 18. Шаг 4: Аватар успешно установлен + Тост</i></p>
</div>

## Заключение

В ходе работы была успешно разработана и протестирована ИС личного словаря. Проект демонстрирует использование подходов к веб-разработке и готов к дальнейшему расширению.

---

<div align="center">
    <a href="#введение">
        <img src="../frontend/public/logo.svg" alt="Logo" width="200" height="200">
    </a>
    <br>
    <sub><b>Веб-приложение // j~dict!^w^</b></sub>
    <br>
    <sup><i>Made with love by <a href="https://github.com/MindlessMuse666" target="_blank" title="MindlessMuse666">MindlessMuse666</a></i></sup>
</div>
