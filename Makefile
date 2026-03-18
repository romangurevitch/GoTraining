include tools.mk

.PHONY: all build build-hello clean test test-hello test-basics test-bank test-challenges lint fmt bench tidy db-up db-down migrate docker-build-hello docker-run-hello

HELLO_IMAGE ?= hello:latest

all: clean tidy lint build test

clean:
	rm -rf bin/
	go clean -testcache

build: clean
	go build -o ./bin/hello ./cmd/hello/...
	go build -o ./bin/bank-api ./cmd/bank-api/...
	go build -o ./bin/bank-cli ./cmd/bank-cli/...

build-hello: clean
   # Building production ready executable
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ./bin/hello ./cmd/hello/main.go

test:
	go test ./...

test-hello:
	go test ./cmd/hello/...
	go test ./internal/hello/...


test-basics:
	go test ./internal/basics/...

test-bank:
	go test ./internal/bank/...

test-challenges:
	go test ./internal/challenges/...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

fmt:
	gofmt -w .

bench:
	go test -bench=. -benchmem ./...

tidy:
	go mod tidy

docker-build-hello:
	docker build -f ./cmd/hello/Dockerfile -t $(HELLO_IMAGE) .

docker-run-hello:
	docker run --rm $(HELLO_IMAGE) $(NAME)

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate:
	@echo "Migration tool not yet configured. See migration/ directory for SQL files."
	@echo "Recommended: use golang-migrate/migrate or goose."
