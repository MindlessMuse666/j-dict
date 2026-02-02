#!/bin/sh

set -e

echo "Запуск приложения..."

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
./main --migrate

# Запускаем сервер
echo "Запуск сервера..."
exec ./main
