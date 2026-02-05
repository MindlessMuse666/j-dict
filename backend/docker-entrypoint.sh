#!/bin/sh

set -e

echo "Запуск приложения..."

# Настройка прав доступа (так как тома могут быть примонтированы как root)
echo "Настройка прав доступа..."
mkdir -p /app/uploads/avatars /app/logs
chown -R appuser:appuser /app/uploads /app/logs

# Ожидаем доступности PostgreSQL
echo "Ожидание PostgreSQL..."
while ! nc -z postgres 5432; do
  sleep 0.1
done
echo "PostgreSQL доступен"

# Ожидаем доступности Kafka
echo "Ожидание Kafka..."
while ! nc -z kafka 9092; do
  sleep 0.1
done
echo "Kafka доступен"

# Выполняем миграции
echo "Выполнение миграций..."
su-exec appuser ./main --migrate

# Запускаем сервер
echo "Запуск сервера..."
exec su-exec appuser ./main
