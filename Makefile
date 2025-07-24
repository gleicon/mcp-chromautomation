# MCP Chrome Automation Service Makefile

.PHONY: build clean test run-server run-ui deps fmt lint help install

# Build variables
BINARY_NAME=mcp-chromautomation
BUILD_DIR=build
VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "Warning: golangci-lint not installed, running basic vet..."; \
		go vet ./...; \
	fi
	@echo "Linting complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./tests/...
	@echo "Tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run the MCP server
run-server: build
	@echo "Starting MCP server..."
	./$(BUILD_DIR)/$(BINARY_NAME) server

# Run the interactive UI
run-ui: build
	@echo "Starting interactive UI..."
	./$(BUILD_DIR)/$(BINARY_NAME) ui

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install the binary to system PATH
install: build
	@echo "Installing $(BINARY_NAME) to system..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development environment ready"

# Run example
run-example: build
	@echo "Running basic usage example..."
	go run examples/basic_usage.go

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Multi-platform builds complete"

# Check dependencies for security vulnerabilities
security-check:
	@echo "Checking for security vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "Warning: govulncheck not installed, install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Show help
help:
	@echo "MCP Chrome Automation Service"
	@echo ""
	@echo "Available commands:"
	@echo "  build          Build the application"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Install dependencies"
	@echo "  fmt            Format code"
	@echo "  lint           Lint code"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  run-server     Run the MCP server"
	@echo "  run-ui         Run the interactive UI"
	@echo "  run-example    Run the basic usage example"
	@echo "  install        Install binary to system PATH"
	@echo "  dev-setup      Set up development environment"
	@echo "  build-all      Build for multiple platforms"
	@echo "  security-check Check for security vulnerabilities"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build for current platform"
	@echo "  make run-ui         # Start the beautiful CLI interface"
	@echo "  make run-server     # Start MCP server for client connections"
	@echo "  make test           # Run all tests"
	@echo "  make install        # Install to system PATH"