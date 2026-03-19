# Build targets

.PHONY: build build-hello

build: generate ## Build all binaries (hello, bank-api, bank-cli)
	@mkdir -p bin
	go build -o bin/hello ./cmd/hello/main.go
	go build -o bin/bank-api ./cmd/bank-api/main.go
	go build -o bin/bank-cli ./cmd/bank-cli/main.go
	@chmod +x bin/*

build-hello: ## Build hello binary for Linux (production)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ./bin/hello ./cmd/hello/main.go
