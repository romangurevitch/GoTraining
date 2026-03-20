# Database targets

.PHONY: db-up db-down migrate

db-up: ## Start PostgreSQL database
	docker compose up -d postgres

db-down: ## Stop PostgreSQL database
	docker compose down postgres

migrate: ## Run SQL migrations
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec -T postgres pg_isready -U gotrainer -d gobank > /dev/null 2>&1; do \
		echo "PostgreSQL is not ready yet, retrying in 1s..."; \
		sleep 1; \
	done
	@echo "Running migrations..."
	@docker compose exec -T postgres sh -c 'set -e; for f in /migration/*.sql; do echo "Applying $$f..."; psql -v ON_ERROR_STOP=1 -h localhost -U gotrainer -d gobank -f "$$f"; done'
