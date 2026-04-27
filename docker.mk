# Docker targets

HELLO_IMAGE ?= hello:latest

.PHONY: docker-build-hello docker-run-hello

docker-build-hello: ## Build docker image for hello world
	docker build -f ./internal/basics/toolchain/Dockerfile -t $(HELLO_IMAGE) ./internal/basics/toolchain

docker-run-hello: ## Run hello world through docker. Usage: make docker-run-hello NAME=<name>
	docker run --rm $(HELLO_IMAGE) $(NAME)
