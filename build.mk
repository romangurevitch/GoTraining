# Build targets

.PHONY: build

build: generate ## Build all binaries (bank-server, bank-cli, temporal-client, temporal-worker)
	@mkdir -p bin
	go build -o bin/bank-server ./cmd/bank/server/main.go
	go build -o bin/bank-cli ./cmd/bank/cli/main.go
	go build -o bin/temporal-client ./cmd/temporal/client/main.go
	go build -o bin/temporal-worker ./cmd/temporal/worker/main.go
	@chmod +x bin/*
