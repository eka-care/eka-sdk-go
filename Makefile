# ABDM SDK Makefile

# Variables
BINARY_NAME=eka-sdk-go
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
GOSEC=gosec

# Directories
SRC_DIR=.
TEST_DIR=./...
COVERAGE_DIR=coverage
DOCS_DIR=docs

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "ABDM SDK for Go - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the SDK
	@echo "Building ABDM SDK..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(SRC_DIR)

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)
	rm -rf $(DOCS_DIR)

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v $(TEST_DIR)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_DIR)/coverage.out $(TEST_DIR)
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race $(TEST_DIR)

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem $(TEST_DIR)

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	$(GOLINT) run

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	$(GOLINT) run --fix

.PHONY: security
security: ## Run security scan
	@echo "Running security scan..."
	$(GOSEC) ./...

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "Verifying go modules..."
	$(GOMOD) verify

.PHONY: install
install: ## Install dependencies
	@echo "Installing dependencies..."
	$(GOGET) -v -t -d ./...
	$(GOMOD) download

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	mkdir -p $(DOCS_DIR)
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

.PHONY: example
example: ## Run example
	@echo "Running example..."
	$(GOCMD) run example/main.go

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -it --rm $(BINARY_NAME):$(VERSION)

.PHONY: release
release: ## Create a new release
	@echo "Creating release for version $(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

.PHONY: check
check: fmt vet lint security ## Run all checks
	@echo "All checks completed successfully!"

.PHONY: ci
ci: install mod-tidy check test-coverage ## Run CI pipeline
	@echo "CI pipeline completed successfully!"

.PHONY: dev-setup
dev-setup: install install-tools ## Setup development environment
	@echo "Development environment setup completed!"

.PHONY: all
all: clean build test-coverage lint security ## Run all targets
	@echo "All targets completed successfully!"