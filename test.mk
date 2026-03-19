# Test targets

.PHONY: test test-hello test-basics test-bank test-challenges bench

test: generate ## Run all tests (excluding student challenges)
	go test $$(go list ./... | grep -v "/internal/basics/challenges")

test-hello: ## Run all hello world tests
	go test ./cmd/hello/... ./internal/hello/...

test-basics: ## Run module 2 (Go basics) tests (excluding challenges)
	go test $$(go list ./internal/basics/... | grep -v "/internal/basics/challenges")

test-bank: generate ## Run module 3 (Go Bank) tests
	go test ./internal/bank/...

test-challenges: ## Run all basics challenge tests
	go test ./internal/basics/challenges/...

bench: ## Run all benchmarks
	go test -bench=. -benchmem ./...
