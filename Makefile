.DEFAULT_GOAL := help

help: ## Показать список всех доступных команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Запуск всех тестов
test: test-unit test-integration ## Запустить все тесты

# Unit тесты (не требуют запущенной БД)
test-unit: ## Unit тесты (не требуют запущенной БД)
	cd backend && go test ./handlers -v

# Интеграционные тесты (требуют запущенную БД)
test-integration: ## Интеграционные тесты (требуют запущенной БД)
	cd backend && go test -v -run "TestFullCRUDWorkflow|TestHealthEndpoint|TestInvalidRequests"

# Сборка проекта
build: ## Сборка проекта
	docker compose build

# Запуск проекта
run: ## Запуск проекта
	docker compose up

# Запуск проекта в фоновом режиме
run-detached: ## Запуск проекта в фоновом режиме
	docker compose up -d

# Остановка проекта
stop: ## Остановка проекта
	docker compose down

# Очистка
clean: ## Очистка
	docker compose down -v
	docker system prune -f

# Тестирование API через curl
test-api: ## Тестирование API через curl
	chmod +x test_api.sh && ./test_api.sh

# Просмотр логов
logs: ## Просмотр логов
	docker compose logs -f

# Просмотр логов конкретного сервиса
logs-backend: ## Просмотр логов backend
	docker compose logs -f backend

logs-frontend: ## Просмотр логов frontend
	docker compose logs -f frontend

logs-db: ## Просмотр логов db
	docker compose logs -f db

logs-nginx: ## Просмотр логов nginx
	docker compose logs -f nginx 