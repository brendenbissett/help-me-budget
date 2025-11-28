# Database management commands

# Database connection string
DB_URL=postgres://budgetuser:budgetpass@localhost:5432/help_me_budget?sslmode=disable

.PHONY: help up down migrate-up migrate-down migrate-create migrate-force migrate-version

help:
	@echo "Database Management Commands:"
	@echo "  make up           - Start PostgreSQL and Redis containers"
	@echo "  make down         - Stop and remove containers"
	@echo "  make migrate-up   - Run all pending migrations"
	@echo "  make migrate-down - Rollback last migration"
	@echo "  make migrate-create NAME=<name> - Create new migration files"
	@echo "  make migrate-force VERSION=<n> - Force database to specific version (fixes dirty state)"
	@echo "  make migrate-version - Show current migration version"

up:
	docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database and Redis are ready!"

down:
	docker-compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=my_migration"; \
		exit 1; \
	fi
	migrate create -ext sql -dir migrations -seq $(NAME)

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make migrate-force VERSION=4"; \
		exit 1; \
	fi
	migrate -path migrations -database "$(DB_URL)" force $(VERSION)

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version
