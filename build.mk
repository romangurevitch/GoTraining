# Build targets

.PHONY: build build-hello

build: ## Build all binaries (hello, bank-api, bank-cli)
	go build -o ./bin/hello ./cmd/hello/...
	go build -o ./bin/bank-api ./cmd/bank-api/...
	go build -o ./bin/bank-cli ./cmd/bank-cli/...

build-hello: ## Build hello binary for Linux (production)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o ./bin/hello ./cmd/hello/main.go
