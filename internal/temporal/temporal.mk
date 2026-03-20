# Temporal targets

.PHONY: temporal-up temporal-down worker-start test-temporal generate.mocks workflow-auto workflow-signal workflow-signal-send

temporal-up: ## Start Temporal dev server and WireMock
	docker compose up -d temporal wiremock

temporal-down: ## Stop Temporal and WireMock
	docker compose down temporal wiremock

worker-start: ## Start the Temporal order processing worker
	go run ./cmd/temporal/worker/... -config="./config/temporal/worker/local/config.yaml"

test-temporal: ## Run module 4 (Temporal) tests
	go test ./internal/temporal/...

generate.mocks: $(MOCKGEN) ## Generate mocks
	$(MOCKGEN) -destination=internal/temporal/activities/mocks/mock_inventory_checker.go \
	           -package=mocks \
	           github.com/romangurevitch/go-training/internal/temporal/activities InventoryChecker

# Helper targets for participants

ORDER_JSON := '{"id": "00000000-0000-0000-0000-000000000001", "line_items": [{"product_id": "00000000-0000-0000-0000-000000000001", "quantity": 1, "price_per_item": "100.00"}]}'

define check_id
	@if [ -z "$(ID)" ]; then \
		echo "\033[31mError: ID is not set.\033[0m"; \
		echo "Please start a workflow first and run the provided 'export ID=...' command,"; \
		echo "or provide it directly: make $(MAKECMDGOALS) ID=order-<uuid>"; \
		exit 1; \
	fi
endef

workflow-auto: ## Start an automated order workflow (usage: make workflow-auto [ID=my-id])
	go run cmd/temporal/client/main.go -workflow=AutoProcessOrder -order=$(ORDER_JSON) -workflow-id=$(ID)

workflow-signal: ## Start a signal-driven order workflow (usage: make workflow-signal [ID=my-id])
	go run cmd/temporal/client/main.go -workflow=ProcessOrder -order=$(ORDER_JSON) -wait=false -workflow-id=$(ID)

workflow-signal-send: ## Send a signal to a workflow (usage: make workflow-signal-send ID=order-<uuid> SIGNAL=pickOrder)
	$(call check_id)
	go run cmd/temporal/client/main.go -signal=$(SIGNAL) -workflow-id=$(ID)

workflow-pick: ## Signal a workflow to 'pick' the order (usage: make workflow-pick)
	@$(MAKE) workflow-signal-send SIGNAL=pickOrder

workflow-ship: ## Signal a workflow to 'ship' the order (usage: make workflow-ship)
	@$(MAKE) workflow-signal-send SIGNAL=shipOrder

workflow-deliver: ## Signal a workflow to 'deliver' the order (usage: make workflow-deliver)
	@$(MAKE) workflow-signal-send SIGNAL=markOrderAsDelivered

workflow-cancel: ## Signal a workflow to 'cancel' the order (usage: make workflow-cancel)
	@$(MAKE) workflow-signal-send SIGNAL=cancelOrder
