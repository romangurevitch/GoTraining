# Docker targets

HELLO_IMAGE ?= hello:latest

.PHONY: docker-build-hello docker-run-hello

docker-build-hello: ## Build docker image for hello world
	docker build -f ./cmd/hello/Dockerfile -t $(HELLO_IMAGE) .

docker-run-hello: ## Run hello world through docker
	docker run --rm $(HELLO_IMAGE) $(NAME)
