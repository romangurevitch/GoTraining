include tools.mk

.PHONY: all build build-hello clean test test-hello test-basics test-bank test-challenges test-temporal lint fmt bench tidy db-up db-down migrate docker-build-hello docker-run-hello temporal-up temporal-down worker-start generate.mocks help

HELLO_IMAGE ?= hello:latest

all: clean tidy lint build test

clean:
	rm -rf bin/
	go clean -testcache

build: ## Build all binaries (hello, bank-api, bank-cli, worker)
	go build -o ./bin/hello ./cmd/hello/...
	go build -o ./bin/bank-api ./cmd/bank-api/...
	go build -o ./bin/bank-cli ./cmd/bank-cli/...
	go build -o ./bin/worker ./cmd/worker/...

build-hello: ## Build hello world binaries
   # Building production ready executable
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ./bin/hello ./cmd/hello/main.go

test: ## Run all tests
	go test ./...

test-hello:  ## Run all hello world tests
	go test ./cmd/hello/...
	go test ./internal/hello/...

test-basics: ## Run module 2 (Go basics) tests
	go test ./internal/basics/...

test-bank: ## Run module 3 (Go Bank) tests
	go test ./internal/bank/...

test-challenges: ## Run all challenge tests
	go test ./internal/challenges/...

lint: $(GOLANGCI_LINT) ## Run linter
	$(GOLANGCI_LINT) run ./...

fmt: ## Format Go code
	gofmt -w .

bench: ## Run all benchmarks
	go test -bench=. -benchmem ./...

tidy: ## Tidy go.mod dependencies
	go mod tidy

docker-build-hello: ## Build docker image for hello world
	docker build -f ./cmd/hello/Dockerfile -t $(HELLO_IMAGE) .

docker-run-hello: ## Run hello world through docker
	docker run --rm $(HELLO_IMAGE) $(NAME)

db-up: ## Start PostgreSQL database
	docker compose up -d postgres

db-down: ## Stop PostgreSQL database
	docker compose down

temporal-up: ## Start Temporal dev server and WireMock
	docker compose up -d temporal wiremock

temporal-down: ## Stop Temporal and WireMock
	docker compose down temporal wiremock

worker-start: ## Start the Temporal order processing worker
	go run ./cmd/worker/... -config="./config/worker/local/config.yaml"

test-temporal: ## Run module 4 (Temporal) tests
	go test ./internal/temporal/...

generate.mocks: $(MOCKGEN) ## Generate mocks
	$(MOCKGEN) -destination=internal/temporal/activities/mocks/mock_inventory_checker.go \
	           -package=mocks \
	           github.com/romangurevitch/go-training/internal/temporal/activities InventoryChecker

migrate: ## Run SQL migrations (instructions only)
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
