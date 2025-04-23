.PHONY: build clean deps lint test run help

# Binary name
BINARY=changelog-generator
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

help:
	@echo "Make targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  lint     - Run linter"
	@echo "  test     - Run tests"
	@echo "  run      - Run the generator"
	@echo "  help     - Show this help message"

# Build the binary
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY) ./cmd/changelog

# Run the generator
run: build
	$(BUILD_DIR)/$(BINARY)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOGET) github.com/manifoldco/promptui

# Install linting tools
lint-deps:
	@command -v $(GOLINT) >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	}

# Run linting
lint: lint-deps
	$(GOLINT) run

# Run tests
test:
	$(GOTEST) -v ./... 