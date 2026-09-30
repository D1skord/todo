include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d todoapp-postgres
env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "This will remove all data in the DB. Are you sure? (y/N) " ans; \
		if [ "$$ans" = "y" ]; then \
			docker compose down todoapp-postgres port-forwarder && \
			sudo rm -rf out/_pgdata; \
				echo "Environment cleaned up."; \
		else \
			echo "Aborting..."; \
		fi

migrate-create:
	@if [ -z "$$seq" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create seq=name"; \
		exit 1; \
	fi;

	@docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$$action" ]; then \
		echo "Error: Migration action is required. Usage: make migrate-up/down action=name"; \
		exit 1; \
	fi;

	@docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/log && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go