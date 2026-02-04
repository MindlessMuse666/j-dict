<div align="center">
    <img src="./public/logo.svg" width="120" height="120" alt="Logo" />
    <h1>Фронтенд проекта <a href="https://github.com/MindlessMuse666/jp-ru-dict/tree/main/frontend">j~dict!^w^</a></h1>
    <p><b><i>Клиентская часть на Vue 3 (p≧w≦q)</i></b></p>
    <br>

![Vue.js](https://img.shields.io/badge/vue-%2335495e.svg?style=for-the-badge&logo=vuedotjs&logoColor=%234FC08D) ![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E) ![Vite](https://img.shields.io/badge/vite-%23646CFF.svg?style=for-the-badge&logo=vite&logoColor=yellow) ![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white)
</div>

## 🚀 Запуск

### Требования

- Node.js 18+
- Запущенный бэкенд на порту 8080 (см. [Бэкенд проекта j~dict!\^w\^](../backend/README.md "Бэкенд проекта j~dict!\^w\^")).

### Команды

```bash
# Установка зависимостей
npm install

# Запуск в режиме разработки
npm run dev

# Сборка для продакшена
npm run build

# Предпросмотр сборки
npm run preview
```

## 🏗 Архитектура

| Область | Технология |
| --- | --- |
| **Фреймворк** | Vue 3 |
| **Сборщик** | Vite |
| **Управление состоянием** | Pinia |
| **Навигация** | Vue Router |
| **Стилизация** | TailwindCSS |
| **HTTP-клиент** | Axios |

## 🎨 Особенности UI

- Адаптивность (Mobile/Desktop);
- Кастомные [кошачьи аватарки](./public/assets/default_avatars "аватарки котиков");
- Японская типографика.

## 📦 Сборка Docker

```bash
docker build -t jp-ru-dict-frontend .
docker run -p 5173:80 jp-ru-dict-frontend
```

---

<div align="center">
    <a href="#jdictw">
        <img src="./public/logo.svg" alt="Logo" width="100" height="100">
    </a>
    <br>
    <sub><b>Веб-приложение // j~dict!^w^</b></sub>
    <br>
    <sup><i>Made with love by <a href="https://github.com/MindlessMuse666" target="_blank" title="MindlessMuse666">MindlessMuse666</a></i></sup>
</div>
