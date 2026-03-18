include tools.mk
include build.mk
include test.mk
include docker.mk
include db.mk

.PHONY: all clean lint fmt tidy help

all: clean tidy lint build test

clean: ## Remove build artifacts and test cache
	rm -rf bin/
	go clean -testcache

lint: $(GOLANGCI_LINT) ## Run linter
	$(GOLANGCI_LINT) run ./...

fmt: ## Format Go code
	gofmt -w .

tidy: ## Tidy go.mod dependencies
	go mod tidy

help: ## Show this help message
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
