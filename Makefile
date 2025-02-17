GOLANGCI_LINT_INSTALL_DIR ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOLANGCI_LINT_INSTALL_DIR)/golangci-lint

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	@echo "golangci-lint not found, installing the latest version..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_INSTALL_DIR)
