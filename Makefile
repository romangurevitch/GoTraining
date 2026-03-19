include tools.mk
include build.mk
include test.mk
include docker.mk
include db.mk
include internal/temporal/temporal.mk

.PHONY: all clean lint fmt tidy run-bank-api run-bank-cli help

all: tools generate clean tidy lint build test

generate: $(MOCKERY) ## Generate mocks
	$(MOCKERY)

clean: ## Remove build artifacts and test cache
	rm -rf bin/
	go clean -testcache

lint: $(GOLANGCI_LINT) ## Run linter
	$(GOLANGCI_LINT) run ./...

fmt: ## Format Go code
	gofmt -w .

tidy: ## Tidy go.mod dependencies
	go mod tidy

# Thin-slice targets
run-bank-api: build ## Start bank API server
	./bin/bank-api

run-bank-cli: build ## Use bank CLI (example: make run-bank-cli ARGS="account create 'John Doe'")
	./bin/bank-cli $(ARGS)
  
help: ## Show this help message
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
