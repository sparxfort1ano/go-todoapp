include .env
export


export PROJECT_ROOT=${shell pwd}

.DEFAULT_GOAL := help


env-up: ## Env: Launch the project environment
	@docker compose up -d todoapp-postgres todoapp-redis

env-down: ## Env: Stop the project environment
	@docker compose down todoapp-postgres todoapp-redis

env-cleanup: ## Env: Clear the project environment
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward: ## Env: Start the socat container for port forwarding
	@docker compose up -d port-forwarder

env-port-close: ## Env: Stop the socat container
	@docker compose down port-forwarder

logs-cleanup: ## Env: Delete the log files from out/logs
	@read -p "Очистить все log файлы окружения? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

swagger-gen: ## Env: Generate the actual Swagger specification
	@docker compose run --rm swagger \
		init  \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps: ## Env: View running Docker Compose services
	@docker compose ps



migrate-create: ## PostgreSQL: Create the migrations
	@mkdir -p ${PROJECT_ROOT}/migrations; \
	if [ -z "$(name)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create name=init"; \
		exit 1; \
	fi; \
	docker compose run --rm --user "$$(id -u):$$(id -g)" todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(name)"

migrate-up: ## PostgreSQL: Apply the migrations
	@$(MAKE) migrate-action action=up

migrate-down: ## PostgreSQL: Rollback the migrations
	@$(MAKE) migrate-action action=down

migrate-action: ## PostgreSQL: Call a migrations command
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm --user "$$(id -u):$$(id -g)" todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"



todoapp-run: ## Go: Execute the Go application locally (for local development and testing)
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export REDIS_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy: ## Go: Start the Go application in the Docker Compose service (for deploying)
	@docker compose up -d --build todoapp

todoapp-undeploy: ## Go: Stop the Go application in the Docker Compose service
	@docker compose down todoapp



help: ## Show help for commands
	@echo "=== Help ==="
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)