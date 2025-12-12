#!/bin/sh

set -e

echo "Ожидание PostgreSQL..."

# Ожидаем, пока PostgreSQL будет готов принимать подключения
while ! nc -z $DB_HOST $DB_PORT; do
  sleep 0.1
done

echo "PostgreSQL доступен"

# Выполняем миграции
echo "Выполнение миграций..."
./main --migrate

# Запускаем приложение
echo "Запуск сервера..."
exec ./main