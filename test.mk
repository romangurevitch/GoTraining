# Test targets

.PHONY: test test-basics test-bank test-challenges bench

test: generate ## Run all tests (excluding student challenges)
	go test $$(go list ./... | grep -v "/internal/challenges/basics")

test-basics: generate ## Run module 2 (Go basics) tests
	go test ./internal/basics/...

test-bank: generate ## Run module 3 (Go Bank) tests
	go test ./internal/bank/...

test-challenges: generate ## Run all basics challenge tests
	go test ./internal/challenges/basics/...

bench: generate ## Run all benchmarks
	go test -bench=. -benchmem ./...
