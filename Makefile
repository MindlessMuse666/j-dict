.PHONY: up down logs clean reset-db help

.DEFAULT_GOAL := help

# Запуск всех сервисов с --no-cache
up:
	docker-compose down
	docker-compose up --build --remove-orphans

# Остановка всех сервисов
down:
	docker-compose down

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

# Открыть pgAdmin
pgadmin:
	open http://localhost:5050 || xdg-open http://localhost:5050 || echo "Откройте http://localhost:5050 в браузере"

# Справка
help:
	@echo "Доступные команды:"
	@echo "  make run           - Запуск всех сервисов"
	@echo "  make up            - Запуск в фоновом режиме"
	@echo "  make down          - Остановка сервисов"
	@echo "  make logs          - Просмотр логов"
	@echo "  make logs-<service>- Просмотр логов конкретного сервиса"
	@echo "  make reset-db      - Полное пересоздание базы данных"
	@echo "  make clean         - Полная очистка"
	@echo "  make status        - Статус сервисов"
	@echo "  make pgadmin       - Открыть pgAdmin"
	@echo "  make help          - Показать эту справку"