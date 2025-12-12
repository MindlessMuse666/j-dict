.PHONY: run build down logs migrate test clean

# Запуск всех сервисов
run:
	docker-compose up --build

# Сборка образов
build:
	docker-compose build

# Остановка всех сервисов
down:
	docker-compose down

# Просмотр логов
logs:
	docker-compose logs -f

# Выполнение миграций
migrate:
	docker-compose run --rm backend ./main --migrate

# Откат миграций
migrate-down:
	docker-compose run --rm backend migrate -path ./db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@postgres:5432/$(DB_NAME)?sslmode=disable" down

# Запуск тестов
test:
	cd backend && go test ./...

# Полная очистка
clean:
	docker-compose down -v
	docker system prune -f

# Запуск бэкенда локально (для разработки)
dev-backend:
	cd backend && go run cmd/server/main.go

# Установка зависимостей бэкенда
backend-deps:
	cd backend && go mod download

# Пересоздание базы данных
reset-db:
	docker-compose down -v
	docker-compose up -d postgres
	sleep 5
	make migrate