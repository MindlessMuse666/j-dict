# Фронтенд проекта j~dict!\^w\^ (p≧w≦q)

![Vue.js](https://img.shields.io/badge/vuejs-%2335495e.svg?style=for-the-badge&logo=vuedotjs&logoColor=%234FC08D) ![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E) ![Vite](https://img.shields.io/badge/vite-%23646CFF.svg?style=for-the-badge&logo=vite&logoColor=white) ![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white)

Фронтенд часть приложения (Vue 3 + Vite).

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
| **HTTP клиент** | Axios |

## 🎨 Особенности UI

- Адаптивность (Mobile/Desktop);
- Кастомные [кошачьи аватарки](./public/assets/default_avatars "аватарки котиков");
- Японская типографика.

## 📦 Сборка Docker

```bash
docker build -t jp-ru-dict-frontend .
docker run -p 5173:80 jp-ru-dict-frontend
```
