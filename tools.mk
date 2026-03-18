# Tools Makefile - Versioned Tool Management
# This file manages all external tools with specific versions for reproducible builds

# Tool versions
GOLANGCI_LINT_VERSION := v1.64.5
MOCKGEN_VERSION := v0.6.0

# Versioned tool paths
GOLANGCI_LINT := ./tools/golangci-lint-$(GOLANGCI_LINT_VERSION)
MOCKGEN := ./tools/mockgen-$(MOCKGEN_VERSION)

# Note: Using versioned binaries directly - no symlinks needed

.PHONY: tools
tools: $(GOLANGCI_LINT) $(MOCKGEN)

.PHONY: clean-tools
clean-tools:
	rm -rf ./tools

# golangci-lint installation with version
$(GOLANGCI_LINT):
	@mkdir -p ./tools
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./tools $(GOLANGCI_LINT_VERSION)
	mv ./tools/golangci-lint $@

# mockgen installation with version
$(MOCKGEN):
	@mkdir -p ./tools
	GOBIN=$(shell pwd)/tools go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
	mv ./tools/mockgen $@

# Help target
.PHONY: help-tools
help-tools:
	@echo "Available tool targets:"
	@echo "  tools        - Install all tools with specific versions"
	@echo "  clean-tools  - Remove all versioned tools"
	@echo ""
	@echo "Tool versions:"
	@echo "  golangci-lint: $(GOLANGCI_LINT_VERSION)"
	@echo "  mockgen:       $(MOCKGEN_VERSION)"
