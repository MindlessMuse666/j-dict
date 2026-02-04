.PHONY: run build down logs migrate clean reset-db help frontend-install frontend-dev frontend-build logs-frontend

.DEFAULT_GOAL := help

# Запуск всех сервисов с использованием buildx для лучшей кешируемости
run:
	docker-compose down
	docker-compose up --build --remove-orphans

# Запуск в фоновом режиме
up:
	docker-compose up --build -d

# Сборка образов с использованием buildx
build:
	docker-compose build --parallel

# Остановка всех сервисов
down:
	docker-compose down

# Остановка и удаление томов
down-volumes:
	docker-compose down -v --remove-orphans

# Просмотр логов
logs:
	docker-compose logs -f --tail=100

# Просмотр логов определенного сервиса
logs-backend:
	docker-compose logs -f backend --tail=100

logs-postgres:
	docker-compose logs -f postgres --tail=100

logs-kafka:
	docker-compose logs -f kafka --tail=100

logs-frontend:
	docker-compose logs -f frontend --tail=100

# Выполнение миграций
migrate:
	docker-compose run --rm backend ./main --migrate

# Пересоздание базы данных
reset-db:
	docker-compose down -v --remove-orphans
	docker volume prune -f
	docker-compose up -d postgres
	sleep 10
	docker-compose run --rm backend ./main --migrate
	echo "База данных пересоздана"

# Очистка
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f --volumes
	docker image prune -f

# Проверка статуса сервисов
status:
	docker-compose ps

# Подключение к PostgreSQL через psql
psql:
	docker-compose exec postgres psql -U app_user -d jp_ru_dict

# Открыть pgAdmin
pgadmin:
	open http://localhost:5050 || xdg-open http://localhost:5050 || echo "Откройте http://localhost:5050 в браузере"

# Frontend команды
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# Справка
help:
	@echo "Доступные команды:"
	@echo "  make run           - Запуск всех сервисов"
	@echo "  make up            - Запуск в фоновом режиме"
	@echo "  make build         - Сборка образов"
	@echo "  make down          - Остановка сервисов"
	@echo "  make down-volumes  - Остановка и удаление томов"
	@echo "  make logs          - Просмотр логов"
	@echo "  make logs-<service>- Просмотр логов конкретного сервиса"
	@echo "  make migrate       - Выполнить миграции"
	@echo "  make reset-db      - Полное пересоздание базы данных"
	@echo "  make clean         - Полная очистка"
	@echo "  make status        - Статус сервисов"
	@echo "  make psql          - Подключение к PostgreSQL"
	@echo "  make pgadmin       - Открыть pgAdmin"
	@echo "  make frontend-install - Установка зависимостей фронтенда"
	@echo "  make frontend-dev  - Запуск фронтенда в режиме разработки"
	@echo "  make frontend-build- Сборка фронтенда"
	@echo "  make help          - Показать эту справку"