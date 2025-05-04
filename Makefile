.PHONY: protogen

protogen:
	protoc --go_out=. --go-grpc_out=. ./api/advert.proto
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/advert.proto

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

GOLANGCI_LINT_INSTALL_DIR ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOLANGCI_LINT_INSTALL_DIR)/golangci-lint

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	@echo "golangci-lint not found, installing the latest version..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_INSTALL_DIR)
