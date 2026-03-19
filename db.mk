# Database targets

.PHONY: db-up db-down migrate

db-up: ## Start PostgreSQL database
	docker compose up -d postgres

db-down: ## Stop PostgreSQL database
	docker compose down

migrate: ## Run SQL migrations (instructions only)
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."
